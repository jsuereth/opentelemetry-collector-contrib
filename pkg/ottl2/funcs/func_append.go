// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - Figure out if we work in the `Val` domain, or the raw domain, or both.
type AppendArguments struct {
	Target types.Var
	Value  any           `ottl:"default=nil"`
	Values pcommon.Slice `ottl:"default=pcommon.Slice()"`
}

func NewAppendFunc() types.Function {
	return stdlib.NewReflectFunc(
		"append",
		&AppendArguments{},
		func(args *AppendArguments) types.Val {
			t := args.Target
			var res []any
			res = appendAny(res, t.Value())
			if args.Value != nil {
				res = appendAny(res, args.Value)
			}
			if args.Values.Len() > 0 {
				for _, a := range args.Values.AsRaw() {
					res = appendAny(res, a)
				}
			}

			// TODO - We may want to support different target types
			// so we don't force everything to be a slice.
			resSlice := pcommon.NewSlice()
			if err := resSlice.FromRaw(res); err != nil {
				return stdlib.NewErrorVal(err)
			}
			err := t.SetValue(stdlib.NewSliceVar(resSlice))
			if err != nil {
				return stdlib.NewErrorVal(err)
			}
			return stdlib.NilVal
		})
}

func appendAny(res []any, value any) []any {
	switch valueType := value.(type) {
	case pcommon.Slice:
		res = append(res, valueType.AsRaw()...)
	// Note: Compared to original OTTL, this PValues are erased to primitives at runtime.
	// We need to handle specific value types.
	// case pcommon.Value:
	case []string:
		res = appendMultiple(res, valueType)
	case []any:
		res = append(res, valueType...)
	case []int64:
		res = appendMultiple(res, valueType)
	case []bool:
		res = appendMultiple(res, valueType)
	case []float64:
		res = appendMultiple(res, valueType)

	case string:
		res = append(res, valueType)
	case int64:
		res = append(res, valueType)
	case bool:
		res = append(res, valueType)
	case float64:
		res = append(res, valueType)
	case any:
		res = append(res, valueType)
	default:
		res = append(res, value)
	}
	return res
}

func appendMultiple[K any](target []any, values []K) []any {
	for _, v := range values {
		target = append(target, v)
	}
	return target
}
