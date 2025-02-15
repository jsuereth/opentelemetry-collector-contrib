// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

// TODO - we should use type-constructor type for lists.
var ListType = runtime.NewPrimitiveType("list")

type listVal[T runtime.Val] []T

// Type implements Val.
func (i listVal[T]) Type() runtime.Type {
	return ListType
}

func (i listVal[T]) ConvertTo(t runtime.Type) (any, error) {
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
func (l listVal[T]) GetIndex(index int64) runtime.Val {
	return l[index]
}

func NewListVal[T runtime.Val](v []T) runtime.Val {
	return (listVal[T])(v)
}
