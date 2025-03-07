// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func NewSpanEnv() (*cel.Env, error) {
	return cel.NewEnv(
		cel.Types(PSpanType),
		cel.Variable("span", PSpanType),
		cel.CustomTypeProvider(&pdataTypeProvider{}),
		cel.Function("SetName",
			cel.MemberOverload("span_set_name",
				[]*cel.Type{PSpanType, cel.StringType},
				cel.BoolType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					name := rhs.Value().(string)
					span := lhs.Value().(ptrace.Span)
					span.SetName(name)
					return types.True
				}),
			),
		),
	)
}

func NewSpanActivation(span ptrace.Span) (cel.Activation, error) {
	bindings := map[string]any{
		"span": NewCelSpan(span),
	}
	return cel.NewActivation(bindings)
}

type pdataTypeProvider struct{}

// EnumValue implements types.Provider.
func (p *pdataTypeProvider) EnumValue(enumName string) ref.Val {
	return types.NewErr("unknown enum name '%s'", enumName)
}

// FindIdent implements types.Provider.
func (p *pdataTypeProvider) FindIdent(identName string) (ref.Val, bool) {
	return nil, false
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i += 1
	}
	return result
}

// FindStructFieldNames implements types.Provider.
func (p *pdataTypeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	switch structType {
	case PSpanType.TypeName():
		return mapKeys(spanFields), true
	}
	return nil, false
}

// FindStructFieldType implements types.Provider.
func (p *pdataTypeProvider) FindStructFieldType(structType string, fieldName string) (*types.FieldType, bool) {
	switch structType {
	case PSpanType.TypeName():
		field, ok := spanFields[fieldName]
		return field, ok
	}
	return nil, false
}

// FindStructType implements types.Provider.
func (p *pdataTypeProvider) FindStructType(structType string) (*types.Type, bool) {
	switch structType {
	case PSpanType.TypeName():
		return PSpanType, true
	}
	return nil, false
}

// NewValue implements types.Provider.
func (p *pdataTypeProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	panic("unimplemented")
}

var _ types.Provider = &pdataTypeProvider{}
