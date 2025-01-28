// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import "reflect"

// TODO - we should use type constructor type here.
var MapType = NewPrimitiveType("map")

type mapVal[V Val] map[string]V

// Type implements Val.
func (m mapVal[V]) Type() Type {
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
func (m mapVal[V]) GetKey(key string) Val {
	return m[key]
}

func NewMapVal[V Val](v map[string]V) Val {
	return (mapVal[V])(v)
}
