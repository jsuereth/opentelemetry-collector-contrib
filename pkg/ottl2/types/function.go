// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"fmt" // This defines function execution within the runtime of OTTL.
)

type Function interface {
	// Calls the given function.
	// TODO - Remove named args?
	Call(args []Val) Val
	// TODO - meta information about a function, including
	// type-tests and expected parameters.

	// The names of arguments in positional order.
	ArgNames() []string
	DefaultArgs() map[string]Val
}

// Calls a function using its default and named arguments.
func CallFunction(f Function, pos []Val, named map[string]Val) Val {
	args, err := createArgs(f, pos, named)
	if err != nil {
		return NewErrorVal(err)
	}
	return f.Call(args)
}

// Takes a function and create the official argument set.
// This will take positional and named arguments, union with default arguments
// and return a purely positional argument list.
func createArgs(f Function, pos []Val, named map[string]Val) ([]Val, error) {
	defaults := f.DefaultArgs()
	names := f.ArgNames()
	result := make([]Val, len(names))
	for i, name := range names {
		if i < len(pos) {
			result[i] = pos[i]
		} else if v, ok := named[name]; name != "" && ok {
			result[i] = v
		} else if v, ok := defaults[name]; name != "" && ok {
			result[i] = v
		} else {
			return result, fmt.Errorf("invalid argument list for %v, missing paramater #%d (%s)", f, i, name)
		}
	}
	return result, nil
}

// This type does not support named args.
type simpleFunc struct {
	numArgs int
	f       func([]Val) Val
}

// ArgNames implements Function.
func (f *simpleFunc) ArgNames() []string {
	return make([]string, f.numArgs)
}

// DefaultArgs implements Function.
func (f *simpleFunc) DefaultArgs() map[string]Val {
	return map[string]Val{}
}

func (f *simpleFunc) Call(args []Val) Val {
	return f.f(args)
}

func NewSimpleFunc(numArgs int, f func([]Val) Val) Function {
	return &simpleFunc{numArgs, f}
}

type advancedFunction struct {
	argNames    []string
	defaultArgs map[string]Val
	f           func([]Val) Val
}

// ArgNames implements Function.
func (a advancedFunction) ArgNames() []string {
	return a.argNames
}

// Call implements Function.
func (a advancedFunction) Call(args []Val) Val {
	return a.f(args)
}

// DefaultArgs implements Function.
func (a advancedFunction) DefaultArgs() map[string]Val {
	return a.defaultArgs
}

// Constructs a function that can have named or default parameters.
func NewFunc(
	// Names of argument values.  Empty strings denote positional only args.
	argNames []string,
	// Default arguments. Must be used with named parameters.
	defaultArgs map[string]Val,
	// A positional-only implementation of the function.
	// All named/default arguments will be turned positional using argNames before calling this.
	f func([]Val) Val,
) Function {
	// TODO - verify named arguments exist in argument list.
	return advancedFunction{argNames, defaultArgs, f}
}
