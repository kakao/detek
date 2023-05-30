package detector

import (
	"fmt"
	"strconv"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/version"
)

var _ detek.Detector = &ApiLifecyclePolicyV1Beta1{}

type ApiLifecyclePolicyV1Beta1 struct{}

func (*ApiLifecyclePolicyV1Beta1) GetMeta() detek.DetectorInfo {
	psp := &v1beta1.PodSecurityPolicy{}
	deprecatedMajor, deprecatedMinor := psp.APILifecycleDeprecated()
	removedMajor, removedMinor := psp.APILifecycleRemoved()

	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "api_lifecycle - policy/v1beta1",
			Description: "Finding api that is unavailable in the future Kubernetes releases.",
			Labels:      []string{"kubernetes", "api"},
		},
		Required: detek.DependencyMeta{
			collector.KeyK8sPolicyV1Beta1PodSecurityPolicyList: {Type: detek.TypeOf(v1beta1.PodSecurityPolicyList{})},
		},
		Level: detek.Warn,
		IfHappened: detek.Description{
			Explanation: fmt.Sprintf("The PodSecurityPolicy in policy/v1beta1 is deprecated in %s, and unavailable in %s",
				fmt.Sprintf("v%d.%d", deprecatedMajor, deprecatedMinor),
				fmt.Sprintf("v%d.%d", removedMajor, removedMinor),
			),
			Solution: "Migrate to Pod Security Admission or a 3rd party admission webhook. " +
				"For more information, please refer the following guide. " +
				"https://kubernetes.io/docs/reference/using-api/deprecation-guide/#psp-v125",
		},
	}
}

func (*ApiLifecyclePolicyV1Beta1) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	var err error

	podSecurityPolicyList, err := detek.Typing[v1beta1.PodSecurityPolicyList](ctx.Get(collector.KeyK8sPolicyV1Beta1PodSecurityPolicyList, nil))
	if err != nil {
		return nil, err
	}

	type Problem struct {
		Resource string
		Name     string
	}
	problems := []Problem{}

	version, err := detek.Typing[version.Info](ctx.Get(collector.KeyK8sVersion, nil))
	if err != nil {
		return nil, err
	}
	currentVersion, err := strconv.ParseFloat(version.Major+"."+version.Minor, 64)
	if err != nil {
		return nil, err
	}

	if currentVersion >= K8S_VERSION_1_21 {
		for _, psp := range podSecurityPolicyList.Items {
			problems = append(problems, Problem{
				Resource: "policy/v1beta1 PodSecurityPolicy",
				Name:     psp.Name,
			})
		}
	}

	return &detek.ReportSpec{
		HasPassed: len(problems) == 0,
		Problem: detek.JSONableData{
			Description: "PodSecurityPolicies using policy/v1beta1 API",
			Data:        problems,
		},
		Attachment: []detek.JSONableData{
			{Description: "# of evaluated PodSecurityPolicies", Data: len(podSecurityPolicyList.Items)},
		},
	}, nil
}
