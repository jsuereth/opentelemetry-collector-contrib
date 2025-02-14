// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/stretchr/testify/assert"
)

func addAll(args []types.Val) types.Val {
	result := int64(0)
	for _, a := range args {
		v, err := a.ConvertTo(IntType)
		if err != nil {
			return NewErrorVal(err)
		}
		result += v.(int64)
	}
	return NewIntVal(result)
}

func TestFunction_CallFunction(t *testing.T) {
	tests := []struct {
		name     string
		f        types.Function
		args     []types.Val
		named    map[string]types.Val
		expected types.Val
	}{
		{
			name: "positional only",
			f: NewFunc(
				"test",
				[]string{"", ""},
				map[string]types.Val{},
				addAll,
			),
			args:     []types.Val{NewIntVal(1), NewIntVal(1)},
			named:    map[string]types.Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]types.Val{},
				addAll,
			),
			args: []types.Val{},
			named: map[string]types.Val{
				"lhs": NewIntVal(1),
				"rhs": NewIntVal(1),
			},
			expected: NewIntVal(2),
		},
		{
			name: "default only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]types.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args:     []types.Val{},
			named:    map[string]types.Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named and default only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]types.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []types.Val{},
			named: map[string]types.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(4),
		},
		{
			name: "named and positional only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]types.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []types.Val{
				NewIntVal(5),
			},
			named: map[string]types.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(8),
		},
		{
			name: "named, defualt and positional",
			f: NewFunc(
				"test",
				[]string{"", "lhs", "rhs"},
				map[string]types.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(4),
				},
				addAll,
			),
			args: []types.Val{
				NewIntVal(5),
			},
			named: map[string]types.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(9),
		},
		{
			name: "reflect, positional",
			f:    NewExampleFunc(),
			args: []types.Val{
				NewIntVal(3),
				NewIntVal(2),
			},
			named:    map[string]types.Val{},
			expected: NewIntVal(1),
		},
		{
			name: "reflect, named",
			f:    NewExampleFunc(),
			args: []types.Val{},
			named: map[string]types.Val{
				"Left":  NewIntVal(3),
				"Right": NewIntVal(2),
			},
			expected: NewIntVal(1),
		},
		{
			name: "reflect, named and posiitonal",
			f:    NewExampleFunc(),
			args: []types.Val{
				NewIntVal(3),
			},
			named: map[string]types.Val{
				"Right": NewIntVal(2),
			},
			expected: NewIntVal(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := callFunction(tt.f, tt.args, tt.named)
			if result.Type() == ErrorType {
				assert.Fail(t, "Found error: %v", result)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Calls a function using its default and named arguments. (for testing only)
func callFunction(f types.Function, pos []types.Val, named map[string]types.Val) types.Val {
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

type ExampleFuncArgs struct {
	Left  int64
	Right int64
}

func NewExampleFunc() types.Function {
	return NewReflectFunc("-", &ExampleFuncArgs{}, func(args *ExampleFuncArgs) types.Val {
		return NewIntVal(args.Left - args.Right)
	})
}

func TestReflect_ArgumentNames(t *testing.T) {
	f := NewExampleFunc()
	assert.ElementsMatch(t, []string{"Left", "Right"}, f.ArgNames())
}

type ExampleOptionalFuncArgs struct {
	Left  Optional[int64]
	Right Optional[int64]
}

func NewExampleOptionalFunc() types.Function {
	return NewReflectFunc("-", &ExampleOptionalFuncArgs{}, func(args *ExampleOptionalFuncArgs) types.Val {
		result := int64(0)
		if !args.Left.IsEmpty() {
			result += args.Left.Get()
		}
		if !args.Right.IsEmpty() {
			result += args.Right.Get()
		}
		return NewIntVal(result)
	})
}

func TestReflect_DefaultArgumentValuesOptional(t *testing.T) {
	f := NewExampleOptionalFunc()
	assert.Equal(t, []string{"Left", "Right"}, f.ArgNames())
	assert.Equal(t, map[string]types.Val{"Left": NilVal, "Right": NilVal}, f.DefaultArgs())
}
