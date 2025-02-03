// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

// TODO - we should use type-constructor type for lists.
var ListType = types.NewPrimitiveType("list")

type listVal[T types.Val] []T

// Type implements Val.
func (i listVal[T]) Type() types.Type {
	return ListType
}

func (i listVal[T]) ConvertTo(t types.Type) (any, error) {
	panic("unimplemented")
}

func (i listVal[T]) Value() any {
	// TODO - memoize this?
	result := make([]any, len(i))
	for idx, v := range i {
		result[idx] = v.Value()
	}
	return result
}

// List is Indexable
func (l listVal[T]) GetIndex(index int64) types.Val {
	return l[index]
}

func NewListVal[T types.Val](v []T) types.Val {
	return (listVal[T])(v)
}
