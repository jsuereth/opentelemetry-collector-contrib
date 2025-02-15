// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type statusCodeEnum struct{}

// FindName implements types.EnumProvider.
func (s statusCodeEnum) FindName(id int64) (string, bool) {
	switch ptrace.StatusCode(id) {
	case ptrace.StatusCodeUnset:
		return "STATUS_CODE_UNSET", true
	case ptrace.StatusCodeOk:
		return "STATUS_CODE_OK", true
	case ptrace.StatusCodeError:
		return "STATUS_CODE_ERROR", true
	}
	return "", false
}

// ResolveName implements types.EnumProvider.
func (s statusCodeEnum) ResolveName(name string) (int64, bool) {
	switch name {
	case "STATUS_CODE_UNSET":
		return int64(ptrace.StatusCodeUnset), true
	case "STATUS_CODE_OK":
		return int64(ptrace.StatusCodeOk), true
	case "STATUS_CODE_ERROR":
		return int64(ptrace.StatusCodeError), true
	}
	return 0, false
}

// TypeName implements types.EnumProvider.
func (s statusCodeEnum) TypeName() string {
	return "ptrace.StatusCode"
}

var StatusCodeEnum runtime.EnumProvider = statusCodeEnum{}
var StatusCodeType runtime.Type = EnumType(StatusCodeEnum)

func NewStatusCodeVal(v ptrace.StatusCode) runtime.Val {
	return NewEnumVal(int64(v), StatusCodeEnum)
}

type statusCodeVar struct {
	getter func() ptrace.StatusCode
	setter func(ptrace.StatusCode)
}

// ConvertTo implements types.Var.
func (s statusCodeVar) ConvertTo(t runtime.Type) (any, error) {
	return NewStatusCodeVal(s.getter()).ConvertTo(t)
}

// SetValue implements types.Var.
func (s statusCodeVar) SetValue(v runtime.Val) error {
	o, err := v.ConvertTo(IntType)
	if err != nil {
		return err
	}
	s.setter(ptrace.StatusCode(o.(int64)))
	return nil
}

// Type implements types.Var.
func (s statusCodeVar) Type() runtime.Type {
	return StatusCodeType
}

// Value implements types.Var.
func (s statusCodeVar) Value() any {
	return int64(s.getter())
}

func NewStatusCodeVar(getter func() ptrace.StatusCode, setter func(ptrace.StatusCode)) runtime.Var {
	return statusCodeVar{getter, setter}
}
