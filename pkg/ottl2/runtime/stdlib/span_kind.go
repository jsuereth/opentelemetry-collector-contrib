// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/traceutil"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type spanKindEnum struct{}

// FindName implements types.EnumProvider.
func (s spanKindEnum) FindName(id int64) (string, bool) {
	switch id {
	case int64(ptrace.SpanKindUnspecified):
		return "Unspecified", true
	case int64(ptrace.SpanKindInternal):
		return "Internal", true
	case int64(ptrace.SpanKindServer):
		return "Server", true
	case int64(ptrace.SpanKindClient):
		return "Client", true
	case int64(ptrace.SpanKindProducer):
		return "Producer", true
	case int64(ptrace.SpanKindConsumer):
		return "Consumer", true
	}
	return "", false
}

// ResolveName implements types.EnumProvider.
func (s spanKindEnum) ResolveName(name string) (int64, bool) {
	switch name {
	case "SPAN_KIND_UNSPECIFIED":
		return int64(ptrace.SpanKindUnspecified), true
	case "SPAN_KIND_INTERNAL":
		return int64(ptrace.SpanKindInternal), true
	case "SPAN_KIND_SERVER":
		return int64(ptrace.SpanKindServer), true
	case "SPAN_KIND_CLIENT":
		return int64(ptrace.SpanKindClient), true
	case "SPAN_KIND_PRODUCER":
		return int64(ptrace.SpanKindProducer), true
	case "SPAN_KIND_CONSUMER":
		return int64(ptrace.SpanKindConsumer), true
	}
	return 0, false
}

// TypeName implements types.EnumProvider.
func (s spanKindEnum) TypeName() string {
	return "ptrace.SpanKind"
}

var SpanKindEnum runtime.EnumProvider = spanKindEnum{}
var SpanKindType runtime.Type = EnumType(SpanKindEnum)

func NewSpanKindVal(v ptrace.SpanKind) runtime.Val {
	return NewEnumVal(int64(v), SpanKindEnum)
}

// TODO - we need to allow both `string` and `deprecated_string`
// for span kind.
type spanKindVar struct {
	getter func() ptrace.SpanKind
	setter func(ptrace.SpanKind)
}

func (e spanKindVar) GetField(field string) runtime.Val {
	switch field {
	case "string":
		return NewStringVar(
			func() string {
				return e.getter().String()
			},
			func(v string) {
				switch v {
				case "Unspecified":
					e.setter(ptrace.SpanKindUnspecified)
				case "Client":
					e.setter(ptrace.SpanKindClient)
				case "Server":
					e.setter(ptrace.SpanKindServer)
				case "Producer":
					e.setter(ptrace.SpanKindProducer)
				case "Consumer":
					e.setter(ptrace.SpanKindConsumer)
				}
			},
		)
	case "deprecated_string":
		return NewStringVar(
			func() string {
				return traceutil.SpanKindStr(e.getter())
			},
			func(v string) {
				// deprecated string uses values within OTTL, oddly.
				if id, ok := SpanKindEnum.ResolveName(v); ok {
					e.setter(ptrace.SpanKind(id))
				}
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field for SpanKind: %s", field))
}

// ConvertTo implements types.Var.
func (s spanKindVar) ConvertTo(t runtime.Type) (any, error) {
	return NewSpanKindVal(s.getter()).ConvertTo(t)
}

// SetValue implements types.Var.
func (s spanKindVar) SetValue(v runtime.Val) error {
	o, err := v.ConvertTo(IntType)
	if err != nil {
		return err
	}
	s.setter(ptrace.SpanKind(o.(int64)))
	return nil
}

// Type implements types.Var.
func (s spanKindVar) Type() runtime.Type {
	return SpanKindType
}

// Value implements types.Var.
func (s spanKindVar) Value() any {
	return int64(s.getter())
}

func NewSpanKindVar(getter func() ptrace.SpanKind, setter func(ptrace.SpanKind)) runtime.Var {
	return spanKindVar{getter, setter}
}
