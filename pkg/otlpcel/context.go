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
		cel.CustomTypeProvider(&pdataTypeProvider{
			structs: []StructTypeProvider{
				testPSpanType,
			},
		}),
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

type pdataTypeProvider struct {
	structs []StructTypeProvider
}

// EnumValue implements types.Provider.
func (p *pdataTypeProvider) EnumValue(enumName string) ref.Val {
	return types.NewErr("unknown enum name '%s'", enumName)
}

// FindIdent implements types.Provider.
func (p *pdataTypeProvider) FindIdent(identName string) (ref.Val, bool) {
	return nil, false
}

// FindStructFieldNames implements types.Provider.
func (p *pdataTypeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	for _, st := range p.structs {
		if st.Name() == structType {
			return st.FindFieldNames(), true
		}
	}
	return nil, false
}

// FindStructFieldType implements types.Provider.
func (p *pdataTypeProvider) FindStructFieldType(structType string, fieldName string) (*types.FieldType, bool) {
	for _, st := range p.structs {
		if st.Name() == structType {
			return st.FindFieldType(fieldName)
		}
	}
	return nil, false
}

// FindStructType implements types.Provider.
func (p *pdataTypeProvider) FindStructType(structType string) (*types.Type, bool) {
	for _, st := range p.structs {
		if st.Name() == structType {
			return st.Type(), true
		}
	}
	return nil, false
}

// NewValue implements types.Provider.
func (p *pdataTypeProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	for _, st := range p.structs {
		if st.Name() == structType {
			return st.Construct(fields)
		}
	}
	return types.NewErr("Cannot find constructor for type: %s", structType)
}

var _ types.Provider = &pdataTypeProvider{}
