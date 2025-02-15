// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

var TimeType = runtime.NewPrimitiveType("time.Time")

type timeVal time.Time

// ConvertTo implements Val.
func (t timeVal) ConvertTo(tpe runtime.Type) (any, error) {
	panic("unimplemented")
}

// Type implements Val.
func (t timeVal) Type() runtime.Type {
	return TimeType
}

// Value implements Val.
func (t timeVal) Value() any {
	return time.Time(t)
}

func NewTimeVal(v time.Time) runtime.Val {
	return timeVal(v)
}

type timeVar struct {
	getter func() time.Time
	setter func(time.Time)
}

func (t timeVar) ConvertTo(tpe runtime.Type) (any, error) {
	panic("unimplemented")
}

func (t timeVar) SetValue(v runtime.Val) error {
	if value, ok := v.Value().(time.Time); ok {
		t.setter(value)
	}
	return nil
}

func (t timeVar) Type() runtime.Type {
	return TimeType
}

func (t timeVar) Value() any {
	return t.getter()
}

func NewTimeVar(getter func() time.Time, setter func(time.Time)) runtime.Var {
	return timeVar{getter, setter}
}
