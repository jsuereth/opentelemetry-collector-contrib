// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
)

func NewSetFunc() runtime.Function {
	return stdlib.NewFunc(
		"set",
		[]string{"target", "value"},
		map[string]runtime.Val{},
		setImpl,
	)
}

func setImpl(args []runtime.Val) runtime.Val {
	if len(args) != 2 {
		return stdlib.NewErrorVal(fmt.Errorf("invalid # of arguments to set, found: %v", args))
	}
	target, ok := args[0].(runtime.Var)
	if !ok {
		return stdlib.NewErrorVal(fmt.Errorf("cannot set: %v", args[0]))
	}
	value := args[1]
	err := target.SetValue(value)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	return stdlib.NilVal
}
