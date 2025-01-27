// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testParser() Parser {
	return Parser{}
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

// func TestParser_mathExpressions(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		expected any
// 	}{
// 		{
// 			name:     "addition",
// 			input:    "1.0+2",
// 			expected: 3,
// 		},
// 	}
// 	p := testParser()
// 	ctx := context.Background()
// 	env := newEnv()
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			expr, err := p.ParseValueString(tt.input)
// 			assert.Nil(t, err)
// 			result, err := expr.Eval(ctx, env)
// 			assert.Nil(t, err)
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
