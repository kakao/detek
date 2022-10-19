package detector

import (
	"fmt"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/core/v1"
)

var _ detek.Detector = &ServicePartiallyAvailable{}

type ServicePartiallyAvailable struct{}

func (*ServicePartiallyAvailable) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "service_partially_available",
			Description: "Finding services that are partially working right now.",
			Labels:      []string{"kubernetes", "service"},
		},
		Required: detek.DependencyMeta{
			collector.KeyK8sCoreV1EndpointList: {Type: detek.TypeOf(v1.EndpointsList{})},
		},
		Level: detek.Error,
		IfHappened: detek.Description{
			Explanation: "Some of Pods are disabled for this Service. This Service is downgraded now.",
			Solution:    "Check all your Pods are functional right now.",
		},
	}
}

func (*ServicePartiallyAvailable) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	epList, err := detek.Typing[v1.EndpointsList](
		ctx.Get(collector.KeyK8sCoreV1EndpointList, nil))
	if err != nil {
		return nil, err
	}

	type Problem struct {
		Name              string
		Namespace         string
		NotReadyEndpoints []string
	}
	problems := []Problem{}

	for _, ep := range epList.Items {
		for _, sub := range ep.Subsets {
			if len(sub.NotReadyAddresses) != 0 {
				notReadies := []string{}
				for _, addr := range sub.NotReadyAddresses {
					text := addr.IP
					if ref := addr.TargetRef; ref != nil {
						text += fmt.Sprintf(" (%s/%s/%s)", ref.Kind, ref.Namespace, ref.Name)
					}
					notReadies = append(notReadies, text)
				}
				problems = append(problems, Problem{
					Name:              ep.Name,
					Namespace:         ep.Namespace,
					NotReadyEndpoints: notReadies,
				})
			}
		}
	}
	return &detek.ReportSpec{
		HasPassed: len(problems) == 0,
		Problem: detek.JSONableData{
			Description: "Unavailable Service List",
			Data:        problems,
		},
		Attachment: []detek.JSONableData{
			{Description: "# of evaluated Endpoints", Data: len(epList.Items)},
		},
	}, nil
}
