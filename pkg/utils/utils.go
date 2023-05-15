package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := []K{}
	for k := range m {
		r = append(r, k)
	}
	return r
}

func GetK8sVersion(kubeconfigPath string) (*version.Info, error) {
	var config *rest.Config
	var err error
	if kubeconfigPath != "" {
		//   1. kubeconfig file located by "--kubeconfig" flag
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
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

	if err != nil {
		return nil, fmt.Errorf("fail to get kubernetes client:%w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("fail to generate client config from kubeconfig:%w", err)
	}

	version, err := clientset.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("fail to get kubernetes server version:%w", err)
	}

	return version, nil
}
