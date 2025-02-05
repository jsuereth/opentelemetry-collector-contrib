// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanStatusType = types.NewStructureType("ptrace.Status", map[string]types.Type{
	"code":    StatusCodeType,
	"message": StringType,
})

type spanStatusVar struct {
	getter func() ptrace.Status
	setter func(ptrace.Status)
}

// ConvertTo implements types.Var.
func (s spanStatusVar) ConvertTo(types.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements types.Var.
func (s spanStatusVar) SetValue(v types.Val) error {
	if st, ok := v.Value().(ptrace.Status); ok {
		s.setter(st)
		return nil
	}
	return fmt.Errorf("invalid span status: %v", v)
}

// Type implements types.Var.
func (s spanStatusVar) Type() types.Type {
	return SpanStatusType
}

// Value implements types.Var.
func (s spanStatusVar) Value() any {
	return s.getter()
}

func (s spanStatusVar) GetField(field string) types.Val {
	switch field {
	case "code":
		return NewStatusCodeVar(
			func() ptrace.StatusCode {
				return s.getter().Code()
			},
			func(v ptrace.StatusCode) {
				// TODO - less copying
				current := s.getter()
				current.SetCode(v)
				s.setter(current)
			},
		)
	case "message":
		return NewStringVar(
			func() string {
				return s.getter().Message()
			},
			func(msg string) {
				current := s.getter()
				current.SetMessage(msg)
				s.setter(current)
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field on span status: %s", field))
}

func NewSpanStatusVar(getter func() ptrace.Status,
	setter func(ptrace.Status)) types.Var {
	return spanStatusVar{getter, setter}
}
