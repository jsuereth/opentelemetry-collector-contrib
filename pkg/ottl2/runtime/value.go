// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package runtime // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

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

// For values that support {name}[index]
type Indexable interface {
	Val
	// Obtains a value at an index
	GetIndex(index int64) Val
}

// For values that support {name}[key]
type KeyIndexable interface {
	Val
	// Obtains a value at a key
	GetKey(key string) Val
}

// For values that support {name}.{field}
type Structure interface {
	Val
	// Obtains a field by its name
	GetField(field string) Val
}
