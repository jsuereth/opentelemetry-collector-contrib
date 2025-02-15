// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func Test_set(t *testing.T) {
	set := NewSetFunc()

	tests := []struct {
		name     string
		target   runtime.Var
		value    runtime.Val
		expected runtime.Val
	}{
		{
			name:     "set name",
			target:   stdlib.NewPvalVar(pcommon.NewValueStr("original")),
			value:    stdlib.NewPvalVar(pcommon.NewValueStr("new name")),
			expected: stdlib.NewPvalVar(pcommon.NewValueStr("new name")),
		},
		{
			name:     "set nil value",
			target:   stdlib.NewPvalVar(pcommon.NewValueStr("original")),
			value:    stdlib.NilVal,
			expected: stdlib.NewPvalVar(pcommon.NewValueStr("original")),
		},
		{
			name:     "set string",
			target:   stdlib.NewPvalVar(pcommon.NewValueStr("original")),
			value:    stdlib.NewStringVal("new name"),
			expected: stdlib.NewPvalVar(pcommon.NewValueStr("new name")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := set.Call([]runtime.Val{tt.target, tt.value})
			// Null return value.
			assert.Equal(t, stdlib.NilVal, result)
			// Target should now match result.
			assert.Equal(t, tt.expected, tt.target)
		})
	}
}
