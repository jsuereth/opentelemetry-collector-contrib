// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

func NewIsIntFunc() types.Function {
	return types.NewSimpleFunc(
		"IsInt",
		1,
		func(v []types.Val) types.Val {
			if len(v) != 1 {
				return types.NewErrorVal(fmt.Errorf("invalid arguments to IsInt"))
			}
			_, err := v[0].ConvertTo(reflect.TypeFor[int64]())
			return types.NewBoolVal(err == nil)
		})
}
