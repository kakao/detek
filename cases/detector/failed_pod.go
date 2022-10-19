package detector

import (
	"fmt"

	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
	v1 "k8s.io/api/core/v1"
)

var _ detek.Detector = &FailedPod{}

type FailedPod struct{}

// GetMeta implements detek.Detector
func (*FailedPod) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "failed_pod",
			Description: "check if there is a pod with a 'Failed' status",
			Labels:      []string{"kubernetes", "pod"},
		},
		Level: detek.Error,
		IfHappened: detek.Description{
			Explanation: `some of pods are in a "Failed" status`,
			Solution:    `check why pods are failed`,
		},
		Required: detek.DependencyMeta{
			collector.KeyK8sCoreV1PodList: {Type: detek.TypeOf(v1.PodList{})},
		},
	}
}

// Do implements detek.Detector
func (i *FailedPod) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	podList, err := detek.Typing[v1.PodList](
		ctx.Get(collector.KeyK8sCoreV1PodList, nil))
	if err != nil {
		return nil, err
	}

	type Problem struct {
		Namespace, Name, Reason string
	}
	problems := []Problem{}

	for _, po := range podList.Items {
		if po.Status.Phase == v1.PodFailed {
			problems = append(problems, Problem{
				Namespace: po.Namespace,
				Name:      po.Name,
				Reason:    i.parseFailReason(po),
			})
		}
	}

	report := &detek.ReportSpec{
		HasPassed:  true,
		Attachment: []detek.JSONableData{{Description: "# of Pods", Data: len(podList.Items)}},
	}
	if len(problems) != 0 {
		report.HasPassed = false
		report.Problem = detek.JSONableData{
			Description: "Failed pod list",
			Data:        problems,
		}
	}
	return report, nil
}

func (*FailedPod) parseFailReason(po v1.Pod) string {
	if po.Status.Message != "" {
		return po.Status.Message
	}

	cs := []v1.ContainerStatus{}
	cs = append(cs, po.Status.InitContainerStatuses...)
	cs = append(cs, po.Status.ContainerStatuses...)
	cs = append(cs, po.Status.EphemeralContainerStatuses...)
	for _, co := range cs {
		if term := co.State.Terminated; term != nil {
			if term.ExitCode != 0 {
				return fmt.Sprintf("container %q, exited with status code %d", co.Name, term.ExitCode)
			}
		}
	}
	return "detek: fail reason not found"
}
