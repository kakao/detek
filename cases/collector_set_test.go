package cases_test

import (
	"fmt"
	"testing"

	"github.com/kakao/detek/cases"
	"github.com/kakao/detek/pkg/detek"
	"github.com/stretchr/testify/assert"
)

func TestValidatingCollectorMeta(t *testing.T) {
	IDMap := make(map[string]bool)
	for _, set := range cases.CollectorSet {
		for _, c := range set(map[string]string{}) {
			meta := c.GetMeta()
			assert.NotEmpty(t, meta.ID, fmt.Sprintf("id for %q is not set", detek.TypeOf(c).String()))
			assert.NotEmpty(t, meta.Description, fmt.Sprintf("description for %q is not set", meta.ID))
			assert.NotEmpty(t, meta.Producing, fmt.Sprintf("collector %q produce nothing", meta.ID))
			if _, ok := IDMap[meta.ID]; ok {
				assert.Fail(t, "duplicated ID detected", meta.ID)
			}
			IDMap[meta.ID] = true
		}
	}
}
