// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"reflect" // Val interface defines the functions supported by all expression values.
)

type Val interface {
	// Understands how to convert this type to others.
	ConvertTo(typeDesc reflect.Type) (any, error)

	// Equal returns true if the `other` value has the same type and content as the implementing struct.
	// Equal(other Val) Val

	// Value returns the raw value of the instance which may not be directly compatible with the expression
	// language types.
	Value() any
}
