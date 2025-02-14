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

func Test_Append(t *testing.T) {
	testCases := []struct {
		Name   string
		Target types.Var
		Value  any
		Values []any
		Want   func(pcommon.Slice)
	}{
		{
			Name:   "Slice: standard []string target - empty",
			Target: stdlib.NewSliceVar(pcommon.NewSlice()),
			Values: []any{"a", "b"},
			Want: func(expectedValue pcommon.Slice) {
				expectedValue.AppendEmpty().SetStr("a")
				expectedValue.AppendEmpty().SetStr("b")
			},
		},
		{
			Name: "Single: standard []string target - value",
			Target: stdlib.NewSliceVar(func() pcommon.Slice {
				r := pcommon.NewSlice()
				r.FromRaw([]any{"5", "6"})
				return r
			}()),
			Value: "a",
			Want: func(expectedValue pcommon.Slice) {
				expectedValue.AppendEmpty().SetStr("5")
				expectedValue.AppendEmpty().SetStr("6")
				expectedValue.AppendEmpty().SetStr("a")
			},
		},
		{
			Name: "Slice: standard []string target - values",
			Target: stdlib.NewSliceVar(func() pcommon.Slice {
				r := pcommon.NewSlice()
				r.FromRaw([]any{"5", "6"})
				return r
			}()),
			Values: []any{"a", "b"},
			Want: func(expectedValue pcommon.Slice) {
				expectedValue.AppendEmpty().SetStr("5")
				expectedValue.AppendEmpty().SetStr("6")
				expectedValue.AppendEmpty().SetStr("a")
				expectedValue.AppendEmpty().SetStr("b")
			},
		},
	}

	f := NewAppendFunc()
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			args := make([]types.Val, 3)
			args[0] = tc.Target
			if tc.Value != nil {
				args[1] = stdlib.NewPvalVar(valueFromAny(tc.Value))
			} else {
				args[1] = f.DefaultArgs()["Value"]
			}
			if tc.Values != nil {
				a := pcommon.NewSlice()
				a.FromRaw(tc.Values)
				args[2] = stdlib.NewSliceVar(a)
			} else {
				args[2] = f.DefaultArgs()["Values"]
			}
			result := f.Call(args)
			if result.Type() == stdlib.ErrorType {
				assert.Fail(t, "append failed: %v", result.Value())
			}
			expectedSlice := pcommon.NewSlice()
			tc.Want(expectedSlice)
			assert.EqualValues(t, expectedSlice, tc.Target.Value())
		})
	}
}

func valueFromAny(v any) pcommon.Value {
	// TODO - Handle all the types we may be appending
	switch valueType := v.(type) {
	case int64:
		return pcommon.NewValueInt(valueType)
	case string:
		return pcommon.NewValueStr(valueType)
	default:
		return pcommon.NewValueEmpty()
	}
}
