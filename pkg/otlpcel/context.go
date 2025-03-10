// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func NewSpanEnv() (*cel.Env, error) {
	pdataConverter := pdataTypeProvider{
		structs: []StructTypeProvider{
			PSpanType,
			PTraceStatusType,
			PValueType,
		},
	}
	return cel.NewEnv(
		// TODO - pdataTypeProvidercan register types recursively if we're feeling lazy.
		cel.CustomTypeProvider(&pdataConverter),
		cel.CustomTypeAdapter(&pdataConverter),
		cel.Variable("span", PSpanType.Type()),
		// TOOD - can we infer this from our type provider as "setter"s?
		cel.Function("SetName",
			cel.MemberOverload("span_set_name",
				[]*cel.Type{PSpanType.Type(), cel.StringType},
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

func internalNativeToValue(value any) ref.Val {
	// TODO - register this.
	switch resultType := value.(type) {
	case ref.Val:
		return resultType
	case bool:
		return types.Bool(resultType)
	case string:
		return types.String(resultType)
	case int64:
		return types.Int(resultType)
	case float64:
		return types.Double(resultType)
	case ptrace.Span:
		return pspanWrapper(resultType)
	case ptrace.StatusCode:
		return types.Int(resultType)
	case pcommon.Map:
		return pmapWrapper(resultType)
	case pcommon.Value:
		return pvalueWrapper(resultType)
	}
	return types.NewErr("cannot convert type to ref.Val: '%v'", reflect.TypeOf(value))
}

// NativeToValue implements ref.TypeAdapter.
func (p *pdataTypeProvider) NativeToValue(value any) ref.Val {
	return internalNativeToValue(value)
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
