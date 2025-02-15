// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package traits // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/traits"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

type Adder interface {
	// Add returns a combination of the current value and other value.
	//
	// If the other value is an unsupported type, an error is returned.
	Add(other runtime.Val) runtime.Val
}

// Divider interface to support '/' operator.
type Divider interface {
	// Divide returns the result of dividing the current value by the input
	// denominator.
	//
	// A denominator value of zero results in an error.
	Divide(denominator runtime.Val) runtime.Val
}

// Multiplier interface to support '*' operator.
type Multiplier interface {
	// Multiply returns the result of multiplying the current and input value.
	Multiply(other runtime.Val) runtime.Val
}

// Negater interface to support unary '-' and '!' operator.
type Negater interface {
	// Negate returns the complement of the current value.
	Negate() runtime.Val
}

// Subtractor interface to support binary '-' operator.
type Subtractor interface {
	// Subtract returns the result of subtracting the input from the current
	// value.
	Subtract(subtrahend runtime.Val) runtime.Val
}
