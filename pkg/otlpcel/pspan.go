// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var PSpanType = NewStructType[ptrace.Span](
	"ptrace.Span",
	[]Field[ptrace.Span]{
		{
			Name:  "name",
			Type:  types.StringType,
			IsSet: func(target ptrace.Span) bool { return true },
			GetFrom: func(target ptrace.Span) (any, error) {
				return target.Name(), nil
			},
		},
		{
			Name: "status",
			Type: PTraceStatusType.Type(),
			IsSet: func(target ptrace.Span) bool {
				return true
			},
			GetFrom: func(target ptrace.Span) (any, error) {
				return target.Status(), nil
			},
		},
		{
			Name: "attributes",
			Type: types.NewMapType(types.StringType, PValueType.Type()),
			IsSet: func(target ptrace.Span) bool {
				return target.Attributes().Len() != 0
			},
			GetFrom: func(target ptrace.Span) (any, error) {
				return target.Attributes(), nil
			},
		},
	},
	func(m map[string]any) (ptrace.Span, error) {
		panic("unimplemented")
	},
	func(s ptrace.Span) ref.Val {
		panic("unimplemented")
	},
	traits.IndexerType,
)

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
	return PSpanType.Type()
}

// Value implements ref.Val.
func (p pspanWrapper) Value() any {
	return ptrace.Span(p)
}

func NewCelSpan(p ptrace.Span) ref.Val {
	return pspanWrapper(p)
}
