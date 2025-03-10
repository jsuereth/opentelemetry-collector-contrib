// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var PMapType *types.Type = types.MapType

type pmapWrapper pcommon.Map

// Find implements traits.Mapper.
func (p pmapWrapper) Find(key ref.Val) (ref.Val, bool) {
	field := key.Value().(string)
	m := pcommon.Map(p)
	v, ok := m.Get(field)
	// TODO - convert value to ref.Val with provider
	if !ok {
		return nil, false
	}
	// TODO
	return internalNativeToValue(v), true
}

// Get implements traits.Mapper.
func (p pmapWrapper) Get(index ref.Val) ref.Val {
	panic("unimplemented")
}

// Iterator implements traits.Mapper.
func (p pmapWrapper) Iterator() traits.Iterator {
	panic("unimplemented")
}

// Size implements traits.Mapper.
func (p pmapWrapper) Size() ref.Val {
	panic("unimplemented")
}

// Contains implements traits.Container.
func (p pmapWrapper) Contains(value ref.Val) ref.Val {
	// TODO - test values
	panic("unimplemented")
}

// ConvertToNative implements ref.Val.
func (p pmapWrapper) ConvertToNative(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// ConvertToType implements ref.Val.
func (p pmapWrapper) ConvertToType(typeValue ref.Type) ref.Val {
	panic("unimplemented")
}

// Equal implements ref.Val.
func (p pmapWrapper) Equal(other ref.Val) ref.Val {
	panic("unimplemented")
}

// Type implements ref.Val.
func (p pmapWrapper) Type() ref.Type {
	return PMapType
}

// Value implements ref.Val.
func (p pmapWrapper) Value() any {
	return pcommon.Map(p)
}

var _ ref.Val = pmapWrapper(pcommon.NewMap())
var _ traits.Mapper = pmapWrapper(pcommon.NewMap())
