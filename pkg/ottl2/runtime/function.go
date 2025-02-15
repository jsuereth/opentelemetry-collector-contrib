// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package runtime // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"

// This defines function execution within the runtime of OTTL.

// A Custom functions defined for OTTL.
type Function interface {
	// The name of the function.
	Name() string
	// Calls the given function.
	// All named or default arguments MUST be turned positional before calling this.
	Call(args []Val) Val

	// The names of arguments in positional order.
	// Empty names are considered positional arguments only.
	ArgNames() []string

	// Default values for arguments.
	// Default arguments MUST have names.
	DefaultArgs() map[string]Val
}
