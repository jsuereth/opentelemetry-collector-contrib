// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
)

// Environment used to parse ASTs
type EvalContext interface {
	// ResolveName returns a value from the context by qualified name, or false if the name
	// could not be structund.
	ResolveName(name string) (types.Val, bool)

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
	value types.Val
}

func NewEvalContext() TransformEnvironment {
	return TransformEnvironment{}
}

func (te *TransformEnvironment) WithVariable(name string, value types.Val) {
	// TODO - don't duplicate values.
	te.variables = append(te.variables, VariableDecl{name, value})
}

func (te TransformEnvironment) Parent() EvalContext {
	return nil
}

func (te TransformEnvironment) ResolveName(name string) (types.Val, bool) {
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
	ResolveName(name string) (types.Type, bool)

	// ResolveFunction returns a function form context by qualified name, or false if the name
	// could not be found.
	ResolveFunction(name string) (types.Function, bool)

	// Resolves an enumeration name into its value.
	ResolveEnum(name string) (types.Val, bool)
}

type EnumDefinition map[string]types.Val

type ParserEnvironment struct {
	variables map[string]types.Type
	functions map[string]types.Function
	enums     []EnumDefinition
}

func (pe ParserEnvironment) String() string {
	return fmt.Sprintf("Env{variables: %v, functions: %v}", pe.variables, pe.functions)
}

func NewParserEnvironemnt(
	variables map[string]types.Type,
	functions map[string]types.Function) ParserEnvironment {
	return ParserEnvironment{
		variables,
		functions,
		[]EnumDefinition{},
	}
}

func (p ParserEnvironment) ResolveName(name string) (types.Type, bool) {
	t, ok := p.variables[name]
	return t, ok
}

func (p ParserEnvironment) ResolveFunction(name string) (types.Function, bool) {
	f, ok := p.functions[name]
	return f, ok
}

func (p ParserEnvironment) ResolveEnum(name string) (types.Val, bool) {
	for _, es := range p.enums {
		for k, v := range es {
			if k == name {
				return v, true
			}
		}
	}
	return stdlib.NewErrorVal(fmt.Errorf("no such enum: %s", name)), false
}
