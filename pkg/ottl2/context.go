// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types" // Context to run transformations

type TransformContext[E any] struct {
	pCtx ParserEnvironment
}

// TransformContext is a ParserContext.HasName
func (e TransformContext[E]) HasName(name string) bool {
	return e.pCtx.HasName(name)
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
