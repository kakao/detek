package detector

import (
	"fmt"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ detek.Detector = &PodWithoutReadinessProbe{}

type PodWithoutReadinessProbe struct{}

// GetMeta implements detek.Detector
func (*PodWithoutReadinessProbe) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "pod_without_readiness_probe",
			Description: "Finding pod without ReadinessProbe (it will check a Pod referenced by Service only)",
			Labels:      []string{"kubernetes", "pod", "probe"},
		},
		Required: detek.DependencyMeta{
			collector.KeyK8sCoreV1PodList:      {Type: detek.TypeOf(v1.PodList{})},
			collector.KeyK8sCoreV1EndpointList: {Type: detek.TypeOf(v1.EndpointsList{})},
		},
		Level: detek.Warn,
		IfHappened: detek.Description{
			Explanation: "Without ReadinessProbe, Kubernetes can not cut traffics to unavailable pod automatically, which could harm your service reliability.",
			Solution:    "Recommend defining a proper ReadinessProbe",
		},
	}
}

// Do implements detek.Detector
func (*PodWithoutReadinessProbe) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	epList, err := detek.Typing[v1.EndpointsList](
		ctx.Get(collector.KeyK8sCoreV1EndpointList, nil))
	if err != nil {
		return nil, err
	}
	podList, err := detek.Typing[v1.PodList](
		ctx.Get(collector.KeyK8sCoreV1PodList, nil))
	if err != nil {
		return nil, err
	}

	podFilter := make(map[types.UID]v1.Endpoints)
	for _, ep := range epList.Items {
		for _, sub := range ep.Subsets {
			for _, addr := range append(sub.Addresses, sub.NotReadyAddresses...) {
				if addr.TargetRef == nil {
					continue
				}
				podFilter[addr.TargetRef.UID] = ep
			}
		}
	}

	targetPods := []v1.Pod{}
	for _, po := range podList.Items {
		if _, ok := podFilter[po.UID]; ok {
			targetPods = append(targetPods, po)
		}
	}

	type Problem struct {
		Namespace    string
		Name         string
		Container    string
		Owner        string
		ReferencedBy string
	}
	problems := []Problem{}

	// Check targeted Pods having ReadinessProbe
	for _, po := range targetPods {
		for _, co := range po.Spec.Containers {
			if co.ReadinessProbe == nil {
				ep := podFilter[po.UID]
				OwnerString := ""
				for _, o := range po.OwnerReferences {
					OwnerString += fmt.Sprintf("%s/%s", o.Kind, o.Name)
				}
				problems = append(problems, Problem{
					Namespace:    po.Namespace,
					Name:         po.Name,
					Container:    co.Name,
					Owner:        OwnerString,
					ReferencedBy: fmt.Sprintf("Service/%s", ep.Name),
				})
			}
		}
	}

	return &detek.ReportSpec{
		HasPassed: len(problems) == 0,
		Problem: detek.JSONableData{
			Description: "Pods without ReadinessProbe",
			Data:        problems,
		},
		Attachment: []detek.JSONableData{
			{Description: "# of evaluated Pod", Data: len(targetPods)},
			{Description: "# of evaluated Endpoints", Data: len(podFilter)},
		},
	}, nil
}
