// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package traits // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/traits"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

type Adder interface {
	// Add returns a combination of the current value and other value.
	//
	// If the other value is an unsupported type, an error is returned.
	Add(other types.Val) types.Val
}

// Divider interface to support '/' operator.
type Divider interface {
	// Divide returns the result of dividing the current value by the input
	// denominator.
	//
	// A denominator value of zero results in an error.
	Divide(denominator types.Val) types.Val
}

// Multiplier interface to support '*' operator.
type Multiplier interface {
	// Multiply returns the result of multiplying the current and input value.
	Multiply(other types.Val) types.Val
}

// Negater interface to support unary '-' and '!' operator.
type Negater interface {
	// Negate returns the complement of the current value.
	Negate() types.Val
}

// Subtractor interface to support binary '-' operator.
type Subtractor interface {
	// Subtract returns the result of subtracting the input from the current
	// value.
	Subtract(subtrahend types.Val) types.Val
}
