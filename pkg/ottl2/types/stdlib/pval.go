// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon" // TODO - define a special type for this as we do lots of conversion
)

// into/out of this in a real type system.
// We need to treat this as a special "top type".
var PvalType = types.NewPrimitiveType("pcommon.Value")

type pvalVal pcommon.Value

// ConvertTo implements Val.
func (p pvalVal) ConvertTo(t types.Type) (any, error) {
	switch t {
	case BoolType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeBool {
			return pcommon.Value(p).Bool(), nil
		} else {
			return nil, fmt.Errorf("%v is not a boolean", p)
		}
	case IntType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeInt {
			return pcommon.Value(p).Int(), nil
		} else {
			return nil, fmt.Errorf("%v is not an integer", p)
		}
	case FloatType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeDouble {
			return pcommon.Value(p).Double(), nil
		} else {
			return nil, fmt.Errorf("%v is not a double", p)
		}
	case StringType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeStr {
			return pcommon.Value(p).Str(), nil
		} else {
			return nil, fmt.Errorf("%v is not a string", p)
		}
	case NilType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeEmpty {
			return nil, nil
		} else {
			return nil, fmt.Errorf("%v is not empty", p)
		}
	}
	// TODO - other pvalue types.
	return nil, fmt.Errorf("unknown type for pcommon.Value: %v", t.Name())
}

func (p pvalVal) Type() types.Type {
	return PvalType
}

func (p pvalVal) Value() any {
	switch pcommon.Value(p).Type() {
	case pcommon.ValueTypeBool:
		return pcommon.Value(p).Bool()
	case pcommon.ValueTypeBytes:
		return pcommon.Value(p).Bytes()
	case pcommon.ValueTypeDouble:
		return pcommon.Value(p).Double()
	case pcommon.ValueTypeEmpty:
		return nil
	case pcommon.ValueTypeInt:
		return pcommon.Value(p).Int()
	case pcommon.ValueTypeMap:
		return pcommon.Value(p).Map()
	case pcommon.ValueTypeSlice:
		// TODO - refine slice into specific array?
		return pcommon.Value(p).Slice()
	case pcommon.ValueTypeStr:
		return pcommon.Value(p).AsString()
	default:
		panic("unexpected pcommon.ValueType")
	}
}

// pvalVal is a Var
func (p pvalVal) SetValue(v types.Val) error {
	switch v.Type() {
	case PvalType:
		other := pcommon.Value(v.(pvalVal))
		other.CopyTo(pcommon.Value(p))
		return nil
	case BoolType:
		pcommon.Value(p).SetBool(v.Value().(bool))
		return nil
	case StringType:
		pcommon.Value(p).SetStr(v.Value().(string))
		return nil
	case FloatType:
		pcommon.Value(p).SetDouble(v.Value().(float64))
		return nil
	case IntType:
		pcommon.Value(p).SetInt(v.Value().(int64))
		return nil
	case NilType:
		// Do nothing.
		return nil
	case ByteSliceType:
		pcommon.Value(p).SetEmptyBytes().FromRaw(v.Value().([]byte))
		return nil
	}
	// TODO - other types...
	panic(fmt.Sprintf("unimplemented conversion %v to pcommon.Value", v.Type()))
}

func NewPvalVar(v pcommon.Value) types.Var {
	return pvalVal(v)
}
