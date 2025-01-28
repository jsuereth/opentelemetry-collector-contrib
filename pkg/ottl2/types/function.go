// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

// This defines function execution within the runtime of OTTL.
type Function interface {
	// Calls the given function.
	// TODO - Remove named args?
	Call(args []Val, namedArgs map[string]Val) (Val, error)
	// TODO - meta information about a function, including
	// type-tests and expected parameters.
}

// This type does not support named args.
type simpleFunc struct {
	f func([]Val) (Val, error)
}

func (f *simpleFunc) Call(args []Val, nargs map[string]Val) (Val, error) {
	// TODO - deal with named args somewhere.
	return f.f(args)
}

func NewRawFunc(f func([]Val) (Val, error)) Function {
	return &simpleFunc{f}
}
