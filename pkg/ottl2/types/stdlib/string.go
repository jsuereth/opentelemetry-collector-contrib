// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

var StringType = types.NewPrimitiveType("string")

type stringVal string

// Type implements Val.
func (s stringVal) Type() types.Type {
	return StringType
}

// Cast strings to valid other types.
func (s stringVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	return convertStringTo(string(s), typeDesc)
}

func (s stringVal) Value() any {
	return (string)(s)
}

func NewStringVal(v string) types.Val {
	return (stringVal)(v)
}

type string64Var struct {
	getter func() string
	setter func(string)
}

func (s string64Var) ConvertTo(typeDesc reflect.Type) (any, error) {
	return convertStringTo(s.getter(), typeDesc)
}

func (i string64Var) SetValue(v types.Val) error {
	// TODO - check types first...
	if value, ok := v.Value().(string); ok {
		i.setter(value)
		return nil
	}
	return fmt.Errorf("value is not a string: %v", v)
}

func (i string64Var) Type() types.Type {
	return StringType
}

func (i string64Var) Value() any {
	return i.getter()
}

func NewStringVar(getter func() string, setter func(string)) types.Var {
	return string64Var{getter, setter}
}

func convertStringTo(s string, typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	// TODO - boolean coercion.
	case reflect.String:
		return reflect.ValueOf(s).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		if typeDesc.Elem().Kind() == reflect.String {
			return &s, nil
		}
	case reflect.Interface:
		sv := s
		if reflect.TypeOf(sv).Implements(typeDesc) {
			return sv, nil
		}
		if reflect.TypeOf(s).Implements(typeDesc) {
			return s, nil
		}
	}
	return nil, fmt.Errorf(
		"unsupported native conversion from string to '%v'", typeDesc)
}
