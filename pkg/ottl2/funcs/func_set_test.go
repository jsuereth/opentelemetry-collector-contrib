// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func Test_set(t *testing.T) {
	set := NewSetFunc()

	tests := []struct {
		name     string
		target   types.Var
		value    types.Val
		expected types.Val
	}{
		{
			name:     "set name",
			target:   types.NewPvalVar(pcommon.NewValueStr("original")),
			value:    types.NewPvalVar(pcommon.NewValueStr("new name")),
			expected: types.NewPvalVar(pcommon.NewValueStr("new name")),
		},
		{
			name:     "set nil value",
			target:   types.NewPvalVar(pcommon.NewValueStr("original")),
			value:    types.NilVal,
			expected: types.NewPvalVar(pcommon.NewValueStr("original")),
		},
		{
			name:     "set string",
			target:   types.NewPvalVar(pcommon.NewValueStr("original")),
			value:    types.NewStringVal("new name"),
			expected: types.NewPvalVar(pcommon.NewValueStr("new name")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := set.Call([]types.Val{tt.target, tt.value})
			// Null return value.
			assert.Equal(t, types.NilVal, result)
			// Target should now match result.
			assert.Equal(t, tt.expected, tt.target)
		})
	}
}
