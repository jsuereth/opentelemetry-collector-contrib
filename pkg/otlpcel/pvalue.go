// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var PValueType = NewStructType("pcommon.Value",
	[]Field[pcommon.Value]{
		{
			Name: "string",
			Type: types.StringType,
			IsSet: func(target pcommon.Value) bool {
				return target.Type() == pcommon.ValueTypeStr
			},
			GetFrom: func(target pcommon.Value) (any, error) {
				if target.Type() == pcommon.ValueTypeStr || target.Type() == pcommon.ValueTypeEmpty {
					return target.AsString(), nil
				}
				return nil, fmt.Errorf("value is not a string: %s", target.Type().String())
			},
		},
		{
			Name: "int",
			Type: types.IntType,
			IsSet: func(target pcommon.Value) bool {
				return target.Type() == pcommon.ValueTypeInt
			},
			GetFrom: func(target pcommon.Value) (any, error) {
				if target.Type() == pcommon.ValueTypeInt || target.Type() == pcommon.ValueTypeEmpty {
					return target.Int(), nil
				}
				return nil, fmt.Errorf("value is not an int: %s", target.Type().String())
			},
		},
		{
			Name: "bool",
			Type: types.BoolType,
			IsSet: func(target pcommon.Value) bool {
				return target.Type() == pcommon.ValueTypeBool
			},
			GetFrom: func(target pcommon.Value) (any, error) {
				if target.Type() == pcommon.ValueTypeBool || target.Type() == pcommon.ValueTypeEmpty {
					return target.Bool(), nil
				}
				return nil, fmt.Errorf("value is not an int: %s", target.Type().String())
			},
		},
		{
			Name: "double",
			Type: types.BoolType,
			IsSet: func(target pcommon.Value) bool {
				return target.Type() == pcommon.ValueTypeDouble
			},
			GetFrom: func(target pcommon.Value) (any, error) {
				if target.Type() == pcommon.ValueTypeDouble || target.Type() == pcommon.ValueTypeEmpty {
					return target.Double(), nil
				}
				return nil, fmt.Errorf("value is not an int: %s", target.Type().String())
			},
		},
	},
	func(m map[string]any) (pcommon.Value, error) {
		result := pcommon.NewValueMap()
		for k, v := range m {
			// TODO - put in shared util somewhere
			switch resultType := v.(type) {
			case int64:
				result.Map().PutInt(k, resultType)
			case string:
				result.Map().PutStr(k, resultType)
			case bool:
				result.Map().PutBool(k, resultType)
			case float64:
				result.Map().PutDouble(k, resultType)
				// TODO: Bytes
				// TODO: Slice
				// TODO: Map
			}
		}
		return result, nil
	},
	func(v pcommon.Value) ref.Val {
		return pvalueWrapper(v)
	},
	// TODO - add all the traits.
	traits.AdderType,
	traits.ComparerType,
	traits.DividerType,
	traits.SubtractorType,
	traits.MultiplierType,
	traits.IterableType,
	traits.ContainerType,
	traits.FieldTesterType,
	traits.FoldableType,
	traits.ListerType,
	traits.MapperType,
	traits.MultiplierType,
	traits.NegatorType,
	// TODO - Receiver type?
)

type pvalueWrapper pcommon.Value

// ConvertToNative implements ref.Val.
func (p pvalueWrapper) ConvertToNative(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// ConvertToType implements ref.Val.
func (p pvalueWrapper) ConvertToType(typeValue ref.Type) ref.Val {
	panic("unimplemented")
}

// Equal implements ref.Val.
func (p pvalueWrapper) Equal(other ref.Val) ref.Val {
	panic("unimplemented")
}

// Type implements ref.Val.
func (p pvalueWrapper) Type() ref.Type {
	return PValueType.Type()
}

// Value implements ref.Val.
func (p pvalueWrapper) Value() any {
	return pcommon.Value(p)
}
