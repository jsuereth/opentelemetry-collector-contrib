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

func Test_IsMap(t *testing.T) {
	tests := []struct {
		name     string
		value    runtime.Val
		expected bool
	}{
		{
			name:     "map",
			value:    stdlib.NewPmapVar(pcommon.NewMap()),
			expected: true,
		},
		{
			name:     "ValueTypeMap",
			value:    stdlib.NewPvalVar(pcommon.NewValueMap()),
			expected: true,
		},
		{
			name:     "not map",
			value:    stdlib.NewStringVal("not a map"),
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
	isMap := NewIsMapFunc()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMap.Call([]runtime.Val{
				tt.value,
			})
			v, err := result.ConvertTo(stdlib.BoolType)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, v.(bool))
		})
	}
}
