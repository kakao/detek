package collector

import (
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	KeyK8sAppsV1DeploymentList = "kubernetes_apps_v1_deployment_list"
)

var _ detek.Collector = &K8sAppsV1Collector{}

type K8sAppsV1Collector struct{}

func (*K8sAppsV1Collector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "kubernetes_apps_v1",
			Description: "collect apps v1 resources from kubernetes",
			Labels:      []string{"kubernetes", "apps/v1", "manifests"},
		},
		Required: detek.DependencyMeta{
			KeyK8sClient: {Type: detek.TypeOf(&kubernetes.Clientset{})},
		},
		Producing: detek.DependencyMeta{
			KeyK8sAppsV1DeploymentList: {Type: detek.TypeOf(v1.DeploymentList{})},
		},
	}
}

func (*K8sAppsV1Collector) Do(ctx detek.DetekContext) error {
	panic("unimplemented")
}
