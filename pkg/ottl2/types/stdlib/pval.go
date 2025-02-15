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
	return convertPval(pcommon.Value(p), t)
}

func (p pvalVal) Type() types.Type {
	return PvalType
}

func (p pvalVal) Value() any {
	return valueOf(pcommon.Value(p))
}

// pvalVal is a Var
func (p pvalVal) SetValue(v types.Val) error {
	return setPValue(pcommon.Value(p), v)
}

// GetIndex implements traits.Indexable
func (s pvalVal) GetIndex(index int64) types.Val {
	// TODO - deal with all slice types
	idx := int(index)
	if idx < pcommon.Value(s).Slice().Len() {
		return NewPvalVar(pcommon.Value(s).Slice().At(int(index)))
	} else {
		return NewLazyPval(func() pcommon.Value {
			for pcommon.Value(s).Slice().Len() < idx {
				pcommon.Value(s).Slice().AppendEmpty()
			}
			return pcommon.Value(s).Slice().At(idx)
		})
	}
}

// GetIndex implements traits.KeyAccessable
func (s pvalVal) GetKey(key string) types.Val {
	// TODO - better error message.
	if r, ok := pcommon.Value(s).Map().Get(key); ok {
		return NewPvalVar(r)
	}
	// Create lazy key.
	return NewLazyPval(func() pcommon.Value {
		return pcommon.Value(s).Map().PutEmpty(key)
	})
}

func NewPvalVar(v pcommon.Value) types.Var {
	return pvalVal(v)
}

type pvalVar func() pcommon.Value

// ConvertTo implements types.Var.
func (p pvalVar) ConvertTo(t types.Type) (any, error) {
	return convertPval((func() pcommon.Value)(p)(), t)
}

// SetValue implements types.Var.
func (p pvalVar) SetValue(v types.Val) error {
	return setPValue((func() pcommon.Value)(p)(), v)
}

func (s pvalVar) GetIndex(index int64) types.Val {
	return NewLazyPval(func() pcommon.Value {
		v := (func() pcommon.Value)(s)()
		// TODO - deal with bytes and other slices...
		if v.Type() != pcommon.ValueTypeSlice {
			// Should we only do this on empty values?
			v.SetEmptySlice()
		}
		idx := int(index)
		for v.Slice().Len() <= idx {
			v.Slice().AppendEmpty()
		}
		return v.Slice().At(idx)
	})
}

func (s pvalVar) GetKey(key string) types.Val {
	return NewLazyPval(func() pcommon.Value {
		v := (func() pcommon.Value)(s)()
		if v.Type() != pcommon.ValueTypeMap {
			// Should we only do this on empty values?
			v.SetEmptyMap()
		}
		if r, ok := v.Map().Get(key); ok {
			return r
		}
		return v.Map().PutEmpty(key)
	})
}

// Type implements types.Var.
func (p pvalVar) Type() types.Type {
	return PvalType
}

// Value implements types.Var.
func (p pvalVar) Value() any {
	return valueOf((func() pcommon.Value)(p)())
}

// Constructs a new Pval type that will generate its "slot" on demand by calling the method.
func NewLazyPval(v func() pcommon.Value) types.Var {
	return pvalVar(v)
}

func convertPval(p pcommon.Value, t types.Type) (any, error) {
	switch t {
	case BoolType:
		if p.Type() == pcommon.ValueTypeBool {
			return p.Bool(), nil
		} else {
			return nil, fmt.Errorf("%v is not a boolean", p)
		}
	case IntType:
		if p.Type() == pcommon.ValueTypeInt {
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
	case SliceType:
		if pcommon.Value(p).Type() == pcommon.ValueTypeSlice {
			return pcommon.Value(p).Slice(), nil
		} else {
			return nil, fmt.Errorf("%v is not a slice", p)
		}
	case PmapType:
		if p.Type() == pcommon.ValueTypeMap {
			return p.Map(), nil
		} else {
			return nil, fmt.Errorf("%v is not a map", p)
		}
	}
	// TODO - other pvalue types.
	return nil, fmt.Errorf("unknown type for pcommon.Value: %v", t.Name())
}

func valueOf(p pcommon.Value) any {
	switch p.Type() {
	case pcommon.ValueTypeBool:
		return p.Bool()
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

func setPValue(p pcommon.Value, v types.Val) error {
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
	case SliceType:
		pcommon.Value(p).Slice().FromRaw(v.Value().(pcommon.Slice).AsRaw())
		return nil
	}
	// TODO - other types...
	panic(fmt.Sprintf("unimplemented conversion %v to pcommon.Value", v.Type()))
}
