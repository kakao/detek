package renderer

import (
	"fmt"

	"github.com/kakao/detek/pkg/detek"
)

type RenderOpts struct {
	Table struct {
		MaxWidth int
	}
	HTML struct{}
	JSON struct {
		Pretty bool
	}
	// YAML struct{}
}

type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatHTML  Format = "html"
)

func (f *Format) IsValid() error {
	if f == nil {
		return fmt.Errorf("this is nil")
	}
	for _, t := range []Format{FormatJSON, FormatTable, FormatHTML} {
		if *f == t {
			return nil
		}
	}
	return fmt.Errorf("%q is not supported format", string(*f))
}

func RenderReports(list *detek.ReportList, format Format, opts RenderOpts) string {
	switch format {
	case FormatHTML:
		return RenderHTMLReports(*list)
	case FormatJSON:
		return RenderJSONReports(*list, opts.JSON.Pretty)
	case FormatTable:
		return RenderTableReports(*list, opts.Table.MaxWidth)
	default:
		return "unsupported format"
	}
}
