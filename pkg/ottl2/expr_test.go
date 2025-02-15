// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"github.com/stretchr/testify/assert"
)

func newEnv() EvalContext {
	return NewEvalContext()
}

type testDef interface {
	RunTest(*testing.T)
}

type exprTest struct {
	name     string
	expr     Interpretable
	expected any
}

func (e *exprTest) RunTest(t *testing.T) {
	ctx := context.Background()
	env := newEnv()
	t.Run(e.name, func(t *testing.T) {
		v := e.expr.Eval(ctx, env)
		assert.NotNil(t, v)
		assert.Equal(t, e.expected, v.Value())
	})
}

func TestExpr_literals(t *testing.T) {
	tests := []testDef{
		&exprTest{
			name:     "nil",
			expr:     NilExpr(),
			expected: nil,
		},
		&exprTest{
			name:     "strings",
			expr:     StringExpr("test"),
			expected: "test",
		},
		&exprTest{
			name:     "ints",
			expr:     IntExpr(5),
			expected: int64(5),
		},
		&exprTest{
			name:     "floats",
			expr:     FloatExpr(5),
			expected: float64(5),
		},
		&exprTest{
			name:     "true",
			expr:     BooleanExpr(true),
			expected: true,
		},
		&exprTest{
			name:     "false",
			expr:     BooleanExpr(false),
			expected: false,
		},
		&exprTest{
			name:     "byteslice",
			expr:     ByteSliceExpr([]byte{0, 1}),
			expected: []byte{0, 1},
		},
		// TODO - we should try to keep the list homogenous
		&exprTest{
			name: "lists",
			expr: ListExpr([]Interpretable{
				StringExpr("one"),
				StringExpr("two"),
			}),
			expected: []any{"one", "two"},
		},
		// TODO - we should try to keep the map homogenous
		&exprTest{
			name: "maps",
			expr: MapExpr(map[string]Interpretable{
				"one": IntExpr(1),
				"two": IntExpr(2),
			}),
			expected: map[string]any{"one": (int64)(1), "two": (int64)(2)},
		},
	}
	for _, tt := range tests {
		tt.RunTest(t)
	}
}

func TestExpr_math(t *testing.T) {
	tests := []testDef{
		&exprTest{
			name:     "add ints",
			expr:     AddExpr(IntExpr(1), IntExpr(2)),
			expected: (int64)(3),
		},
		&exprTest{
			name:     "add floats",
			expr:     AddExpr(FloatExpr(1), FloatExpr(2)),
			expected: (float64)(3),
		},
	}
	for _, tt := range tests {
		tt.RunTest(t)
	}
}

func addAll(args []runtime.Val) runtime.Val {
	result := int64(0)
	for _, a := range args {
		v, err := a.ConvertTo(stdlib.IntType)
		if err != nil {
			return stdlib.NewErrorVal(err)
		}
		result += v.(int64)
	}
	return stdlib.NewIntVal(result)
}

func TestExpr_FunctionCal(t *testing.T) {
	tests := []struct {
		name     string
		f        runtime.Function
		args     []runtime.Val
		named    map[string]runtime.Val
		expected runtime.Val
	}{
		{
			name: "positional only",
			f: stdlib.NewFunc(
				"test",
				[]string{"", ""},
				map[string]runtime.Val{},
				addAll,
			),
			args:     []runtime.Val{stdlib.NewIntVal(1), stdlib.NewIntVal(1)},
			named:    map[string]runtime.Val{},
			expected: stdlib.NewIntVal(2),
		},
		{
			name: "named only",
			f: stdlib.NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{},
				addAll,
			),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
				"lhs": stdlib.NewIntVal(1),
				"rhs": stdlib.NewIntVal(1),
			},
			expected: stdlib.NewIntVal(2),
		},
		{
			name: "default only",
			f: stdlib.NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": stdlib.NewIntVal(1),
					"rhs": stdlib.NewIntVal(1),
				},
				addAll,
			),
			args:     []runtime.Val{},
			named:    map[string]runtime.Val{},
			expected: stdlib.NewIntVal(2),
		},
		{
			name: "named and default only",
			f: stdlib.NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": stdlib.NewIntVal(1),
					"rhs": stdlib.NewIntVal(1),
				},
				addAll,
			),
			args: []runtime.Val{},
			named: map[string]runtime.Val{
				"rhs": stdlib.NewIntVal(3),
			},
			expected: stdlib.NewIntVal(4),
		},
		{
			name: "named and positional only",
			f: stdlib.NewFunc(
				"test",
				[]string{"lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": stdlib.NewIntVal(1),
					"rhs": stdlib.NewIntVal(1),
				},
				addAll,
			),
			args: []runtime.Val{
				stdlib.NewIntVal(5),
			},
			named: map[string]runtime.Val{
				"rhs": stdlib.NewIntVal(3),
			},
			expected: stdlib.NewIntVal(8),
		},
		{
			name: "named, defualt and positional",
			f: stdlib.NewFunc(
				"test",
				[]string{"", "lhs", "rhs"},
				map[string]runtime.Val{
					"lhs": stdlib.NewIntVal(1),
					"rhs": stdlib.NewIntVal(4),
				},
				addAll,
			),
			args: []runtime.Val{
				stdlib.NewIntVal(5),
			},
			named: map[string]runtime.Val{
				"rhs": stdlib.NewIntVal(3),
			},
			expected: stdlib.NewIntVal(9),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			named := map[string]Interpretable{}
			for n, v := range tt.named {
				named[n] = ValExpr(v)
			}
			args := make([]Interpretable, len(tt.args))
			for i, v := range tt.args {
				args[i] = ValExpr(v)
			}
			expr := FuncCallExpr(tt.f, args, named)
			result := expr.Eval(context.Background(), newEnv())
			assert.Equal(t, tt.expected, result)
		})
	}
}
