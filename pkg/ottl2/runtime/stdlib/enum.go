// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

type enumType string

func (e enumType) Name() string {
	return string(e)
}

func EnumType(e runtime.EnumProvider) runtime.Type {
	return enumType(e.TypeName())
}

type enumVal struct {
	v        int64
	provider runtime.EnumProvider
}

// ConvertTo implements Val.
func (e enumVal) ConvertTo(t runtime.Type) (any, error) {
	switch t {
	case enumType(t.Name()):
		return e.v, nil
	case IntType:
		return e.v, nil
	}
	return nil, fmt.Errorf("cannot convert enum %s to %s", e.provider.TypeName(), t.Name())
}

// Type implements Val.
func (e enumVal) Type() runtime.Type {
	return EnumType(e.provider)
}

// Value implements Val.
func (e enumVal) Value() any {
	return e.v
}

// TODO - Is this actually something we want consistently across enums?
// GetField implements StructureAccessible
func (e enumVal) GetField(field string) runtime.Val {
	if field == "string" {
		// Convert the value to  its name.
		if name, ok := e.provider.FindName(e.v); ok {
			return NewStringVal(name)
		}
	}
	return NewErrorVal(fmt.Errorf("invalid field on enum: %s", field))
}

func NewEnumVal(id int64, e runtime.EnumProvider) runtime.Val {
	return enumVal{id, e}
}
