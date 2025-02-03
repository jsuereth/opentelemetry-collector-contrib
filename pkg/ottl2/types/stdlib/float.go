// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

var FloatType = types.NewPrimitiveType("float")

type float64Val float64

// Type implements Val.
func (f float64Val) Type() types.Type {
	return FloatType
}

// ConvertTo implements Val.
func (f float64Val) ConvertTo(t types.Type) (any, error) {
	switch t {
	case FloatType:
		return float64(f), nil
	}
	return nil, fmt.Errorf("type conversion error from Double to '%v'", t)
}

func (f float64Val) Value() any {
	return (float64)(f)
}

// Flaots are addable
func (f float64Val) Add(o types.Val) types.Val {
	rhs, err := o.ConvertTo(FloatType)
	if err != nil {
		return NewErrorVal(err)
	}
	return (float64Val)((float64)(f) + rhs.(float64))
}

func NewFloatVal(v float64) types.Val {
	h := (float64Val)(v)
	return &h
}
