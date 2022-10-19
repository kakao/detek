package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/kakao/detek/pkg/detek"
)

func RenderHTMLReports(r detek.ReportList) string {
	tmpl, err := template.New("html").Funcs(template.FuncMap{
		"marshal": func(v interface{}) template.JS {
			a, _ := json.MarshalIndent(v, "", "  ")
			return template.JS(a)
		},
	}).Parse(_HTMLReportTemplate)
	if err != nil {
		panic(fmt.Errorf("this is bug: %w", err))
	}

	reportMap := make(map[detek.SeverityLevel][]detek.Report)

	for _, r := range r.Reports {
		reportMap[r.Level] = append(reportMap[r.Level], r)
	}

	type LevelSet struct {
		Level   detek.SeverityLevel
		Reports []detek.Report
	}
	sortedReport := []LevelSet{}
	for _, level := range []detek.SeverityLevel{
		detek.Fatal, detek.Error, detek.Warn, detek.Normal, detek.Unknown,
	} {
		if v, ok := reportMap[level]; ok {
			sortedReport = append(sortedReport, LevelSet{
				Level:   level,
				Reports: v,
			})
		}
	}

	data := struct {
		StartedAt  time.Time
		FinishedAt time.Time
		Reports    []LevelSet
		TotalCount int
	}{
		StartedAt:  r.StartedAt,
		FinishedAt: r.FinishedAt,
		Reports:    sortedReport,
		TotalCount: len(r.Reports),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(fmt.Errorf("this is bug: %w", err))
	}
	return buf.String()
}

const _HTMLReportTemplate = _HTMLReportTemplateCSS + `
<body>
	<div class="main">
		<h1>detek Cluster Report</h1>
		<!--METADATA-->
		<div class="container">			
			<b>started at:</b>
			<p>{{.StartedAt.Format "Jan 02, 2006 15:04:05 UTC"}}</p>
			<b>finished at:</b>
			<p>{{.FinishedAt.Format "Jan 02, 2006 15:04:05 UTC"}}</p>
			<b>reports:</b>
			<div>
				{{range $_, $val := .Reports}}<p>{{$val.Level}}: {{len $val.Reports}}</p>{{end}}
				<p>Total: {{.TotalCount}}</p>
			</div>
		</div>
		<div>
			<!--REPORTS-->
			{{range $_, $val := .Reports}}
			{{$reports := $val.Reports}}
			{{$key := $val.Level}}
			<h2>Level: {{$key}}</h2>
			<div class="reports">
				{{range $_, $report := $reports}}
				<div>
					<!--A SINGLE REPORT-->
					<input id="{{$report.ID}}" type="checkbox" class="hide" />
					<label for="{{$report.ID}}" class="reportHeader">
						<h3>{{$report.ID}}</h3>
						<p>{{$report.Description}}</p>
					</label>
					<div class="reportBody">
						<div class="tags">
							{{range $_, $label := $report.Labels}}<label>{{$label}}</label>{{end}}
						</div>
						<div>
							<h2>Current State</h2>
							<p>{{$report.CurrentState.Explanation}}</p>
							{{if $report.CurrentState.Solution}}
							<h2>Solution</h2>
							<p>{{$report.CurrentState.Solution}}</p>
							{{end}}
							{{if $report.ReportSpec.Problem}}
							<h2>{{$report.ReportSpec.Problem.Description}}</h2>
							<div class="jsonWrap">
							<pre>{{marshal $report.ReportSpec.Problem.Data}}</pre>
							</div>
							{{end}}
							{{if $report.ReportSpec.Attachment}}
							<h2>More Info</h2>
							{{range $_, $data := $report.ReportSpec.Attachment}}
							<h3>{{$data.Description}}</h3>
							<div class="jsonWrap">
							<pre>{{marshal $data.Data}}</pre>
							</div>
							{{end}}
							{{end}}						
							<pre></pre>
						</div>
					</div>
				</div>
				{{end}}
			</div>
			{{end}}			
		</div>
	</div>
</body>
`

const _HTMLReportTemplateCSS = `<head><style>
	body {
		font-family: system-ui;
	}

	h1 {
		margin: 5px 0;
	}

	pre {
		overflow: auto;
		max-height: 250px;
	}

	.jsonWrap {
		background: black;
		color: white;
		padding: 10px;
		border-radius: 10px;
	}


	.tags {
		height: 40px;
		display: flex;
	}

	.tags>label {
		background-color: #999999;
		color: #ffffff;
		padding: 4px 5px;
		border-radius: 5px;
		margin: auto 5px auto 0;
	}

	.hide {
		display: none;
	}

	.main {
		margin: auto;
		max-width: 768px;
		align-items: center;
	}

	.container {
		display: inline-grid;
		grid-template-columns: 100px 200px;
		align-items: center;
		border: solid 1px;
		padding: 0 10px;
		border-radius: 5px;
	}

	.reports {
		border-radius: 5px;
		box-shadow: rgb(0 0 0 / 16%) 0px 0px 4px 1px;
	}

	.reports>div:not(:last-child) {
		border-bottom: solid #EDF1FD;
	}
	.reports>div:last-child>label {
		border-radius: 0 0 5px 5px;
	}
	.reports>div:first-child>label {
		border-radius: 5px 5px 0 0;
	}
	
	.reportHeader {
		display: flex;
		align-items: center;
		background-color: white;
		cursor: pointer;
		padding: 20px;		
	}

	.reportHeader:hover {
		background-color: ghostwhite !important;
	}

	input:checked~.reportHeader {
		background-color: #e9ecef;
	}

	.reportHeader>* {
		margin: 0;
	}

	.reportHeader>h3 {
		margin-right: 10px;
	}

	.reportBody {
		display: none;
		background-color: #fcfcfc;
		padding: 10px;
	}

	input:checked~.reportBody {
		display: block;
	}
</style></head>`
