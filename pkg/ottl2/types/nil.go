// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"
)

// TODO - Nil should be a special type.
var NilType = NewPrimitiveType("nil")

type nilVal struct{}

// Type implements Val.
func (n *nilVal) Type() Type {
	return NilType
}

// ConvertTo implements Val.
func (n *nilVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	switch typeDesc.Kind() {
	case reflect.Bool:
		return false, nil
	case reflect.Ptr:
		return nil, nil
	case reflect.Interface:
		return nil, nil
	}
	// If the type conversion isn't supported return an error.
	return nil, fmt.Errorf("type conversion error from 'nil' to '%v'", typeDesc)
}

func (n *nilVal) Value() any {
	return nil
}

var NilVal Val = &nilVal{}
