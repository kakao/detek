package detector

import (
	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/core/v1"
)

var _ detek.Detector = &PodWithoutRequests{}

type PodWithoutRequests struct {
	DoNotCheckCPU    bool
	DoNotChekcMemory bool
}

// GetMeta implements detek.Detector
func (d *PodWithoutRequests) GetMeta() detek.DetectorInfo {
	desc := "Finding pod without requests"
	if d.DoNotCheckCPU {
		desc += " (will not check CPU requests)"
	}
	if d.DoNotChekcMemory {
		desc += " (will not check Mem requests)"
	}
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "pod_without_requests",
			Description: desc,
			Labels:      []string{"kubernetes", "pod"},
		},
		Level: detek.Warn,
		IfHappened: detek.Description{
			Explanation: "There is a Pod without memory or cpu requests (which could cause unwanted overcommit)",
			Solution:    "Set requests for those Pods",
		},
		Required: detek.DependencyMeta{
			collector.KeyK8sCoreV1PodList: {Type: detek.TypeOf(v1.PodList{})},
		},
	}
}

// Do implements detek.Detector
func (d *PodWithoutRequests) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	podList, err := detek.Typing[v1.PodList](
		ctx.Get(collector.KeyK8sCoreV1PodList, nil))
	if err != nil {
		return nil, err
	}

	type Problem struct {
		Namespace   string
		Name        string
		Container   string
		CPURequests bool
		MemRequests bool
	}
	problems := []Problem{}

	for _, po := range podList.Items {
		for _, co := range po.Spec.Containers {
			p := Problem{
				Namespace:   po.Namespace,
				Name:        po.Name,
				Container:   co.Name,
				CPURequests: true,
				MemRequests: true,
			}
			if _, ok := co.Resources.Requests[v1.ResourceCPU]; !ok {
				p.CPURequests = false
			}
			if _, ok := co.Resources.Requests[v1.ResourceMemory]; !ok {
				p.MemRequests = false
			}
			if !p.CPURequests && !d.DoNotCheckCPU {
				problems = append(problems, p)
				continue
			}
			if !p.MemRequests && !d.DoNotChekcMemory {
				problems = append(problems, p)
				continue
			}
		}
	}

	return &detek.ReportSpec{
		HasPassed: len(problems) == 0,
		Problem: detek.JSONableData{
			Description: "list of Pods with no limits",
			Data:        problems,
		},
		Attachment: []detek.JSONableData{{Description: "# of evaluated Pods", Data: len(podList.Items)}},
	}, nil
}
