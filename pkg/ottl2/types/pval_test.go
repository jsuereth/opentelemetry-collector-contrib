// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestPval_setValue(t *testing.T) {
	tests := []struct {
		name     string
		input    Val
		expected pcommon.Value
	}{
		{
			name:     "PVal accepts booleans",
			input:    NewBoolVal(true),
			expected: pcommon.NewValueBool(true),
		},
		{
			name:     "PVal accepts ints",
			input:    NewIntVal(24),
			expected: pcommon.NewValueInt(24),
		},
		{
			name:     "PVal accepts floats",
			input:    NewFloatVal(4.34),
			expected: pcommon.NewValueDouble(4.34),
		},
		{
			name:     "PVal accepts strings",
			input:    NewStringVal("my test"),
			expected: pcommon.NewValueStr("my test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pval := NewPvalVar(pcommon.NewValueEmpty())
			err := pval.SetValue(tt.input)
			assert.Nil(t, err)
			result := pval.Value()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPval_ConvertTo(t *testing.T) {
	tests := []struct {
		name     string
		input    pcommon.Value
		typeCast reflect.Type
		expected any
	}{
		{
			name:     "converts to bool",
			input:    pcommon.NewValueBool(false),
			typeCast: reflect.TypeFor[bool](),
			expected: false,
		},
		{
			name:     "converts to int",
			input:    pcommon.NewValueInt(42),
			typeCast: reflect.TypeFor[int64](),
			expected: int64(42),
		},
		{
			name:     "converts to double",
			input:    pcommon.NewValueDouble(34.5),
			typeCast: reflect.TypeFor[float64](),
			expected: float64(34.5),
		},
		{
			name:     "converts to string",
			input:    pcommon.NewValueStr("test val"),
			typeCast: reflect.TypeFor[string](),
			expected: "test val",
		},
		{
			name:     "empty return nil",
			input:    pcommon.NewValueEmpty(),
			typeCast: reflect.TypeOf(nil),
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pval := NewPvalVar(tt.input)
			result, err := pval.ConvertTo(tt.typeCast)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
