// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
)

func NewSetFunc() types.Function {
	return stdlib.NewFunc(
		"set",
		[]string{"target", "value"},
		map[string]types.Val{},
		setImpl,
	)
}

func setImpl(args []types.Val) types.Val {
	if len(args) != 2 {
		return stdlib.NewErrorVal(fmt.Errorf("invalid # of arguments to set, found: %v", args))
	}
	target, ok := args[0].(types.Var)
	if !ok {
		return stdlib.NewErrorVal(fmt.Errorf("cannot set: %v", args[0]))
	}
	value := args[1]
	err := target.SetValue(value)
	if err != nil {
		stdlib.NewErrorVal(err)
	}
	return stdlib.NilVal
}
