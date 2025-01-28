// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import "reflect"

// TODO - we should use type-constructor type for lists.
var ListType = NewPrimitiveType("list")

type listVal[T Val] []T

// Type implements Val.
func (i listVal[T]) Type() Type {
	return ListType
}

func (i listVal[T]) ConvertTo(typeDesc reflect.Type) (any, error) {
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
func (l listVal[T]) GetIndex(index int64) Val {
	return l[index]
}

func NewListVal[T Val](v []T) Val {
	return (listVal[T])(v)
}
