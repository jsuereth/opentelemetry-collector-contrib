// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"reflect"
	"time"
)

var TimeType = NewPrimitiveType("time.Time")

type timeVal time.Time

// ConvertTo implements Val.
func (t timeVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// Type implements Val.
func (t timeVal) Type() Type {
	return TimeType
}

// Value implements Val.
func (t timeVal) Value() any {
	return time.Time(t)
}

func NewTimeVal(v time.Time) Val {
	return timeVal(v)
}
