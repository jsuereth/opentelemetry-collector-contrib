// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

var IntType = types.NewPrimitiveType("int")

type int64Val int64

// Type implements Val.
func (i int64Val) Type() types.Type {
	return IntType
}

func (i int64Val) Value() any {
	return (int64)(i)
}

func (i int64Val) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	case reflect.Int64:
		return reflect.ValueOf(i).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		switch typeDesc.Elem().Kind() {
		case reflect.Int64:
			v := int64(i)
			p := reflect.New(typeDesc.Elem())
			p.Elem().Set(reflect.ValueOf(v).Convert(typeDesc.Elem()))
			return p.Interface(), nil
		}
	case reflect.Interface:
		iv := i.Value()
		if reflect.TypeOf(iv).Implements(typeDesc) {
			return iv, nil
		}
		if reflect.TypeOf(i).Implements(typeDesc) {
			return i, nil
		}
	}
	return nil, fmt.Errorf("unsupported type conversion from 'int' to %v", typeDesc)
}

var int64Type = reflect.TypeOf(int64(0))

// Integers are addable
func (i int64Val) Add(o types.Val) types.Val {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) + rhs.(int64))
}

// Integers are subtractable
func (i int64Val) Subtract(o types.Val) types.Val {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) - rhs.(int64))
}

// Integers are comparable
func (i int64Val) Equals(o types.Val) bool {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return false
	}
	return ((int64)(i) == rhs.(int64))
}

func (i int64Val) LessThan(o types.Val) bool {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return false
	}
	return ((int64)(i) < rhs.(int64))
}

func NewIntVal(v int64) types.Val {
	return (int64Val)(v)
}

type int64Var struct {
	getter func() int64
	setter func(int64)
}

// ConvertTo implements Var.
func (i int64Var) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements Var.
func (i int64Var) SetValue(v types.Val) error {
	// TODO - check types first...
	if value, ok := v.Value().(int64); ok {
		i.setter(value)
	}
	return nil
}

// Type implements Var.
func (i int64Var) Type() types.Type {
	return int64Type
}

// Value implements Var.
func (i int64Var) Value() any {
	return i.getter()
}

func NewIntVar(getter func() int64, setter func(int64)) types.Var {
	return int64Var{getter, setter}
}
