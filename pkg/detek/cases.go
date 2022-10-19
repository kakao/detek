package detek

type MetaInfo struct {
	// Naming Convention (for Dectector)
	//   (the cluster has) abnormal_pod
	//   (the cluster has) not_ready_node
	//   (the cluster has) deployment_without_pdb
	//   (the cluster has) obsolete_dns_mode
	// omit "it has" and "/^kubernetes/g".
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Labels      []string `json:"labels,omitempty"`
}

// Detector Definitions
type Detector interface {
	GetMeta() DetectorInfo
	Do(ctx DetekContext) (*ReportSpec, error)
}

type DetectorInfo struct {
	MetaInfo
	Required DependencyMeta

	// Show how severe this case is when the thing has happened.
	Level SeverityLevel

	// Show what user can do when the thing has happened.
	IfHappened Description `json:"-"`
}

type Description struct {
	// Show what current situation is
	Explanation string `json:"explanation"`
	// Show how to solve this
	Solution string `json:"solution,omitempty"`
}

var (
	NormalStatus Description = Description{
		Explanation: "Everything is normal",
	}
	NoDepStatus Description = Description{
		Explanation: "Not executed, some of required data is not provided",
		Solution:    "Check if every Collectors are executed properly",
	}
	ErrOnDetectorStatus Description = Description{
		Explanation: "Detector is failed with an error",
		Solution:    "This may be a bug in this program. Check an error message",
	}
)

// SeverityLevelDescription is a definition of descriptions that show what is the meaning of each level.
type SeverityLevelDescription struct {
	// Fatal if some feature or statuses are 100% fully not functional.
	// It is likely that the cluster already in catastrophic situation.
	Fatal *Description `json:"fatal,omitempty"`

	// Error if some feature or statuses are partially worked.
	// Somehow kubernetes magic makes some services works, but need some fix.
	Error *Description `json:"error,omitempty"`

	// Warn if some features or statuses are fully workes, but there's something not recommended there.
	// e.g, not setting limits in pod spec is not recommended, but doesn't harm your service.
	Warn *Description `json:"warn,omitempty"`

	// Normal if everything is fine
	Normal *Description `json:"normal,omitempty"`
}

// Collector Definitions
type CollectorInfo struct {
	MetaInfo
	Required  DependencyMeta
	Producing DependencyMeta
}

type Collector interface {
	GetMeta() CollectorInfo
	Do(ctx DetekContext) error
}
