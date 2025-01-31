// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"reflect"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
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
			value:    types.NewIntVal(0),
			expected: true,
		},
		{
			name:     "ValueTypeInt",
			value:    types.NewPvalVar(pcommon.NewValueInt(0)),
			expected: true,
		},
		{
			name:     "float64",
			value:    types.NewFloatVal(2.7),
			expected: false,
		},
		{
			name:     "ValueTypeString",
			value:    types.NewPvalVar(pcommon.NewValueStr("a string")),
			expected: false,
		},
		{
			name:     "not Int",
			value:    types.NewStringVal("string"),
			expected: false,
		},
		{
			name:     "string number",
			value:    types.NewStringVal("0"),
			expected: false,
		},
		{
			name:     "ValueTypeSlice",
			value:    types.NewPvalVar(pcommon.NewValueSlice()),
			expected: false,
		},
		{
			name:     "nil",
			value:    types.NilVal,
			expected: false,
		},
	}
	isInt := NewIsIntFunc()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.CallFunction(isInt, []types.Val{
				tt.value,
			}, map[string]types.Val{})
			v, err := result.ConvertTo(reflect.TypeFor[bool]())
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, v.(bool))
		})
	}
}
