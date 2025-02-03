// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"github.com/stretchr/testify/assert"
)

// Define the context we'll evaluate statements against.
type testContext struct {
	name string
	age  int64
}

// Define the structure of the context for the parser to understand.
var testContextType = types.NewStructureType(
	"MyContext",
	map[string]types.Type{
		"name": stdlib.StringType,
		"age":  stdlib.IntType,
	},
)

// Note: MyContext MUST implement types.Val and traits.StructureAccessible
func (m *testContext) ConvertTo(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("unable to convert context")
}

// We define how to access values here.
func (m *testContext) GetField(field string) types.Val {
	switch field {
	case "name":
		return stdlib.NewStringVar(
			func() string {
				return m.name
			},
			func(s string) {
				m.name = s
			},
		)
	case "age":
		return stdlib.NewIntVar(
			func() int64 {
				return m.age
			},
			func(i int64) {
				m.age = i
			},
		)
	}
	return stdlib.NewErrorVal(fmt.Errorf("unknown field: %s", field))
}
func (m *testContext) Type() types.Type {
	return testContextType
}
func (m *testContext) Value() any {
	return m
}

func IsEmptyFunc() types.Function {
	return stdlib.NewSimpleFunc("IsEmpty", 1, func(args []types.Val) types.Val {
		r, err := args[0].ConvertTo(reflect.TypeFor[string]())
		if err != nil {
			return stdlib.NewErrorVal(err)
		}
		return stdlib.NewBoolVal(len(r.(string)) == 0)
	})
}

func RouteFunc() types.Function {
	return stdlib.NewSimpleFunc("route", 0, func(args []types.Val) types.Val {
		return stdlib.NewBoolVal(true)
	})
}

func Test_simple_e2e_readonly(t *testing.T) {
	env := NewTransformContext[testContext](
		testContextType,
		func(v *testContext) types.Val { return v },
		WithFunctions[testContext]([]types.Function{
			IsEmptyFunc(),
			RouteFunc(),
		}),
	)
	stmt, err := ParseStatement(env, "route() where IsEmpty(name)")
	assert.Nil(t, err)
	ctx := testContext{name: "test"}
	result, cond, err := stmt.Execute(context.Background(), &ctx)
	assert.Nil(t, err)
	assert.False(t, cond)
	assert.Nil(t, result)
}

func Test_simple_e2e_mutable(t *testing.T) {
	env := NewTransformContext[testContext](
		testContextType,
		func(v *testContext) types.Val { return v },
		WithFunctions[testContext]([]types.Function{
			IsEmptyFunc(),
			RouteFunc(),
		}),
		WithFunctions[testContext](funcs.StandardFuncs()),
	)
	ctx := testContext{name: "test", age: 21}
	stmt, err := ParseStatement(env, "set(age, 40)")
	assert.Nil(t, err)
	result, cond, err := stmt.Execute(context.Background(), &ctx)
	assert.Nil(t, err)
	assert.True(t, cond)
	assert.Nil(t, result)
	assert.Equal(t, int64(40), ctx.age, "context: %s", ctx)
}

func Benchmark_statement(t *testing.B) {
	tests := []struct {
		name       string
		ottl       string
		expectErr  bool
		expectCond bool
		expected   any
		expect     func(t *testing.B, ctx testContext)
	}{
		{
			name:       "simple route",
			ottl:       "route() where IsEmpty(name)",
			expectErr:  false,
			expectCond: false,
		},
		{
			name:       "set value with constant",
			ottl:       "set(name, \"hello world\")",
			expectErr:  false,
			expectCond: true,
			expect: func(t *testing.B, ctx testContext) {
				assert.Equal(t, "hello world", ctx.name)
			},
		},
		// {
		// 	name:       "set value with mutated calculation",
		// 	ottl:       "set(age, age + 5) where not IsEmpty(name)",
		// 	expectErr:  false,
		// 	expectCond: true,
		// 	expect: func(t *testing.B, ctx testContext) {
		// 		assert.Equal(t, int64(26), ctx.age)
		// 	},
		// },
	}
	env := NewTransformContext[testContext](
		testContextType,
		func(v *testContext) types.Val { return v },
		WithFunctions[testContext]([]types.Function{
			IsEmptyFunc(),
			RouteFunc(),
		}),
		WithFunctions[testContext](funcs.StandardFuncs()),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			ctx := testContext{name: "test", age: 21}
			var stmt *Statement[testContext]
			// Benchmark compiling different types of statements.
			t.Run("parse", func(t *testing.B) {
				s, err := ParseStatement(env, tt.ottl)
				assert.Nil(t, err)
				stmt = &s
			})
			// If compile was succesful, benchmark evaluation.
			if stmt != nil {
				t.Run("eval", func(t *testing.B) {
					result, cond, err := stmt.Execute(context.Background(), &ctx)
					if tt.expectErr {
						assert.NotNil(t, err)
					} else {
						assert.Nil(t, err)
					}
					assert.Equal(t, tt.expectCond, cond)
					assert.Equal(t, tt.expected, result)
					if tt.expect != nil {
						tt.expect(t, ctx)
					}
				})
			}
		})
	}
}
