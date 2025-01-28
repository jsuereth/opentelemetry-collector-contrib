// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"
)

type float64Val float64

// ConvertTo implements Val.
func (f float64Val) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	case reflect.Float64:
		v := float64(f)
		return reflect.ValueOf(v).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		switch typeDesc.Elem().Kind() {
		case reflect.Float64:
			v := float64(f)
			p := reflect.New(typeDesc.Elem())
			p.Elem().Set(reflect.ValueOf(v).Convert(typeDesc.Elem()))
			return p.Interface(), nil
		}
	case reflect.Interface:
		dv := f.Value()
		if reflect.TypeOf(dv).Implements(typeDesc) {
			return dv, nil
		}
		if reflect.TypeOf(f).Implements(typeDesc) {
			return f, nil
		}
	}
	return nil, fmt.Errorf("type conversion error from Double to '%v'", typeDesc)
}

func (f float64Val) Value() any {
	return (float64)(f)
}

var float64Type = reflect.TypeOf((float64)(0))

// Flaots are addable
func (f float64Val) Add(o Val) Val {
	rhs, err := o.ConvertTo(float64Type)
	if err != nil {
		return NewErrorVal(err)
	}
	return (float64Val)((float64)(f) + rhs.(float64))
}

func NewFloatVal(v float64) Val {
	h := (float64Val)(v)
	return &h
}
