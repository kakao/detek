package cases_test

import (
	"fmt"
	"testing"

	"github.com/kakao/detek/cases"
	"github.com/kakao/detek/pkg/detek"
	"github.com/stretchr/testify/assert"
)

func TestValidatingDetectorMeta(t *testing.T) {
	IDMap := make(map[string]bool)
	for _, set := range cases.DetectorSet {
		for _, d := range set(map[string]string{}) {
			meta := d.GetMeta()
			assert.NotEmpty(t, meta.ID, fmt.Sprintf("id for %q is not set", detek.TypeOf(d).String()))
			assert.NotEmpty(t, meta.Description, fmt.Sprintf("description for %q is not set", meta.ID))
			assert.NotEmpty(t, meta.IfHappened.Explanation, fmt.Sprintf("explanation for %q is not set", meta.ID))
			assert.NotEmpty(t, meta.Level, fmt.Sprintf("level for %q is not set", meta.ID))
			if _, ok := IDMap[meta.ID]; ok {
				assert.Fail(t, "duplicated ID detected", meta.ID)
			}
			IDMap[meta.ID] = true
		}
	}
}
