// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func Test_IsInt(t *testing.T) {
	tests := []struct {
		name     string
		value    types.Val
		expected bool
	}{
		{
			name:     "int",
			value:    stdlib.NewIntVal(0),
			expected: true,
		},
		{
			name:     "ValueTypeInt",
			value:    stdlib.NewPvalVar(pcommon.NewValueInt(0)),
			expected: true,
		},
		{
			name:     "float64",
			value:    stdlib.NewFloatVal(2.7),
			expected: false,
		},
		{
			name:     "ValueTypeString",
			value:    stdlib.NewPvalVar(pcommon.NewValueStr("a string")),
			expected: false,
		},
		{
			name:     "not Int",
			value:    stdlib.NewStringVal("string"),
			expected: false,
		},
		{
			name:     "string number",
			value:    stdlib.NewStringVal("0"),
			expected: false,
		},
		{
			name:     "ValueTypeSlice",
			value:    stdlib.NewPvalVar(pcommon.NewValueSlice()),
			expected: false,
		},
		{
			name:     "nil",
			value:    stdlib.NilVal,
			expected: false,
		},
	}
	isInt := NewIsIntFunc()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInt.Call([]types.Val{
				tt.value,
			})
			v, err := result.ConvertTo(stdlib.BoolType)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, v.(bool))
		})
	}
}
