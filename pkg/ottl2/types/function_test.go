// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func addAll(args []Val) Val {
	result := int64(0)
	for _, a := range args {
		v, err := a.ConvertTo(reflect.TypeFor[int64]())
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
		f        Function
		args     []Val
		named    map[string]Val
		expected Val
	}{
		{
			name: "positional only",
			f: NewFunc(
				[]string{"", ""},
				map[string]Val{},
				addAll,
			),
			args:     []Val{NewIntVal(1), NewIntVal(1)},
			named:    map[string]Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named only",
			f: NewFunc(
				[]string{"lhs", "rhs"},
				map[string]Val{},
				addAll,
			),
			args: []Val{},
			named: map[string]Val{
				"lhs": NewIntVal(1),
				"rhs": NewIntVal(1),
			},
			expected: NewIntVal(2),
		},
		{
			name: "default only",
			f: NewFunc(
				[]string{"lhs", "rhs"},
				map[string]Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args:     []Val{},
			named:    map[string]Val{},
			expected: NewIntVal(2),
		},
		{
			name: "named and default only",
			f: NewFunc(
				[]string{"lhs", "rhs"},
				map[string]Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []Val{},
			named: map[string]Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(4),
		},
		{
			name: "named and positional only",
			f: NewFunc(
				[]string{"lhs", "rhs"},
				map[string]Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(1),
				},
				addAll,
			),
			args: []Val{
				NewIntVal(5),
			},
			named: map[string]Val{
				"rhs": NewIntVal(3),
			},
			expected: NewIntVal(8),
		},
		{
			name: "named, defualt and positional",
			f: NewFunc(
				[]string{"", "lhs", "rhs"},
				map[string]Val{
					"lhs": NewIntVal(1),
					"rhs": NewIntVal(4),
				},
				addAll,
			),
			args: []Val{
				NewIntVal(5),
			},
			named: map[string]Val{
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
