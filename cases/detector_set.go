package cases

import (
	"github.com/kakao/detek/cases/detector"
	"github.com/kakao/detek/pkg/detek"
)

type DetectorSetInitiator func(map[string]string) []detek.Detector

var (
	DetectorSet map[string]DetectorSetInitiator = map[string]DetectorSetInitiator{
		DefaultSet: func(m map[string]string) []detek.Detector {
			return []detek.Detector{
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
		},
		// add more preset here
	}
)
