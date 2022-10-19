package renderer

import (
	"testing"
	"time"

	"github.com/kakao/detek/pkg/detek"
)

func TestRenderTableReports(t *testing.T) {
	type args struct {
		list     detek.ReportList
		MaxWidth int
	}
	tests := []struct {
		name string
		args args
		want func(string) bool
	}{
		{
			name: "test 1",
			args: args{
				list: detek.ReportList{
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
			},
			want: func(s string) bool { return len(s) != 0 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderTableReports(tt.args.list, tt.args.MaxWidth); !tt.want(got) {
				t.Errorf("want() return false")
			}
		})
	}
}
