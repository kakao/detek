package renderer

import (
	"fmt"
	"testing"

	"github.com/kakao/detek/pkg/detek"
)

func TestRenderTablePlan(t *testing.T) {
	type args struct {
		collectors []detek.Collector
		detectors  []detek.Detector
	}
	tests := []struct {
		name string
		args args
		want func(string) bool
	}{
		{
			name: "empty",
			args: args{},
			want: func(s string) bool {
				fmt.Println(s)
				return len(s) != 0
			},
		},
		{
			name: "something",
			args: args{
				collectors: []detek.Collector{dummyCollector{}},
				detectors:  []detek.Detector{dummyDetector{}},
			},
			want: func(s string) bool {
				fmt.Println(s)
				return len(s) != 0
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTablePlan(tt.args.collectors, tt.args.detectors); !tt.want(got) {
				t.Errorf("want() return false")
			}
		})
	}
}

var _ detek.Collector = dummyCollector{}

type dummyCollector struct{}

// Do implements detek.Collector
func (dummyCollector) Do(ctx detek.DetekContext) error {
	panic("unimplemented")
}

// GetMeta implements detek.Collector
func (dummyCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "dummy",
			Description: "dummydescription",
			Labels:      []string{"label"},
		},
		Required: detek.DependencyMeta{
			"dummy": detek.DependencyInfo{Type: detek.TypeOf("")},
		},
		Producing: detek.DependencyMeta{
			"dummydummy": detek.DependencyInfo{Type: detek.TypeOf("")},
		},
	}
}

var _ detek.Detector = dummyDetector{}

type dummyDetector struct{}

// Do implements detek.Detector
func (dummyDetector) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	panic("unimplemented")
}

// GetMeta implements detek.Detector
func (dummyDetector) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "dummy",
			Description: "dummydescription",
			Labels:      []string{"label"},
		},
		Required: detek.DependencyMeta{
			"dummydummy": detek.DependencyInfo{Type: detek.TypeOf("")},
		},
		IfHappened: detek.NormalStatus,
		Level:      detek.Error,
	}
}
