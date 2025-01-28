// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
) // Environment used to parse ASTs

type EvalContext interface {
	// ResolveName returns a value from the context by qualified name, or false if the name
	// could not be found.
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
	return types.NilVal, false
}

// Context we need when evaluating parsed ASTs before turning them into Interpretable.
type ParserContext interface {
	// Returns true if the given name exists in the current context.
	HasName(name string) bool

	// ResolveFunction returns a function form context by qualified name, or false if the name
	// could not be found.
	ResolveFunction(name string) (types.Function, bool)

	// Resolves an enumeration name into its value.
	ResolveEnum(name string) (types.Val, bool)
}

type ParserEnvironment struct {
	variableNames []string
	functions     map[string]types.Function
}

func (pe ParserEnvironment) String() string {
	return fmt.Sprintf("Env{variabless: %v, functions: %v}", pe.variableNames, pe.functions)
}

func NewParserEnvironemnt(
	variableNames []string,
	functions map[string]types.Function) ParserEnvironment {
	return ParserEnvironment{
		variableNames,
		functions,
	}
}

func (p ParserEnvironment) HasName(name string) bool {
	for _, n := range p.variableNames {
		if n == name {
			return true
		}
	}
	return false
}

func (p ParserEnvironment) ResolveFunction(name string) (types.Function, bool) {
	f, ok := p.functions[name]
	return f, ok
}

func (p ParserEnvironment) ResolveEnum(name string) (types.Val, bool) {
	return types.NewErrorVal(fmt.Errorf("no such enum: %s", name)), false
}
