// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"errors"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - define a special pccomon.Map type
// We want this to denote it can be keyed
var PmapType = types.NewPrimitiveType("pcommon.Map")

type pmapVal pcommon.Map

func (m pmapVal) Type() types.Type {
	return PmapType
}

func (m pmapVal) ConvertTo(t types.Type) (any, error) {
	switch t {
	case PmapType:
		return m.Value(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion from 'pcommon.Map' to %s", t.Name())
}

func (m pmapVal) Value() any {
	return pcommon.Map(m)
}

// PmapVal is KeyAccessible
func (m pmapVal) GetKey(key string) types.Val {
	// Note: pcommon.Map and pcommon.Value are references that mutate the underlying value.
	// so we don't need to use a getter/setter appraoch.
	v, ok := pcommon.Map(m).Get(key)
	if !ok {
		return newEmptyMapKeyVar(pcommon.Map(m), key)
	}
	return NewPvalVar(v)
}

// SetValue implements Var.
func (m pmapVal) SetValue(o types.Val) error {
	other, err := o.ConvertTo(PmapType)
	if err != nil {
		return err
	}
	other.(pcommon.Map).MoveTo(pcommon.Map(m))
	return nil
}

type emptyMapKeyVar struct {
	m   pcommon.Map
	key string
}

// ConvertTo implements Var.
func (e *emptyMapKeyVar) ConvertTo(t types.Type) (any, error) {
	value, ok := e.m.Get(e.key)
	if ok {
		return NewPvalVar(value).ConvertTo(t)
	}
	return nil, errors.New("cannot convert missing pcommon.Map key")
}

// SetValue implements Var.
func (e *emptyMapKeyVar) SetValue(v types.Val) error {
	// TODO - check all supported types and write one.
	switch v.Type() {
	case BoolType:
		value, err := v.ConvertTo(BoolType)
		if err != nil {
			return err
		}
		e.m.PutBool(e.key, value.(bool))
		return err
	case IntType:
		value, err := v.ConvertTo(IntType)
		if err != nil {
			return err
		}
		e.m.PutInt(e.key, value.(int64))
		return err
	case FloatType:
		value, err := v.ConvertTo(FloatType)
		if err != nil {
			return err
		}
		e.m.PutDouble(e.key, value.(float64))
		return err
	// TODO - handle slices and embedded maps.
	case PvalType:
	}
	panic("unimplemented")
}

func (e *emptyMapKeyVar) Type() types.Type {
	return PvalType
}

func (e *emptyMapKeyVar) Value() any {
	result, ok := e.m.Get(e.key)
	if ok {
		return NewPvalVar(result).Value()
	}
	return nil
}

func newEmptyMapKeyVar(m pcommon.Map, key string) types.Var {
	return &emptyMapKeyVar{m, key}
}

func NewPmapVar(m pcommon.Map) types.Var {
	return pmapVal(m)
}
