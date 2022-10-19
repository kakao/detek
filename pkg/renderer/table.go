package renderer

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kakao/detek/pkg/detek"
)

func RenderTableReports(list detek.ReportList, MaxWidth int) string {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.SeparateRows = true

	tw.AppendHeader(table.Row{"ID", "LEVEL", "TYPE", "DESCRIPTION"})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", AutoMerge: true},
		{Name: "LEVEL", AutoMerge: true},
	})

	for _, level := range []detek.SeverityLevel{detek.Unknown, detek.Fatal, detek.Error, detek.Warn, detek.Normal} {
		for _, r := range list.Reports {
			if r.Level != level {
				continue
			}

			tw.AppendRow(table.Row{r.ID, r.Level, "Desc", r.Description})
			if r.Level != detek.Unknown && r.Level != detek.Normal {
				tw.AppendRow(table.Row{r.ID, r.Level, "Expl", r.CurrentState.Explanation})
				tw.AppendRow(table.Row{r.ID, r.Level, "Sol", r.CurrentState.Solution})
			}
			if r.Problem.Data != nil {
				tw.AppendRow(table.Row{r.ID, r.Level, "Prob", r.Problem.String()})
			}
			for _, attach := range r.Attachment {
				tw.AppendRow(table.Row{r.ID, r.Level, "Atta", attach.String()})
			}
		}
	}

	if MaxWidth != 0 {
		tw.SetAllowedRowLength(MaxWidth)
	}

	return tw.Render()
}
