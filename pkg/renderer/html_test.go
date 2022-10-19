package renderer

import (
	"strings"
	"testing"
	"time"

	"github.com/kakao/detek/pkg/detek"
	"golang.org/x/net/html"
)

func TestRenderHTMLReports(t *testing.T) {
	type args struct {
		r detek.ReportList
	}
	tests := []struct {
		name string
		args args
		want func(string) bool
	}{
		{
			name: "test - 1",
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
			},
			want: func(s string) bool {
				if _, err := html.Parse(strings.NewReader(s)); err != nil {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RenderHTMLReports(tt.args.r); !tt.want(got) {
				t.Errorf("want() returned false")
			}
		})
	}
}

func generateDummyReport(id string, level detek.SeverityLevel) detek.Report {
	return detek.Report{
		MetaInfo: detek.MetaInfo{
			ID:          id,
			Description: "this is test description " + id,
			Labels:      []string{"test", "kubernetes", id},
		},
		Level: level,
		CurrentState: detek.Description{
			Explanation: string(level) + " state",
			Solution:    "do nothing",
		},
		ReportSpec: detek.ReportSpec{
			HasPassed: false,
			Problem: detek.JSONableData{
				Description: "something is happend",
				Data: map[string]string{
					"Hello":       "World",
					string(level): id,
				},
			},
			Attachment: []detek.JSONableData{
				{
					Description: "# of something",
					Data: map[string]string{
						"This": "Attachment",
					},
				},
				{
					Description: "# of something",
					Data: map[string]string{
						"This": "Attachment",
					},
				},
			},
		},
		CreatedAt: time.Now(),
	}
}
