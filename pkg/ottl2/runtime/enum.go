// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package runtime // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

// Provides an Enumerated value to OTTL.
type EnumProvider interface {
	// A type name for the enumeration.
	TypeName() string
	// Find the name of an enumerated value.
	FindName(id int64) (string, bool)
	// Find the value of enumeration name.
	ResolveName(name string) (int64, bool)
}
