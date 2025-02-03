// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

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

func (i int64Val) ConvertTo(t types.Type) (any, error) {
	return convertIntTo(int64(i), t)
}

// Integers are addable
func (i int64Val) Add(o types.Val) types.Val {
	rhs, err := o.ConvertTo(IntType)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) + rhs.(int64))
}

// Integers are subtractable
func (i int64Val) Subtract(o types.Val) types.Val {
	rhs, err := o.ConvertTo(IntType)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) - rhs.(int64))
}

// Integers are comparable
func (i int64Val) Equals(o types.Val) bool {
	rhs, err := o.ConvertTo(IntType)
	if err != nil {
		return false
	}
	return ((int64)(i) == rhs.(int64))
}

func (i int64Val) LessThan(o types.Val) bool {
	rhs, err := o.ConvertTo(IntType)
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
func (i int64Var) ConvertTo(t types.Type) (any, error) {
	return convertIntTo(i.getter(), t)
}

// SetValue implements Var.
func (i int64Var) SetValue(v types.Val) error {
	// TODO - check types first...
	if value, ok := v.Value().(int64); ok {
		i.setter(value)
		return nil
	}
	return fmt.Errorf("cannot set integer from value: %v", v)
}

// Type implements Var.
func (i int64Var) Type() types.Type {
	return IntType
}

// Value implements Var.
func (i int64Var) Value() any {
	return i.getter()
}

func (i int64Var) Add(o types.Val) types.Val {
	// TODO - just use Type() directly on o and .Value().(int64)
	rhs, err := o.ConvertTo(IntType)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i.getter()) + rhs.(int64))
}

func NewIntVar(getter func() int64, setter func(int64)) types.Var {
	return int64Var{getter, setter}
}

func convertIntTo(i int64, t types.Type) (any, error) {
	switch t {
	case IntType:
		return int64(i), nil
	}
	return nil, fmt.Errorf("unsupported type conversion from 'int' to %v", t)
}
