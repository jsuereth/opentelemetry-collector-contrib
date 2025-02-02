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

	// The type of this Val.
	Type() Type

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
