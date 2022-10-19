package collector

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	KeyK8sCoreV1NodeList     = "kubernetes_core_v1_nodelist"
	KeyK8sCoreV1PodList      = "kubernetes_core_v1_podlist"
	KeyK8sCoreV1ServiceList  = "kubernetes_core_v1_servicelist"
	KeyK8sCoreV1EndpointList = "kubernetes_core_v1_endpointlist"
)

var _ detek.Collector = &K8sCoreV1Collector{}

type K8sCoreV1Collector struct{}

func (*K8sCoreV1Collector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "kubernetes_core_v1",
			Description: "collect core v1 resources from kubernetes",
			Labels:      []string{"kubernetes", "core/v1", "manifest"},
		},
		Required: detek.DependencyMeta{
			KeyK8sClient: {Type: detek.TypeOf(&kubernetes.Clientset{})},
		},
		Producing: detek.DependencyMeta{
			KeyK8sCoreV1PodList:      {Type: detek.TypeOf(v1.PodList{})},
			KeyK8sCoreV1NodeList:     {Type: detek.TypeOf(v1.NodeList{})},
			KeyK8sCoreV1EndpointList: {Type: detek.TypeOf(v1.EndpointsList{})},
			KeyK8sCoreV1ServiceList:  {Type: detek.TypeOf(v1.ServiceList{})},
		},
	}
}

func (*K8sCoreV1Collector) Do(dctx detek.DetekContext) error {
	c, err := detek.Typing[*kubernetes.Clientset](
		dctx.Get(KeyK8sClient, nil),
	)
	if err != nil {
		return fmt.Errorf("fail to get kubernetes client: %w", err)
	}
	var errs = &multierror.Error{}

	ctx := dctx.Context()

	podList, err := c.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	errs = multierror.Append(errs, err)

	errs = multierror.Append(errs,
		dctx.Set(KeyK8sCoreV1PodList, *podList),
	)

	nodeList, err := c.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	errs = multierror.Append(errs, err)
	errs = multierror.Append(errs,
		dctx.Set(KeyK8sCoreV1NodeList, *nodeList),
	)

	serviceList, err := c.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	errs = multierror.Append(errs, err)
	errs = multierror.Append(errs,
		dctx.Set(KeyK8sCoreV1ServiceList, *serviceList),
	)

	epList, err := c.CoreV1().Endpoints("").List(ctx, metav1.ListOptions{})
	errs = multierror.Append(errs, err)
	errs = multierror.Append(errs,
		dctx.Set(KeyK8sCoreV1EndpointList, *epList),
	)

	return errs.ErrorOrNil()
}
