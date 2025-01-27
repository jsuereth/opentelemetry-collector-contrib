// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package traits // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/traits"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

// For types that support {name}[index]
type Indexable interface {
	// Obtains a value at an index
	GetIndex(index int64) types.Val
}

// For types that support {name}[key]
type KeyAccessable interface {
	// Obtains a value at a key
	GetKey(key string) types.Val
}

// For types that support {name}.{field}
type StructureAccessible interface {
	// Obtains a field by its name
	GetField(field string) types.Val
}
