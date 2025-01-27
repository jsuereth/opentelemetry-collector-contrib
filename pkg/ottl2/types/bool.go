// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect" // Stores literal boolean values
)

type boolVal struct {
	value bool
}

// How we coeerce between known types in OTLP.
func (b *boolVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(b).Convert(typeDesc).Interface(), nil
	case reflect.Ptr:
		switch typeDesc {
		default:
			if typeDesc.Elem().Kind() == reflect.Bool {
				p := bool(b.value)
				return &p, nil
			}
		}
	case reflect.Interface:
		bv := b.Value()
		if reflect.TypeOf(bv).Implements(typeDesc) {
			return bv, nil
		}
		if reflect.TypeOf(b).Implements(typeDesc) {
			return b, nil
		}
	}
	return nil, fmt.Errorf("type conversion error from bool to '%v'", typeDesc)
}

func (b *boolVal) Value() any {
	return (bool)(b.value)
}

func NewBoolVal(v bool) Val {
	return &boolVal{v}
}

var (
	TrueVal  = NewBoolVal(true)
	FalseVal = NewBoolVal(false)
)
