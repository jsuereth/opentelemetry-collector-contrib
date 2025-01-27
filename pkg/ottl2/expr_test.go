// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEvalContext struct{}

// Parent implements EvalContext.
func (t *testEvalContext) Parent() EvalContext {
	panic("unimplemented")
}

// ResolveName implements EvalContext.
func (t *testEvalContext) ResolveName(name string) (any, bool) {
	panic("unimplemented")
}

func newEnv() EvalContext {
	return &testEvalContext{}
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
