package detek

import (
	"encoding/json"
	"time"
)

type SeverityLevel string

type JSONableData struct {
	Description string `json:"description"`

	// You can set any "JSON-Marshal-able" data in here
	// one depth JSON (object or array) is recommended, but not required.
	Data any `json:"data"`
}

func (j *JSONableData) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		b = []byte(err.Error())
	}
	return string(b)
}

const (
	// Fatal if some feature or statuses are 100% fully not functional.
	// It is likely that the cluster already in catastrophic situation.
	Fatal SeverityLevel = "Fatal"

	// Error if some feature or statuses are partially worked.
	// Somehow Kubernetes magic makes services work, but need some fix.
	Error SeverityLevel = "Error"

	// Warn if some features or statuses are fully works, but there's something not recommended there.
	// e.g, not setting limits in pod spec is not recommended, but doesn't harm your service.
	Warn SeverityLevel = "Warn"

	// Normal if everything is fine
	Normal SeverityLevel = "Normal"

	// Something unexpected occur, can not examine severity
	Unknown SeverityLevel = "Unknown"
)

func (s *SeverityLevel) ToInt() int {
	if s == nil {
		return 0
	}
	v, ok := map[SeverityLevel]int{
		Fatal:  4,
		Error:  3,
		Warn:   2,
		Normal: 1,
	}[*s]
	if !ok {
		return 0
	} else {
		return v
	}
}

type Report struct {
	MetaInfo
	CreatedAt time.Time     `json:"created_at"`
	Level     SeverityLevel `json:"level"`

	CurrentState Description `json:"-"`
	ReportSpec
}

type ReportSpec struct {
	// Is this Passed?
	HasPassed bool `json:"-"`

	// Attachment to show the causes of problem.
	Problem JSONableData `json:"problem,omitempty"`

	// Attachment for debugging purpose
	Attachment []JSONableData `json:"attachment,omitempty"`
}

type ReportList struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	Reports []Report `json:"reports"`
}
