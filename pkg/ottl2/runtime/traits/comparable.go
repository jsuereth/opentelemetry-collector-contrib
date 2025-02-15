// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package traits // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/traits"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

// For types that can be compared to each other.
type Comparable interface {
	// Whether this value is equal to another.
	Equals(other runtime.Val) bool
	// Whether this value is less than another.
	LessThan(other runtime.Val) bool
}
