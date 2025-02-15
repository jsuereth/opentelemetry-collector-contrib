// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

// TODO - we should use type constructor type here.
var MapType = runtime.NewPrimitiveType("map")

type mapVal[V runtime.Val] map[string]V

// Type implements Val.
func (m mapVal[V]) Type() runtime.Type {
	return MapType
}

// ConvertTo implements Val.
func (m mapVal[V]) ConvertTo(t runtime.Type) (any, error) {
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
func (m mapVal[V]) GetKey(key string) runtime.Val {
	return m[key]
}

func NewMapVal[V runtime.Val](v map[string]V) runtime.Val {
	return (mapVal[V])(v)
}
