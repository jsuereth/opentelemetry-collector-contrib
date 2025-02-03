// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CallFunction(tt.f, tt.args, tt.named)
			assert.Equal(t, tt.expected, result)
		})
	}
}
