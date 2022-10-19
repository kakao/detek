package renderer

import (
	"encoding/json"
	"fmt"

	"github.com/kakao/detek/pkg/detek"
)

func RenderJSONReports(r detek.ReportList, pretty bool) string {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(r, "", "  ")
	} else {
		b, err = json.Marshal(r)
	}
	if err != nil {
		panic(fmt.Errorf("this is a bug: %w", err))
	}
	return string(b)
}
