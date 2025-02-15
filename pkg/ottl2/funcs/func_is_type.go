// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
)

type isTypeFunc struct {
	name string
	tpe  types.Type
}

// ArgNames implements types.Function.
func (i *isTypeFunc) ArgNames() []string {
	return []string{""}
}

// Call implements types.Function.
func (i *isTypeFunc) Call(args []types.Val) types.Val {
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
func (i *isTypeFunc) DefaultArgs() map[string]types.Val {
	return map[string]types.Val{}
}

// Name implements types.Function.
func (i *isTypeFunc) Name() string {
	return i.name
}

func NewIsTypeFunc(name string, tpe types.Type) types.Function {
	return &isTypeFunc{name, tpe}
}

func NewIsBoolFunc() types.Function {
	return NewIsTypeFunc("IsBool", stdlib.BoolType)
}

func NewIsIntFunc() types.Function {
	return NewIsTypeFunc("IsInt", stdlib.IntType)
}

func NewIsDoubleFunc() types.Function {
	return NewIsTypeFunc("IsDouble", stdlib.FloatType)
}

func NewIsStringFunc() types.Function {
	return NewIsTypeFunc("IsString", stdlib.StringType)
}

func NewIsListFunc() types.Function {
	return NewIsTypeFunc("IsList", stdlib.SliceType)
}

func NewIsMapFunc() types.Function {
	return NewIsTypeFunc("IsMap", stdlib.PmapType)
}
