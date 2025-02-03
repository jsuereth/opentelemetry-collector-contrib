// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

// read-only access to a value.
type Val interface {
	// Understands how to convert this type to others.
	ConvertTo(Type) (any, error)

	// The type of this Val.
	Type() Type

	// Value returns the raw value of the instance which may not be directly compatible with the expression
	// language types.
	Value() any
}

// read/write access to a value.
type Var interface {
	Val
	// Sets the Value at a location.
	SetValue(Val) error
}
