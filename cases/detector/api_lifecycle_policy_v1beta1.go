package detector

import (
	"strconv"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	"k8s.io/api/policy/v1beta1"
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
			Explanation: "The policy/v1beta1 PodSecurityPolicy is deprecated in v" + strconv.Itoa(deprecatedMajor) + "." + strconv.Itoa(deprecatedMinor) + "+, " +
				"unavailable in v" + strconv.Itoa(removedMajor) + "." + strconv.Itoa(removedMinor) + "+",
			Solution: "Use the policy/v1 not policy/v1beta1.",
		},
	}
}

func (*ApiLifecyclePolicyV1Beta1) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	podSecurityPolicyList, err := detek.Typing[v1beta1.PodSecurityPolicyList](ctx.Get(collector.KeyK8sPolicyV1Beta1PodSecurityPolicyList, nil))
	if err != nil {
		return nil, err
	}

	type Problem struct {
		Resource  string
		Name      string
	}
	problems := []Problem{}

	for _, psp := range podSecurityPolicyList.Items {
		problems = append(problems, Problem{
			Resource:  "policy/v1beta1 PodSecurityPolicy",
			Name:      psp.Name,
		})
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
