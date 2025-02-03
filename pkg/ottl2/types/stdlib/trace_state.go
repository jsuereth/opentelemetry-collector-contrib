// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/otel/trace"
)

var TraceStateType = types.NewPrimitiveType("pcommon.TraceState")

// Note: trace state HIDES a string.
type traceStateVar struct {
	getter func() string
	setter func(string)
}

// ConvertTo implements types.Var.
func (t traceStateVar) ConvertTo(tpe types.Type) (any, error) {
	switch tpe {
	case TraceStateType:
		return t.getter(), nil
	case StringType:
		return t.getter(), nil
	}
	return nil, fmt.Errorf("cannot convert tracestate to: %s", tpe.Name())
}

// SetValue implements types.Var.
func (t traceStateVar) SetValue(o types.Val) error {
	v, err := o.ConvertTo(StringType)
	if err != nil {
		return err
	}
	t.setter(v.(string))
	return nil
}

// Type implements types.Var.
func (t traceStateVar) Type() types.Type {
	return TraceStateType
}

// Value implements types.Var.
func (t traceStateVar) Value() any {
	return t.getter()
}

// GetKey implements KeyAccessable.
func (t traceStateVar) GetKey(key string) types.Val {
	return NewStringVar(
		func() string {
			if ts, err := trace.ParseTraceState(t.getter()); err == nil {
				return ts.Get(key)
			}
			// TODO - error here?
			return ""
		},
		func(s string) {
			if ts, err := trace.ParseTraceState(t.getter()); err == nil {
				if updated, err := ts.Insert(key, s); err == nil {
					t.setter(updated.String())
				}
			}
		},
	)
}

func NewTraceStateVar(getter func() string,
	setter func(string)) types.Var {
	return traceStateVar{getter, setter}
}
