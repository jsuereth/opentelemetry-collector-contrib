// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

var TimeType = types.NewPrimitiveType("time.Time")

type timeVal time.Time

// ConvertTo implements Val.
func (t timeVal) ConvertTo(tpe types.Type) (any, error) {
	panic("unimplemented")
}

// Type implements Val.
func (t timeVal) Type() types.Type {
	return TimeType
}

// Value implements Val.
func (t timeVal) Value() any {
	return time.Time(t)
}

func NewTimeVal(v time.Time) types.Val {
	return timeVal(v)
}

type timeVar struct {
	getter func() time.Time
	setter func(time.Time)
}

func (t timeVar) ConvertTo(tpe types.Type) (any, error) {
	panic("unimplemented")
}

func (t timeVar) SetValue(v types.Val) error {
	if value, ok := v.Value().(time.Time); ok {
		t.setter(value)
	}
	return nil
}

func (t timeVar) Type() types.Type {
	return TimeType
}

func (t timeVar) Value() any {
	return t.getter()
}

func NewTimeVar(getter func() time.Time, setter func(time.Time)) types.Var {
	return timeVar{getter, setter}
}
