package renderer

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kakao/detek/pkg/detek"
)

func RenderTablePlan(collectors []detek.Collector, detectors []detek.Detector) string {
	tw := table.NewWriter()
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.SeparateRows = true

	tw.AppendHeader(table.Row{"SEQ", "ID", "TYPE", "DESCRIPTION"})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Name: "SEQ", AutoMerge: true},
		{Name: "ID", AutoMerge: true},
		{Name: "TYPE", AutoMerge: true},
	})
	seq := 0
	for _, c := range collectors {
		seq += 1
		meta := c.GetMeta()
		tw.AppendRow(table.Row{fmt.Sprintf("collector-%d", seq), meta.ID, "desc", "-", meta.Description})
		for key, info := range meta.Required {
			tw.AppendRow(table.Row{fmt.Sprintf("collector-%d", seq), meta.ID, "consume", key, info.Type.String()})
		}
		for key, info := range meta.Producing {
			tw.AppendRow(table.Row{fmt.Sprintf("collector-%d", seq), meta.ID, "produce", key, info.Type.String()})
		}
	}
	seq = 0
	for _, d := range detectors {
		seq += 1
		meta := d.GetMeta()
		tw.AppendRow(table.Row{fmt.Sprintf("detctor-%d", seq), meta.ID, "desc", "description", meta.Description})
		tw.AppendRow(table.Row{fmt.Sprintf("detctor-%d", seq), meta.ID, "desc", "severity", meta.Level})
		for key, info := range meta.Required {
			tw.AppendRow(table.Row{fmt.Sprintf("detctor-%d", seq), meta.ID, "consume", key, info.Type.String()})
		}
	}
	return tw.Render()
}
