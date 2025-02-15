// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
)

type isTypeFunc struct {
	name string
	tpe  runtime.Type
}

// ArgNames implements types.Function.
func (i *isTypeFunc) ArgNames() []string {
	return []string{""}
}

// Call implements types.Function.
func (i *isTypeFunc) Call(args []runtime.Val) runtime.Val {
	if len(args) != 1 {
		return stdlib.NewErrorVal(fmt.Errorf("invalid arguments to %s", i.name))
	}
	arg := args[0]
	switch arg.Type() {
	case i.tpe:
		return stdlib.TrueVal
	// pcommon.Value could be anything so check it directly.
	case stdlib.PvalType:
		// Check
		_, err := arg.ConvertTo(i.tpe)
		return stdlib.NewBoolVal(err == nil)
	}
	return stdlib.FalseVal
}

// DefaultArgs implements types.Function.
func (i *isTypeFunc) DefaultArgs() map[string]runtime.Val {
	return map[string]runtime.Val{}
}

// Name implements types.Function.
func (i *isTypeFunc) Name() string {
	return i.name
}

func NewIsTypeFunc(name string, tpe runtime.Type) runtime.Function {
	return &isTypeFunc{name, tpe}
}

func NewIsBoolFunc() runtime.Function {
	return NewIsTypeFunc("IsBool", stdlib.BoolType)
}

func NewIsIntFunc() runtime.Function {
	return NewIsTypeFunc("IsInt", stdlib.IntType)
}

func NewIsDoubleFunc() runtime.Function {
	return NewIsTypeFunc("IsDouble", stdlib.FloatType)
}

func NewIsStringFunc() runtime.Function {
	return NewIsTypeFunc("IsString", stdlib.StringType)
}

func NewIsListFunc() runtime.Function {
	return NewIsTypeFunc("IsList", stdlib.SliceType)
}

func NewIsMapFunc() runtime.Function {
	return NewIsTypeFunc("IsMap", stdlib.PmapType)
}
