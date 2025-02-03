// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
)

func NewIsIntFunc() types.Function {
	return stdlib.NewSimpleFunc(
		"IsInt",
		1,
		func(v []types.Val) types.Val {
			if len(v) != 1 {
				return stdlib.NewErrorVal(fmt.Errorf("invalid arguments to IsInt"))
			}
			_, err := v[0].ConvertTo(reflect.TypeFor[int64]())
			return stdlib.NewBoolVal(err == nil)
		})
}
