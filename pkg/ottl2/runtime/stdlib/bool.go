// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"
	// Stores literal boolean values
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

var BoolType = runtime.NewPrimitiveType("bool")

type boolVal bool

// Type implements Val.
func (b boolVal) Type() runtime.Type {
	return BoolType
}

// How we coeerce between known types in OTLP.
func (b boolVal) ConvertTo(t runtime.Type) (any, error) {
	switch t {
	case BoolType:
		return bool(b), nil
	}
	return nil, fmt.Errorf("type conversion error from bool to '%v'", t)
}

func (b boolVal) Value() any {
	return (bool)(b)
}

func NewBoolVal(v bool) runtime.Val {
	return (boolVal)(v)
}

var (
	TrueVal  = NewBoolVal(true)
	FalseVal = NewBoolVal(false)
)
