// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"
)

type int64Val int64

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
func (i int64Val) Add(o Val) Val {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) + rhs.(int64))
}

// Integers are subtractable
func (i int64Val) Subtract(o Val) Val {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return NewErrorVal(err)
	}
	return (int64Val)((int64)(i) - rhs.(int64))
}

// Integers are comparable
func (i int64Val) Equals(o Val) bool {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return false
	}
	return ((int64)(i) == rhs.(int64))
}

func (i int64Val) LessThan(o Val) bool {
	rhs, err := o.ConvertTo(int64Type)
	if err != nil {
		return false
	}
	return ((int64)(i) < rhs.(int64))
}

func NewIntVal(v int64) Val {
	return (int64Val)(v)
}
