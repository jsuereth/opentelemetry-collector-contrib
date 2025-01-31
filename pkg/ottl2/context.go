// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types" // Context to run transformations
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/traits"
)

type TransformContext[E any] struct {
	pCtx      ParserEnvironment
	constants map[string]types.Val
	converter func(E) types.Val
}

// TODO - see if we cna infer a types.StructType from a go interface...
func NewTransformContext[E any](
	ctxType types.StructType,
	converter func(E) types.Val, // Converts the raw context type into a `types.Val` with `traits.StructureAccessible` matching the `ctxType`.
	opts ...Option[E]) TransformContext[E] {
	contextFields := map[string]types.Type{}
	for _, field := range ctxType.FieldNames() {
		t, _ := ctxType.GetField(field)
		contextFields[field] = t
	}
	result := TransformContext[E]{
		pCtx:      NewParserEnvironemnt(contextFields, map[string]types.Function{}),
		constants: map[string]types.Val{},
		converter: converter,
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
func WithFunction[E any](f types.Function) Option[E] {
	return Option[E]{
		func(c *TransformContext[E]) {
			c.pCtx.functions[f.Name()] = f
		},
	}
}

// Registers a set of functions function (editor or convertor).
func WithFunctions[E any](fs []types.Function) Option[E] {
	return Option[E]{
		func(c *TransformContext[E]) {
			for _, f := range fs {
				c.pCtx.functions[f.Name()] = f
			}
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

type valDrivenEvalContext struct {
	source types.Val
}

// Parent implements EvalContext.
func (v *valDrivenEvalContext) Parent() EvalContext {
	return nil
}

func (v *valDrivenEvalContext) String() string {
	return fmt.Sprintf("Context{%v}", v.source)
}

// ResolveName implements EvalContext.
func (v *valDrivenEvalContext) ResolveName(name string) (types.Val, bool) {
	r := v.source.(traits.StructureAccessible).GetField(name)
	return r, r.Type() != types.ErrorType
}

// Constructs an evaluation context for E.
func (e TransformContext[E]) NewEvalContext(ctx E) EvalContext {
	v := e.converter(ctx)
	return &valDrivenEvalContext{v}
}
