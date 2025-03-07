// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var PSpanType *cel.Type = types.NewObjectType("ptrace.Span", traits.IndexerType)
var spanFields = map[string]*types.FieldType{
	"name": {
		Type:  types.StringType,
		IsSet: func(target any) bool { return true },
		GetFrom: func(target any) (any, error) {
			// TODO - cast and ereturn error
			return target.(ptrace.Span).Name(), nil
		},
	},
	"status": {
		// TODO - status type.
		Type:  types.IntType,
		IsSet: func(target any) bool { return true },
		GetFrom: func(target any) (any, error) {
			return target.(ptrace.Span).Status(), nil
		},
	},
	"attributes": {
		// TODO pcommon.Map
		Type: types.MapType,
		IsSet: func(target any) bool {
			return true
		},
		GetFrom: func(target any) (any, error) {
			return target.(ptrace.Span).Attributes(), nil
		},
	},
}

type pspanWrapper ptrace.Span

// ConvertToNative implements ref.Val.
func (p pspanWrapper) ConvertToNative(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// ConvertToType implements ref.Val.
func (p pspanWrapper) ConvertToType(typeValue ref.Type) ref.Val {
	panic("unimplemented")
}

// Equal implements ref.Val.
func (p pspanWrapper) Equal(other ref.Val) ref.Val {
	panic("unimplemented")
}

// Type implements ref.Val.
func (p pspanWrapper) Type() ref.Type {
	return PSpanType
}

// Value implements ref.Val.
func (p pspanWrapper) Value() any {
	return ptrace.Span(p)
}

func NewCelSpan(p ptrace.Span) ref.Val {
	return pspanWrapper(p)
}
