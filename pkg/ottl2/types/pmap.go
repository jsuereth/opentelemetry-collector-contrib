// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"errors"
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - define a special pccomon.Map type
// We want this to denote it can be keyed
var PmapType = NewPrimitiveType("pcommon.Map")

type pmapVal pcommon.Map

func (m pmapVal) Type() Type {
	return PmapType
}

func (m pmapVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	case reflect.TypeFor[pcommon.Map]().Kind():
		return reflect.ValueOf(m).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		switch typeDesc.Elem().Kind() {
		case reflect.TypeFor[pcommon.Map]().Kind():
			v := pcommon.Map(m)
			p := reflect.New(typeDesc.Elem())
			p.Elem().Set(reflect.ValueOf(v).Convert(typeDesc.Elem()))
			return p.Interface(), nil
		}
	case reflect.Interface:
		iv := m.Value()
		if reflect.TypeOf(iv).Implements(typeDesc) {
			return iv, nil
		}
		if reflect.TypeOf(m).Implements(typeDesc) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("unsupported type conversion from 'pcommon.Map' to %v", typeDesc)
}

func (m pmapVal) Value() any {
	return pcommon.Map(m)
}

// PmapVal is KeyAccessible
func (m pmapVal) GetKey(key string) Val {
	// Note: pcommon.Map and pcommon.Value are references that mutate the underlying value.
	// so we don't need to use a getter/setter appraoch.
	v, ok := pcommon.Map(m).Get(key)
	if !ok {
		return newEmptyMapKeyVar(pcommon.Map(m), key)
	}
	return NewPvalVar(v)
}

// SetValue implements Var.
func (m pmapVal) SetValue(o Val) error {
	other, err := o.ConvertTo(reflect.TypeFor[pcommon.Map]())
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
func (e *emptyMapKeyVar) ConvertTo(typeDesc reflect.Type) (any, error) {
	value, ok := e.m.Get(e.key)
	if ok {
		return NewPvalVar(value).ConvertTo(typeDesc)
	}
	return nil, errors.New("cannot convert missing pcommon.Map key")
}

// SetValue implements Var.
func (e *emptyMapKeyVar) SetValue(v Val) error {
	// TODO - check all supported types and write one.
	switch v.Type() {
	case BoolType:
		value, err := v.ConvertTo(reflect.TypeFor[bool]())
		if err != nil {
			return err
		}
		e.m.PutBool(e.key, value.(bool))
		return err
	case IntType:
		value, err := v.ConvertTo(reflect.TypeFor[int64]())
		if err != nil {
			return err
		}
		e.m.PutInt(e.key, value.(int64))
		return err
	case FloatType:
		value, err := v.ConvertTo(reflect.TypeFor[float64]())
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

func (e *emptyMapKeyVar) Type() Type {
	return PvalType
}

func (e *emptyMapKeyVar) Value() any {
	result, ok := e.m.Get(e.key)
	if ok {
		return NewPvalVar(result).Value()
	}
	return nil
}

func newEmptyMapKeyVar(m pcommon.Map, key string) Var {
	return &emptyMapKeyVar{m, key}
}

func NewPmapVar(m pcommon.Map) Var {
	return pmapVal(m)
}
