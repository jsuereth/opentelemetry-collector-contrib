// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/pdata/pcommon" // TODO - define a special type for this as we do lots of conversion
)

// into/out of this in a real type system.
// We need to treat this as a special "top type".
var PvalType = NewPrimitiveType("pcommon.Value")

type pvalVal pcommon.Value

// ConvertTo implements Val.
func (p pvalVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	// TODO - appropriately case this to zero value for primitives?
	if typeDesc == nil && pcommon.Value(p).Type() == pcommon.ValueTypeEmpty {
		return nil, nil
	}
	switch typeDesc.Kind() {
	case reflect.Int64:
		if pcommon.Value(p).Type() == pcommon.ValueTypeInt {
			return pcommon.Value(p).Int(), nil
		} else {
			return nil, fmt.Errorf("%v is not an integer", p)
		}
	case reflect.Bool:
		if pcommon.Value(p).Type() == pcommon.ValueTypeBool {
			return pcommon.Value(p).Bool(), nil
		} else {
			return nil, fmt.Errorf("%v is not a boolean", p)
		}
	case reflect.Float64:
		if pcommon.Value(p).Type() == pcommon.ValueTypeDouble {
			return pcommon.Value(p).Double(), nil
		} else {
			return nil, fmt.Errorf("%v is not a double", p)
		}
	case reflect.String:
		if pcommon.Value(p).Type() == pcommon.ValueTypeStr {
			return pcommon.Value(p).Str(), nil
		} else {
			return nil, fmt.Errorf("%v is not a string", p)
		}
	}
	// TODO - other pvalue types.
	return nil, fmt.Errorf("unknown type for pcommon.Value: %v", typeDesc)
}

func (p pvalVal) Type() Type {
	return PvalType
}

func (p pvalVal) Value() any {
	// TODO - we should probably erase to a 'primitive' value here if we can.
	return pcommon.Value(p)
}

// pvalVal is a Var
func (p pvalVal) SetValue(v Val) error {
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
	}
	// TODO - other types...
	panic(fmt.Sprintf("unimplemented conversion %v to pcommon.Value", v.Type()))
}

func NewPvalVar(v pcommon.Value) Var {
	return pvalVal(v)
}
