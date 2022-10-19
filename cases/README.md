detek has two main components, `Collector`and `Detector`. The `Collector` collects information by fetching data from external APIs (e.g, Kubernetes API). The `Detector`, on the other hand, by using data collected by the `Collector`, determines whether a specific situation has happened or not.


## Collector
The `Collector` is one of the interfaces in detek, which only has two methods.

```go
type Collector interface {
	// Give detek metadata of this Collector
	GetMeta() CollectorInfo

	// Collect Data using external API and store it in DetekContext
	Do(ctx DetekContext) error 
}
```

below is the example of the `Pod Collector`, which collects Pod manifests using Kubernetes API and save it to `DetekContext`.

```go
const (
	KeyK8sCoreV1PodList  = "kubernetes_core_v1_podlist"
)

// just type hinting. Not necessary.
// "K8sCoreV1PodCollector" is a implemntation of the "Collector"
var _ detek.Collector = &K8sCoreV1PodCollector{}

type K8sCoreV1PodCollector struct{}

// Give detek a metadata of this collector
// GetMeta implements detek.Collector
func (*K8sCoreV1PodCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "kubernetes_core_v1_pod",
			Description: "collect core v1 pod resources from kubernetes",
			Labels:      []string{"kubernetes", "core/v1", "pod", "manifest"},
		},
		Required: detek.DependencyMeta{
			// K8sCoreV1PodCollector need K8sClient (may collected by the other Collector)
			KeyK8sClient: {Type: detek.TypeOf(&kubernetes.Clientset{})},
		},
		Producing: detek.DependencyMeta{
			// K8sCoreV1PodCollector will produce PodList
			KeyK8sCoreV1PodList:  {Type: detek.TypeOf(v1.PodList{})},
		},
	}
}

// Collecting logic here
// detek will execute this, and will expect to get PodList as declared at GetMeta
// Do implements detek.Collector
func (*K8sCoreV1Collector) Do(dctx detek.DetekContext) error {
	// Get Kubernetes Client (collected by the other Collector)
	c, err := detek.Typing[*kubernetes.Clientset]( // <- This is a syntactic sugar.
		dctx.Get(KeyK8sClient, nil), // <- Get data from a detek Store
	)
	if err != nil {
		return fmt.Errorf("fail to get kubernetes client: %w", err)
	}
  
  // Get Pod List
  podList, err := c.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
    return fmt.Errorf("fail to get Pod list from kubernetes: %w", err)
	}
  
  // Give Pod List to detek (as declared at GetMeta)
  return dctx.Set(KeyK8sCoreV1PodList, *podList) // <- Set data to a detek Store
}
```

`Collector` (`Pod Collector`, in this case) can be used by appending it on [cases/collector_set.go](./collector_set.go).

```go
var CollectorSet map[string]CollectorSetInitiator = map[string]CollectorSetInitiator{
	DefaultSet: func(m map[string]string) []detek.Collector {
		return []detek.Collector{
			&collector.K8sClientCollector{KubeconfigPath: m[CONFIG_KUBECONFIG]},
			&collector.K8sCoreV1PodCollector{}, // APPENDED
		}
	},
	// add more preset here
}
```


## Detector

The `Detector` is also one of the interfaces in detek, with two methods.

```go
type Detector interface {
	// Give detek metadata of this Collector
	GetMeta() DetectorInfo

	// Find issues using data stored by Collector. 
	Do(ctx DetekContext) (*ReportSpec, error)
}
```

With collected data (stored in `DetekContext`), `Detector` can use those data to validate whether the cluster has an issue or not. Below is the example of the `failed_pod Detector`, which tries to find any pod with a **Failed** state.

```go
// just type hinting. Not necessary.
var _ detek.Detector = &FailedPod{}

type FailedPod struct{}

// Give detek a metadata of this detector
// GetMeta implements detek.Detector
func (*FailedPod) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "failed_pod",
			Description: "check if there is a pod with a 'Failed' status",
			Labels:      []string{"kubernetes", "pod"},
		},
		// Severity of this case is (if happened)
		Level: detek.Error,
		// What user can do (if happened)
		IfHappened: detek.Description{
			Explanation: `some of pods are in a "Failed" status`,
			Solution:    `check why pods are failed`,
		},
		Required: detek.DependencyMeta{
			// This Detector requires following dependency
			// (if the dependency not exists, it will not be run by detek)
			collector.KeyK8sCoreV1PodList: {Type: detek.TypeOf(v1.PodList{})},
		},
	}
}

// Do implements detek.Detector
func (i *FailedPod) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	// fetch PodList from a detek
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
		// if Pod's Phase is Failed, than append to the Problem List
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
  // if found some problem
  if len(problems) != 0 {
		// Detector reports this is not passed 
		report.HasPassed = false
		// and shows users what is the problem.
		report.Problem = detek.JSONableData{
			Description: "Failed pod list",
			Data:        problems,
		}
	}
	return report, nil
}
```

`Detector` (`FailedPod Detector`, in this case) can be used by appending it on [cases/detector_set.go](./detector_set.go).

```go
var (
	DetectorSet map[string]DetectorSetInitiator = map[string]DetectorSetInitiator{
		DefaultSet: func(m map[string]string) []detek.Detector {
			return []detek.Detector{
				&detector.FailedPod{},
			}
		},
		// add more preset here
	}
)
```