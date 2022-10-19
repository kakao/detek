package cases

import (
	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
)

type CollectorSetInitiator func(map[string]string) []detek.Collector

var CollectorSet map[string]CollectorSetInitiator = map[string]CollectorSetInitiator{
	DefaultSet: func(m map[string]string) []detek.Collector {
		return []detek.Collector{
			&collector.K8sClientCollector{KubeconfigPath: m[CONFIG_KUBECONFIG]},
			&collector.K8sCoreV1Collector{},
		}
	},
	// add more preset here
}
