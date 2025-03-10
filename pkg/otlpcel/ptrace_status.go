// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var PTraceStatusType = NewStructType[ptrace.Status](
	"ptrace.Status",
	[]Field[ptrace.Status]{
		{
			Name: "code",
			Type: types.IntType,
			IsSet: func(target ptrace.Status) bool {
				return target.Code() != ptrace.StatusCodeUnset
			},
			GetFrom: func(target ptrace.Status) (any, error) {
				return target.Code(), nil
			},
		},
		{
			Name: "message",
			Type: types.StringType,
			IsSet: func(target ptrace.Status) bool {
				return len(target.Message()) != 0
			},
			GetFrom: func(target ptrace.Status) (any, error) {
				return target.Message(), nil
			},
		},
	},
	func(m map[string]any) (ptrace.Status, error) {
		return ptrace.NewStatus(), fmt.Errorf("unimplemented")
	},
	func(s ptrace.Status) ref.Val {
		return pTraceStatusWrapper(s)
	},
)

type pTraceStatusWrapper ptrace.Status

// ConvertToNative implements ref.Val.
func (p pTraceStatusWrapper) ConvertToNative(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// ConvertToType implements ref.Val.
func (p pTraceStatusWrapper) ConvertToType(typeValue ref.Type) ref.Val {
	panic("unimplemented")
}

// Equal implements ref.Val.
func (p pTraceStatusWrapper) Equal(other ref.Val) ref.Val {
	panic("unimplemented")
}

// Type implements ref.Val.
func (p pTraceStatusWrapper) Type() ref.Type {
	return PTraceStatusType.Type()
}

// Value implements ref.Val.
func (p pTraceStatusWrapper) Value() any {
	return ptrace.Status(p)
}
