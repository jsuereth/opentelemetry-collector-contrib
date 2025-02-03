// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
) // This type does not support named args.
type simpleFunc struct {
	name    string
	numArgs int
	f       func([]types.Val) types.Val
}

// Name implements Function.
func (f *simpleFunc) Name() string {
	return f.name
}

// ArgNames implements Function.
func (f *simpleFunc) ArgNames() []string {
	return make([]string, f.numArgs)
}

// DefaultArgs implements Function.
func (f *simpleFunc) DefaultArgs() map[string]types.Val {
	return map[string]types.Val{}
}

func (f *simpleFunc) Call(args []types.Val) types.Val {
	return f.f(args)
}

func NewSimpleFunc(name string, numArgs int, f func([]types.Val) types.Val) types.Function {
	return &simpleFunc{name, numArgs, f}
}

type advancedFunction struct {
	name        string
	argNames    []string
	defaultArgs map[string]types.Val
	f           func([]types.Val) types.Val
}

// Name implements Function.
func (a advancedFunction) Name() string {
	return a.name
}

// ArgNames implements Function.
func (a advancedFunction) ArgNames() []string {
	return a.argNames
}

// Call implements Function.
func (a advancedFunction) Call(args []types.Val) types.Val {
	return a.f(args)
}

// DefaultArgs implements Function.
func (a advancedFunction) DefaultArgs() map[string]types.Val {
	return a.defaultArgs
}

// Constructs a function that can have named or default parameters.
func NewFunc(
	// Name of the function.
	name string,
	// Names of argument values.  Empty strings denote positional only args.
	argNames []string,
	// Default arguments. Must be used with named parameters.
	defaultArgs map[string]types.Val,
	// A positional-only implementation of the function.
	// All named/default arguments will be turned positional using argNames before calling this.
	f func([]types.Val) types.Val,
) types.Function {
	// TODO - verify named arguments exist in argument list.
	return advancedFunction{name, argNames, defaultArgs, f}
}

// Calls a function using its default and named arguments.
func CallFunction(f types.Function, pos []types.Val, named map[string]types.Val) types.Val {
	args, err := createArgs(f, pos, named)
	if err != nil {
		return NewErrorVal(err)
	}
	return f.Call(args)
}

// Takes a function and create the official argument set.
// This will take positional and named arguments, union with default arguments
// and return a purely positional argument list.
func createArgs(f types.Function, pos []types.Val, named map[string]types.Val) ([]types.Val, error) {
	defaults := f.DefaultArgs()
	names := f.ArgNames()
	result := make([]types.Val, len(names))
	for i, name := range names {
		if i < len(pos) {
			result[i] = pos[i]
		} else if v, ok := named[name]; name != "" && ok {
			result[i] = v
		} else if v, ok := defaults[name]; name != "" && ok {
			result[i] = v
		} else {
			if name != "" {
				return result, fmt.Errorf("invalid argument list for %s, missing paramater #%d: %s", f.Name(), i, name)
			} else {
				return result, fmt.Errorf("invalid argument list for %s, missing paramater #%d", f.Name(), i)
			}
		}
	}
	return result, nil
}
