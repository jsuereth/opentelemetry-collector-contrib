// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"
)

var StringType = NewPrimitiveType("string")

type stringVal string

// Type implements Val.
func (s stringVal) Type() Type {
	return StringType
}

// Cast strings to valid other types.
func (s stringVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	// TODO - boolean coercion.
	case reflect.String:
		return reflect.ValueOf(s).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		if typeDesc.Elem().Kind() == reflect.String {
			return &s, nil
		}
	case reflect.Interface:
		sv := s.Value()
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

func (s stringVal) Value() any {
	return (string)(s)
}

func NewStringVal(v string) Val {
	return (stringVal)(v)
}
