// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

var StringType = types.NewPrimitiveType("string")

type stringVal string

// Type implements Val.
func (s stringVal) Type() types.Type {
	return StringType
}

// Cast strings to valid other types.
func (s stringVal) ConvertTo(t types.Type) (any, error) {
	return convertStringTo(string(s), t)
}

func (s stringVal) Value() any {
	return (string)(s)
}

func NewStringVal(v string) types.Val {
	return (stringVal)(v)
}

type string64Var struct {
	getter func() string
	setter func(string)
}

func (s string64Var) ConvertTo(t types.Type) (any, error) {
	return convertStringTo(s.getter(), t)
}

func (i string64Var) SetValue(v types.Val) error {
	// TODO - check types first...
	if value, ok := v.Value().(string); ok {
		i.setter(value)
		return nil
	}
	return fmt.Errorf("value is not a string: %v", v)
}

func (i string64Var) Type() types.Type {
	return StringType
}

func (i string64Var) Value() any {
	return i.getter()
}

func NewStringVar(getter func() string, setter func(string)) types.Var {
	return string64Var{getter, setter}
}

func convertStringTo(s string, t types.Type) (any, error) {
	switch t {
	case StringType:
		return s, nil
	}
	return nil, fmt.Errorf(
		"unsupported native conversion from string to '%v'", t.Name())
}
