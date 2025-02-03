// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

// TODO - we should use type constructor type here.
var MapType = types.NewPrimitiveType("map")

type mapVal[V types.Val] map[string]V

// Type implements Val.
func (m mapVal[V]) Type() types.Type {
	return MapType
}

// ConvertTo implements Val.
func (m mapVal[V]) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

func (m mapVal[V]) Value() any {
	// TODO - unwrap early or memoize this.
	result := map[string]any{}
	for k, v := range m {
		result[k] = v.Value()
	}
	return result
}

// maps are Keyable
func (m mapVal[V]) GetKey(key string) types.Val {
	return m[key]
}

func NewMapVal[V types.Val](v map[string]V) types.Val {
	return (mapVal[V])(v)
}
