package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kakao/detek/pkg/detek"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	KeyK8sVersion    = "kubernetes_version"
	KeyK8sRestConfig = "kubernetes_rest_config"
	KeyK8sClient     = "kubernetes_client"
)

var _ detek.Collector = &K8sClientCollector{}

type K8sClientCollector struct {
	KubeconfigPath string
}

func (*K8sClientCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "kubernetes_client",
			Description: "generate kubernetes client and some basic information using ENV",
			Labels:      []string{"kubernetes", "client"},
		},
		Required: detek.DependencyMeta{ /* NOTHING */ },
		Producing: detek.DependencyMeta{
			KeyK8sClient:     {Type: detek.TypeOf(&kubernetes.Clientset{})},
			KeyK8sRestConfig: {Type: detek.TypeOf(&rest.Config{})},
			KeyK8sVersion:    {Type: detek.TypeOf(version.Info{})},
		},
	}
}

func (c *K8sClientCollector) Do(ctx detek.DetekContext) error {

	var config *rest.Config
	var err error
	if c.KubeconfigPath != "" {
		//   1. kubeconfig file located by "--kubeconfig" flag
		config, err = clientcmd.BuildConfigFromFlags("", c.KubeconfigPath)
	} else if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
		//   2. kubeconfig file located by "KUBECONFIG" env
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		//   3. in-cluster client configuration (useful when using detek in a kubernetes cluster)
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		//   4. kubeconfig file located in default directory ($HOME/.kube/config)`,
		kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	config.WarningHandler = rest.NoWarnings{}

	if err != nil {
		return fmt.Errorf("fail to get kubernetes client:%w", err)
	}
	if err := ctx.Set(KeyK8sRestConfig, config); err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("fail to generate client config from kubeconfig:%w", err)
	}
	if err := ctx.Set(KeyK8sClient, clientset); err != nil {
		return err
	}

	vi, err := clientset.ServerVersion()
	if err != nil {
		return fmt.Errorf("fail to get kubernetes server version:%w", err)
	}
	if err := ctx.Set(KeyK8sVersion, *vi); err != nil {
		return err
	}
	return nil
}
