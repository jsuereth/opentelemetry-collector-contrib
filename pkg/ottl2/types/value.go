// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"reflect" // Val interface defines the functions supported by all expression values.
)

// read-only access to a value.
type Val interface {
	// Understands how to convert this type to others.
	ConvertTo(typeDesc reflect.Type) (any, error)

	// Equal returns true if the `other` value has the same type and content as the implementing struct.
	// Equal(other Val) Val

	// Value returns the raw value of the instance which may not be directly compatible with the expression
	// language types.
	Value() any
}

// A holder for a value that can also be set.
type Var interface {
	Val
	// Sets the Value at a location.
	SetValue(Val) error
}

type Getter = func() Val
type Setter = func(Val) error
type getterSetterVar struct {
	getter Getter
	setter Setter
}

// ConvertTo implements Var.
func (g getterSetterVar) ConvertTo(typeDesc reflect.Type) (any, error) {
	return g.getter().ConvertTo(typeDesc)
}

// SetValue implements Var.
func (g getterSetterVar) SetValue(v Val) error {
	return g.setter(v)
}

// Value implements Var.
func (g getterSetterVar) Value() any {
	return g.getter().Value()
}

func NewGetterSetterVar(getter Getter, setter Setter) Var {
	return getterSetterVar{getter, setter}
}
