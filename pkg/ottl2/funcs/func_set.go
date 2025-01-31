// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

func NewSetFunc() types.Function {
	return types.NewFunc(
		"set",
		[]string{"target", "value"},
		map[string]types.Val{},
		setImpl,
	)
}

func setImpl(args []types.Val) types.Val {
	if len(args) != 2 {
		return types.NewErrorVal(fmt.Errorf("invalid # of arguments to set, found: %v", args))
	}
	target, ok := args[0].(types.Var)
	if !ok {
		return types.NewErrorVal(fmt.Errorf("cannot set: %v", args[0]))
	}
	value := args[1]
	err := target.SetValue(value)
	if err != nil {
		types.NewErrorVal(err)
	}
	return types.NilVal
}
