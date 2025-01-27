// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

// This defines function execution within the runtime of OTTL.
type Function interface {
	// Calls the given function.
	Call(args []Val) (Val, error)
	// TODO - meta information about a function, including
	// type-tests and expected parameters.
}

type simpleFunc struct {
	f func([]Val) (Val, error)
}

func (f *simpleFunc) Call(args []Val) (Val, error) {
	return f.f(args)
}

func NewRawFunc(f func([]Val) (Val, error)) Function {
	return &simpleFunc{f}
}
