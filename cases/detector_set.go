package cases

import (
	"strconv"

	"github.com/kakao/detek/cases/detector"
	"github.com/kakao/detek/pkg/detek"
	"github.com/kakao/detek/pkg/utils"
)

type DetectorSetInitiator func(map[string]string) []detek.Detector

var (
	DetectorSet map[string]DetectorSetInitiator = map[string]DetectorSetInitiator{
		DefaultSet: func(m map[string]string) []detek.Detector {
			detectors := []detek.Detector{
				&detector.FailedPod{},
				&detector.PodWithoutLimits{
					DoNotCheckCPU:    true, // Disable Checking CPU Limits
					DoNotChekcMemory: false,
				},
				&detector.PodWithoutRequests{
					DoNotCheckCPU:    false,
					DoNotChekcMemory: false,
				},
				&detector.PodWithoutLivenessProbe{},
				&detector.PodWithoutReadinessProbe{},
				&detector.ServiceNoAvailableTarget{},
				&detector.ServicePartiallyAvailable{},
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
				detectors = append(detectors, &detector.ApiLifecyclePolicyV1Beta1{})
			}

			return detectors
		},
		// add more preset here
	}
)
