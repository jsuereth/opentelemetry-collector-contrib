// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func testParser() Parser {
	return NewParser(NewParserEnvironemnt(map[string]runtime.Type{}, map[string]runtime.Function{}, []runtime.EnumProvider{}))
}

func TestParser_literals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "nil",
			input:    "nil",
			expected: nil,
		},
		{
			name:     "boolean true",
			input:    "true",
			expected: true,
		},
		{
			name:     "boolean false",
			input:    "false",
			expected: false,
		},
		{
			name:  "lists",
			input: "[true, true, false]",
			// TODO - we want this to be a list of booleans.
			expected: []any{true, true, false},
		},
		{
			name:     "strings",
			input:    "\"test\"",
			expected: "test",
		},
		{
			name:     "float64",
			input:    "1.0",
			expected: 1.0,
		},
		{
			name:     "int64",
			input:    "1045",
			expected: int64(1045),
		},
		{
			name:     "any maps",
			input:    "{\"key\": \"value\"}",
			expected: map[string]any{"key": "value"},
		},
	}

	p := testParser()
	ctx := context.Background()
	env := newEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Parsing expr: %s", tt.input)
			expr, err := p.ParseValueString(tt.input)
			assert.Nil(t, err)
			result := expr.Eval(ctx, env).Value()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_mathExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "addition",
			input:    "1 + 2",
			expected: (int64)(3),
		},
		{
			name:     "subtraction",
			input:    "2 - 1",
			expected: (int64)(1),
		},
	}
	p := testParser()
	ctx := context.Background()
	env := newEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := p.ParseValueString(tt.input)
			assert.Nil(t, err)
			result := expr.Eval(ctx, env).Value()
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

var testStructType = runtime.NewStructureType(
	"testStruct",
	map[string]runtime.Type{
		"value": stdlib.StringType,
	},
)

type testStruct struct {
	value string
}

func (t *testStruct) Type() runtime.Type {
	return testStructType
}

func (t *testStruct) Value() any {
	return *t
}

func (t *testStruct) ConvertTo(tpe runtime.Type) (any, error) {
	panic("unimplemented")
}

func (t *testStruct) GetField(field string) runtime.Val {
	return stdlib.NewStringVal(t.value)
}

func newTestStruct(value string) runtime.Val {
	return &testStruct{value}
}

func newTestEnv(cfg func(*TransformEnvironment)) *EvalContext {
	result := NewEvalContext()
	cfg(&result)
	var ctx EvalContext = result
	return &ctx
}
func newTestParserEnv(cfg func() ParserEnvironment) *ParserContext {
	result := cfg()
	var ctx ParserContext = result
	return &ctx
}

func TestParser_environment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		penv     *ParserContext
		env      *EvalContext
	}{
		{
			name:     "field",
			input:    "some.value",
			expected: "result",
			penv: newTestParserEnv(func() ParserEnvironment {
				return NewParserEnvironemnt(
					map[string]runtime.Type{"some": runtime.NewStructureType("custom", map[string]runtime.Type{
						"value": stdlib.StringType,
					})},
					map[string]runtime.Function{},
					[]runtime.EnumProvider{},
				)
			}),
			env: newTestEnv(func(te *TransformEnvironment) {
				te.WithVariable("some", newTestStruct("result"))
			}),
		},
		{
			name:     "index",
			input:    "list[1]",
			expected: "two",
			penv: newTestParserEnv(func() ParserEnvironment {
				return NewParserEnvironemnt(
					map[string]runtime.Type{"list": stdlib.ListType},
					map[string]runtime.Function{},
					[]runtime.EnumProvider{},
				)
			}),
			env: newTestEnv(func(te *TransformEnvironment) {
				te.WithVariable("list", stdlib.NewListVal([]runtime.Val{
					stdlib.NewStringVal("one"), stdlib.NewStringVal("two"),
				}))
			}),
		},
		{
			name:     "key",
			input:    "dict[\"hi\"]",
			expected: "one",
			penv: newTestParserEnv(func() ParserEnvironment {
				return NewParserEnvironemnt(
					map[string]runtime.Type{"dict": stdlib.MapType},
					map[string]runtime.Function{},
					[]runtime.EnumProvider{},
				)
			}),
			env: newTestEnv(func(te *TransformEnvironment) {
				te.WithVariable("dict", stdlib.NewMapVal(map[string]runtime.Val{
					"hi":  stdlib.NewStringVal("one"),
					"bye": stdlib.NewStringVal("two"),
				}))
			}),
		},
		{
			name:     "editor",
			input:    "doSomething(test)",
			expected: "test",
			penv: newTestParserEnv(func() ParserEnvironment {
				return NewParserEnvironemnt(
					map[string]runtime.Type{"test": stdlib.StringType},
					map[string]runtime.Function{
						"doSomething": stdlib.NewSimpleFunc("doSomething", 1, func(v []runtime.Val) runtime.Val {
							return v[0]
						}),
					},
					[]runtime.EnumProvider{},
				)
			}),
			env: newTestEnv(func(te *TransformEnvironment) {
				te.WithVariable("test", stdlib.NewStringVal("test"))
			}),
		},
		{
			name:     "enum",
			input:    "STATUS_CODE_OK",
			expected: int64(ptrace.StatusCodeOk),
			penv: newTestParserEnv(func() ParserEnvironment {
				return NewParserEnvironemnt(
					map[string]runtime.Type{},
					map[string]runtime.Function{},
					[]runtime.EnumProvider{stdlib.StatusCodeEnum},
				)
			}),
			env: newTestEnv(func(te *TransformEnvironment) {}),
		},
	}
	ctx := context.Background()
	var p Parser
	var env EvalContext
	for _, tt := range tests {
		if tt.penv != nil {
			p = NewParser(*tt.penv)
		} else {
			p = NewParser(NewParserEnvironemnt(
				map[string]runtime.Type{},
				map[string]runtime.Function{},
				[]runtime.EnumProvider{},
			))
		}
		if tt.env != nil {
			env = *tt.env
		} else {
			env = newEnv()
		}
		t.Run(tt.name, func(t *testing.T) {
			expr, err := p.ParseValueString(tt.input)
			assert.Nil(t, err, "Failed to parse: %s, error: %v, parser: %v", tt.input, err, p)
			result := expr.Eval(ctx, env).Value()
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParser_conditions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
		penv     *ParserContext
		env      *EvalContext
	}{
		{
			name:     "literal true",
			input:    "true",
			expected: true,
		},
		{
			name:     "literal false",
			input:    "false",
			expected: false,
		},
		{
			name:     "or expressions",
			input:    "false or true",
			expected: true,
		},
		{
			name:     "and expressions",
			input:    "false and true",
			expected: false,
		},
		{
			name:     "lt comparison",
			input:    "0 < 5",
			expected: true,
		},
		{
			name:     "lte comparison",
			input:    "5 <= 5",
			expected: true,
		},
		{
			name:     "gt comparison",
			input:    "0 > 5",
			expected: false,
		},
		{
			name:     "gte comparison",
			input:    "0 >= 5",
			expected: false,
		},
		{
			name:     "ne comparison",
			input:    "0 != 5",
			expected: true,
		},
	}
	ctx := context.Background()
	var p Parser
	var env EvalContext
	for _, tt := range tests {
		if tt.penv != nil {
			p = NewParser(*tt.penv)
		} else {
			p = NewParser(NewParserEnvironemnt(
				map[string]runtime.Type{},
				map[string]runtime.Function{},
				[]runtime.EnumProvider{},
			))
		}
		if tt.env != nil {
			env = *tt.env
		} else {
			env = newEnv()
		}
		t.Run(tt.name, func(t *testing.T) {
			expr, err := p.ParseConditionString(tt.input)
			assert.Nil(t, err)
			result := expr.Eval(ctx, env).Value()
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
