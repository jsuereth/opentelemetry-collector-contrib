// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

// Environment used to parse ASTs
type TransformEnvironment struct {
	variables []VariableDecl
}

// VariableDecl defines a variable declaration which may optionally have a constant value.
type VariableDecl struct {
	name  string
	value Val
}

type EvalContext interface {
	// ResolveName returns a value from the activation by qualified name, or false if the name
	// could not be found.
	ResolveName(name string) (any, bool)

	// Parent returns the parent of the current activation, may be nil.
	// If non-nil, the parent will be searched during resolve calls.
	Parent() EvalContext
}
