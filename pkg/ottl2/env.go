// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
)

// Environment used to parse ASTs
type EvalContext interface {
	// ResolveName returns a value from the context by qualified name, or false if the name
	// could not be structund.
	ResolveName(name string) (runtime.Val, bool)

	// Parent returns the parent of the current activation, may be nil.
	// If non-nil, the parent will be searched during resolve calls.
	Parent() EvalContext
}

// Default implementation of EvalContext
type TransformEnvironment struct {
	variables []VariableDecl
}

// VariableDecl defines a variable declaration which may optionally have a constant value.
type VariableDecl struct {
	name  string
	value runtime.Val
}

func NewEvalContext() TransformEnvironment {
	return TransformEnvironment{}
}

func (te *TransformEnvironment) WithVariable(name string, value runtime.Val) {
	// TODO - don't duplicate values.
	te.variables = append(te.variables, VariableDecl{name, value})
}

func (te TransformEnvironment) Parent() EvalContext {
	return nil
}

func (te TransformEnvironment) ResolveName(name string) (runtime.Val, bool) {
	for _, v := range te.variables {
		if v.name == name {
			return v.value, true
		}
	}
	return stdlib.NilVal, false
}

// Context we need when evaluating parsed ASTs before turning them into Interpretable.
type ParserContext interface {
	// Returns true if the given name exists in the current context.
	ResolveName(name string) (runtime.Type, bool)

	// ResolveFunction returns a function form context by qualified name, or false if the name
	// could not be found.
	ResolveFunction(name string) (runtime.Function, bool)

	// Resolves an enumeration name into its value.
	ResolveEnum(name string) (runtime.Val, bool)
}

type ParserEnvironment struct {
	variables map[string]runtime.Type
	functions map[string]runtime.Function
	enums     []runtime.EnumProvider
}

func (pe ParserEnvironment) String() string {
	return fmt.Sprintf("Env{variables: %v, functions: %v}", pe.variables, pe.functions)
}

func NewParserEnvironemnt(
	variables map[string]runtime.Type,
	functions map[string]runtime.Function,
	enums []runtime.EnumProvider) ParserEnvironment {
	return ParserEnvironment{
		variables,
		functions,
		enums,
	}
}

func (p ParserEnvironment) ResolveName(name string) (runtime.Type, bool) {
	t, ok := p.variables[name]
	return t, ok
}

func (p ParserEnvironment) ResolveFunction(name string) (runtime.Function, bool) {
	f, ok := p.functions[name]
	return f, ok
}

func (p ParserEnvironment) ResolveEnum(name string) (runtime.Val, bool) {
	for _, provider := range p.enums {
		if id, ok := provider.ResolveName(name); ok {
			return stdlib.NewEnumVal(id, provider), true
		}
	}
	return stdlib.NewErrorVal(fmt.Errorf("no such enum name: %s", name)), false
}
