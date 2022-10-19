package detek

import (
	"fmt"
)

// FD is short for Fake Dependency
type FD struct {
	Key           string
	Value         interface{}
	ShouldProduce bool
	ShouldConsume bool
}

type FakeCollector struct {
	Name      string
	Required  []FD
	Producing []FD
	IsError   bool
	IsPanic   bool
}

func (i FakeCollector) GetMeta() CollectorInfo {
	Required := make(DependencyMeta)
	for _, d := range i.Required {
		Required[d.Key] = DependencyInfo{Type: TypeOf(d.Value)}
	}
	Producing := make(DependencyMeta)
	for _, d := range i.Producing {
		Producing[d.Key] = DependencyInfo{Type: TypeOf(d.Value)}
	}

	return CollectorInfo{
		MetaInfo:  MetaInfo{ID: i.Name},
		Required:  Required,
		Producing: Producing,
	}
}
func (i FakeCollector) Do(ctx DetekContext) error {
	if i.IsError {
		return fmt.Errorf("dummy error: %s", i.Name)
	} else if i.IsPanic {
		panic(i.Name)
	}

	for _, r := range i.Required {
		if r.ShouldConsume {
			val, err := ctx.Get(r.Key, nil)
			if err != nil {
				return err
			}
			if val.Type.Kind() != TypeOf(r.Value).Kind() {
				return fmt.Errorf("Wrong Type Received")
			}
		}
	}
	for _, p := range i.Producing {
		if p.ShouldProduce {
			err := ctx.Set(p.Key, p.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type FakeDetector struct {
	Name        string
	Required    []FD
	IsError     bool
	IsPanic     bool
	ShoudPassed bool
}

func (i FakeDetector) GetMeta() DetectorInfo {
	Required := make(DependencyMeta)
	for _, d := range i.Required {
		Required[d.Key] = DependencyInfo{Type: TypeOf(d.Value)}
	}
	return DetectorInfo{
		MetaInfo: MetaInfo{ID: i.Name},
		Required: Required,
		Level:    Error,
		IfHappened: Description{
			Explanation: "Intended Failure",
			Solution:    "Detect this properly",
		},
	}
}
func (i FakeDetector) Do(ctx DetekContext) (*ReportSpec, error) {
	if i.IsError {
		return nil, fmt.Errorf("dummy error: %s", i.Name)
	} else if i.IsPanic {
		panic(i.Name)
	}

	for _, r := range i.Required {
		if r.ShouldConsume {
			val, err := ctx.Get(r.Key, nil)
			if err != nil {
				return nil, err
			}
			if val.Type.Kind() != TypeOf(r.Value).Kind() {
				return nil, fmt.Errorf("Wrong Type Received")
			}
		}
	}
	return &ReportSpec{
		HasPassed: i.ShoudPassed,
	}, nil
}
