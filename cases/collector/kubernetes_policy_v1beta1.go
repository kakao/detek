package collector

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/kakao/detek/pkg/detek"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	KeyK8sPolicyV1Beta1PodSecurityPolicyList = "kubernetes_policy_v1beta1_pod_security_policy_list"
)

var _ detek.Collector = &K8sPolicyV1Beta1Collector{}

type K8sPolicyV1Beta1Collector struct{}

func (*K8sPolicyV1Beta1Collector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "kubernetes_policy_v1beta1",
			Description: "collect extensions v1beta1 resources from kubernetes",
			Labels:      []string{"kubernetes", "policy/v1beta1", "manifests"},
		},
		Required: detek.DependencyMeta{
			KeyK8sClient: {Type: detek.TypeOf(&kubernetes.Clientset{})},
		},
		Producing: detek.DependencyMeta{
			KeyK8sPolicyV1Beta1PodSecurityPolicyList: {Type: detek.TypeOf(v1beta1.PodSecurityPolicyList{})},
		},
	}
}

func (*K8sPolicyV1Beta1Collector) Do(dctx detek.DetekContext) error {
	c, err := detek.Typing[*kubernetes.Clientset](
		dctx.Get(KeyK8sClient, nil),
	)
	if err != nil {
		return fmt.Errorf("fail to get kubernetes client: %w", err)
	}
	var errs = &multierror.Error{}

	ctx := dctx.Context()

	podSecurityPolicyList, err := c.PolicyV1beta1().PodSecurityPolicies().List(ctx, metav1.ListOptions{})
	errs = multierror.Append(errs, err)
	errs = multierror.Append(errs,
		dctx.Set(KeyK8sPolicyV1Beta1PodSecurityPolicyList, *podSecurityPolicyList),
	)

	return errs.ErrorOrNil()
}
