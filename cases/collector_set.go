package cases

import (
	"strconv"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	"github.com/kakao/detek/pkg/utils"
)

type CollectorSetInitiator func(map[string]string) []detek.Collector

var CollectorSet map[string]CollectorSetInitiator = map[string]CollectorSetInitiator{
	DefaultSet: func(m map[string]string) []detek.Collector {
		collectors := []detek.Collector{
			&collector.K8sClientCollector{KubeconfigPath: m[CONFIG_KUBECONFIG]},
			&collector.K8sCoreV1Collector{},
		}

		var err error
		k8sVersionInfo, err := utils.GetK8sVersion(m[CONFIG_KUBECONFIG])
		if err != nil {
			return nil
		}
		k8sVersion, err := strconv.ParseFloat(k8sVersionInfo.Major+"."+k8sVersionInfo.Minor, 64)
		if err != nil {
			return nil
		}
		if k8sVersion >= K8S_VERSION_1_21 && k8sVersion < K8S_VERSION_1_25 {
			collectors = append(collectors, &collector.K8sPolicyV1Beta1Collector{})
		}

		return collectors
	},
	// add more preset here
}
