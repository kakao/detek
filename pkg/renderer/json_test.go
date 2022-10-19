package renderer

import (
	"testing"
	"time"

	"github.com/kakao/detek/pkg/detek"
)

func TestRenderJSONReports(t *testing.T) {
	type args struct {
		r      detek.ReportList
		pretty bool
	}
	tests := []struct {
		name string
		args args
		want func(string) bool
	}{
		{
			name: "pretty",
			args: args{
				r: detek.ReportList{
					StartedAt:  time.Now(),
					FinishedAt: time.Now(),
					Reports: []detek.Report{
						generateDummyReport("1", detek.Fatal),
						generateDummyReport("2", detek.Fatal),
						generateDummyReport("3", detek.Warn),
						generateDummyReport("4", detek.Warn),
						generateDummyReport("5", detek.Error),
						generateDummyReport("6", detek.Error),
						generateDummyReport("7", detek.Normal),
						generateDummyReport("8", detek.Normal),
						generateDummyReport("9", detek.Unknown),
						generateDummyReport("10", detek.Unknown),
					},
				},
				pretty: true,
			},
			want: func(s string) bool { return len(s) != 0 },
		},
		{
			name: "ugly",
			args: args{
				r: detek.ReportList{
					StartedAt:  time.Now(),
					FinishedAt: time.Now(),
					Reports: []detek.Report{
						generateDummyReport("1", detek.Fatal),
						generateDummyReport("2", detek.Fatal),
						generateDummyReport("3", detek.Warn),
						generateDummyReport("4", detek.Warn),
						generateDummyReport("5", detek.Error),
						generateDummyReport("6", detek.Error),
						generateDummyReport("7", detek.Normal),
						generateDummyReport("8", detek.Normal),
						generateDummyReport("9", detek.Unknown),
						generateDummyReport("10", detek.Unknown),
					},
				},
				pretty: false,
			},
			want: func(s string) bool { return len(s) != 0 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderJSONReports(tt.args.r, tt.args.pretty); !tt.want(got) {
				t.Errorf("want() return false")
			}
		})
	}
}
