// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types" // Context to run transformations

type TransformContext[E any] struct {
	pCtx      ParserEnvironment
	constants map[string]types.Val
}

func NewTransformContext[E any](ctxType any, opts ...Option[E]) TransformContext[E] {
	// TODO - take the "todo" type for E and use it.
	result := TransformContext[E]{
		pCtx:      NewParserEnvironemnt(map[string]types.Type{}, map[string]types.Function{}),
		constants: map[string]types.Val{},
	}
	for _, opt := range opts {
		opt.f(&result)
	}
	return result
}

type Option[E any] struct {
	f func(c *TransformContext[E])
}

// Registers a function (editor or convertor) usable in this transformation.
func WithFunction[E any](name string, f types.Function) Option[E] {
	return Option[E]{
		func(c *TransformContext[E]) {
			c.pCtx.functions[name] = f
		},
	}
}

// Registers a constant usable at "root" for this transformation.
func WithConstant[E any](name string, value types.Val) Option[E] {
	return Option[E]{
		func(c *TransformContext[E]) {
			c.pCtx.variables[name] = value.Type()
			c.constants[name] = value
		},
	}
}

// Registers a complete enumerated type.
func WithEnum[E any](enumDef EnumDefinition) Option[E] {
	return Option[E]{
		func(c *TransformContext[E]) {
			c.pCtx.enums = append(c.pCtx.enums, enumDef)
		},
	}
}

// TransformContext is a ParserContext.HasName
func (e TransformContext[E]) ResolveName(name string) (types.Type, bool) {
	return e.pCtx.ResolveName(name)
}

// TransformContext is a ParserContext.ResolveFunctoin
func (e TransformContext[E]) ResolveFunction(name string) (types.Function, bool) {
	return e.pCtx.ResolveFunction(name)
}

// TransformContext is a ParserContext.ResolveEnum
func (e TransformContext[E]) ResolveEnum(name string) (types.Val, bool) {
	return e.pCtx.ResolveEnum(name)
}

// Constructs an evaluation context for E.
func (e TransformContext[E]) NewEvalContext(ctx E) TransformEnvironment {
	panic("Unimplemented")
}
