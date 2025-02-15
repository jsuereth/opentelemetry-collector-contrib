// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func addAll(args []runtime.Val) runtime.Val {
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
		f        runtime.Function
		args     []runtime.Val
		named    map[string]runtime.Val
		expected runtime.Val
	}{
		{
			name: "positional only",
			f: NewFunc(
				"test",
				[]string{"", ""},
				map[string]runtime.Val{},
				addAll,
			),
			args:     []runtime.Val{NewIntVal(1), NewIntVal(1)},
			named:    map[string]runtime.Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{},
				addAll,
			),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
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
				map[string]runtime.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args:     []runtime.Val{},
			named:    map[string]runtime.Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named and default only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(4),
		},
		{
			name: "named and positional only",
			f: NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []runtime.Val{
				NewIntVal(5),
			},
			named: map[string]runtime.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(8),
		},
		{
			name: "named, defualt and positional",
			f: NewFunc(
				"test",
				[]string{"", "lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(4),
				},
				addAll,
			),
			args: []runtime.Val{
				NewIntVal(5),
			},
			named: map[string]runtime.Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(9),
		},
		{
			name: "reflect, positional",
			f:    NewExampleFunc(),
			args: []runtime.Val{
				NewIntVal(3),
				NewIntVal(2),
			},
			named:    map[string]runtime.Val{},
			expected: NewIntVal(1),
		},
		{
			name: "reflect, named",
			f:    NewExampleFunc(),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
				"Left":  NewIntVal(3),
				"Right": NewIntVal(2),
			},
			expected: NewIntVal(1),
		},
		{
			name: "reflect, named and posiitonal",
			f:    NewExampleFunc(),
			args: []runtime.Val{
				NewIntVal(3),
			},
			named: map[string]runtime.Val{
				"Right": NewIntVal(2),
			},
			expected: NewIntVal(1),
		},
		{
			name:     "reflect, defaults",
			f:        NewExampleOptionalFunc(),
			args:     []runtime.Val{},
			named:    map[string]runtime.Val{},
			expected: NewIntVal(2),
		},
		{
			name: "reflect, defaults and named",
			f:    NewExampleOptionalFunc(),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
				"Right": NewIntVal(2),
			},
			expected: NewIntVal(4),
		},
		{
			name: "reflect, positional",
			f:    NewExampleOptionalFunc(),
			args: []runtime.Val{
				NewIntVal(3),
				NewIntVal(3),
			},
			named:    map[string]runtime.Val{},
			expected: NewIntVal(6),
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
func callFunction(f runtime.Function, pos []runtime.Val, named map[string]runtime.Val) runtime.Val {
	args, err := createArgs(f, pos, named)
	if err != nil {
		return NewErrorVal(err)
	}
	return f.Call(args)
}

// Takes a function and create the official argument set.
// This will take positional and named arguments, union with default arguments
// and return a purely positional argument list.
func createArgs(f runtime.Function, pos []runtime.Val, named map[string]runtime.Val) ([]runtime.Val, error) {
	defaults := f.DefaultArgs()
	names := f.ArgNames()
	result := make([]runtime.Val, len(names))
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

func NewExampleFunc() runtime.Function {
	return NewReflectFunc("-", &ExampleFuncArgs{}, func(args *ExampleFuncArgs) runtime.Val {
		return NewIntVal(args.Left - args.Right)
	})
}

func TestReflect_ArgumentNames(t *testing.T) {
	f := NewExampleFunc()
	assert.ElementsMatch(t, []string{"Left", "Right"}, f.ArgNames())
}

type ExampleOptionalFuncArgs struct {
	Left  int64 `ottl:"default=2"`
	Right int64 `ottl:"default=0"`
}

func NewExampleOptionalFunc() runtime.Function {
	return NewReflectFunc("+", &ExampleOptionalFuncArgs{}, func(args *ExampleOptionalFuncArgs) runtime.Val {
		return NewIntVal(args.Left + args.Right)
	})
}

func TestReflect_DefaultArgumentValuesOptional(t *testing.T) {
	f := NewExampleOptionalFunc()
	assert.Equal(t, []string{"Left", "Right"}, f.ArgNames())
	assert.Equal(t, map[string]runtime.Val{"Left": NewIntVal(2), "Right": NewIntVal(0)}, f.DefaultArgs())
}

type AdvancedDefaultFuncArgs struct {
	Slice pcommon.Slice `ottl:"default=pcommon.Slice()"`
}

func TestReflect_DefaultAdvancedValues(t *testing.T) {
	f := NewReflectFunc("test", &AdvancedDefaultFuncArgs{},
		func(t *AdvancedDefaultFuncArgs) runtime.Val {
			return NewSliceVar(t.Slice)
		})
	args := f.DefaultArgs()
	assert.Contains(t, args, "Slice")
	sliceDefault := args["Slice"]
	sliceDefaultValue, ok := sliceDefault.Value().(pcommon.Slice)
	assert.True(t, ok, "default is not a slice", sliceDefault)
	assert.Equal(t, []any{}, sliceDefaultValue.AsRaw())

	mySlice := pcommon.NewSlice()
	mySlice.AppendEmpty().SetInt(1)
	mySlice.AppendEmpty().SetInt(2)
	result := f.Call([]runtime.Val{NewSliceVar(mySlice)})
	resultSlice, ok := result.Value().(pcommon.Slice)
	assert.True(t, ok, "function call result is not a slice", result)
	assert.Equal(t, []any{int64(1), int64(2)}, resultSlice.AsRaw())
}
