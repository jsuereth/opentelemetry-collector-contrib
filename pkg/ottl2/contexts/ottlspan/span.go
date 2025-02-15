// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottlspan // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/contexts/ottlspan"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	// Experimental: *NOTE* this constant is subject to change or removal in the future.
	ContextName = "span"
)

var (
	// Experimental: *NOTE* this constant is subject to change or removal in the future.
	ContextType = runtime.NewStructureType(
		"SpanContext",
		map[string]runtime.Type{
			"cache":    stdlib.PmapType,
			"resource": stdlib.ResourceType,
			"span":     stdlib.SpanType,
		},
	)
)

type SpanContext struct {
	span                 ptrace.Span
	instrumentationScope pcommon.InstrumentationScope
	resource             pcommon.Resource
	cache                pcommon.Map
	scopeSpans           ptrace.ScopeSpans
	resourceSpans        ptrace.ResourceSpans
}

// Helper to make sure we pull resource and schema-url of resource from the right
// location.
type resourceContextWrapper SpanContext

func (s *resourceContextWrapper) SchemaUrl() string {
	return s.resourceSpans.SchemaUrl()
}

func (s *resourceContextWrapper) SetSchemaUrl(v string) {
	s.resourceSpans.SetSchemaUrl(v)
}

func (s *resourceContextWrapper) GetResource() pcommon.Resource {
	return s.resource
}

func (s *resourceContextWrapper) GetResourceSchemaURLItem() stdlib.SchemaURLItem {
	// TODO - implement this.
	return s
}

// GetField implements runtime.Structure.
func (s *SpanContext) GetField(field string) runtime.Val {
	switch field {
	case "span":
		return stdlib.NewSpanVal(s.span)
	case "scope":
		// TODO - implement this.
	case "resource":
		return stdlib.NewResourceVal((*resourceContextWrapper)(s))
	case "cache":
		return stdlib.NewPmapVar(s.cache)
	}
	return stdlib.NewErrorVal(fmt.Errorf("cannot find field on SpanContext: %s", field))
}

// ConvertTo implements runtime.Val.
func (s *SpanContext) ConvertTo(t runtime.Type) (any, error) {
	return nil, fmt.Errorf("cannot convert SpanContext to %s", t.Name())
}

// Type implements runtime.Val.
func (s *SpanContext) Type() runtime.Type {
	return ContextType
}

// Value implements runtime.Val.
func (s *SpanContext) Value() any {
	return s
}

// TODO - pull in internal/logging
// func (tCtx TransformContext) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
// 	err := encoder.AddObject("resource", logging.Resource(tCtx.resource))
// 	err = errors.Join(err, encoder.AddObject("scope", logging.InstrumentationScope(tCtx.instrumentationScope)))
// 	err = errors.Join(err, encoder.AddObject("span", logging.Span(tCtx.span)))
// 	err = errors.Join(err, encoder.AddObject("cache", logging.Map(tCtx.cache)))
// 	return err
// }

type SpanContextOption func(*SpanContext)

func NewSpanContext(span ptrace.Span, instrumentationScope pcommon.InstrumentationScope, resource pcommon.Resource, scopeSpans ptrace.ScopeSpans, resourceSpans ptrace.ResourceSpans, options ...SpanContextOption) SpanContext {
	tc := SpanContext{
		span:                 span,
		instrumentationScope: instrumentationScope,
		resource:             resource,
		cache:                pcommon.NewMap(),
		scopeSpans:           scopeSpans,
		resourceSpans:        resourceSpans,
	}
	for _, opt := range options {
		opt(&tc)
	}
	return tc
}

// Experimental: *NOTE* this option is subject to change or removal in the future.
func WithCache(cache *pcommon.Map) SpanContextOption {
	return func(p *SpanContext) {
		if cache != nil {
			p.cache = *cache
		}
	}
}

// Constructs a new Transform context for parsing and executing `span` context OTTL statements.
func NewSpanTransformContext() ottl2.TransformContext[SpanContext] {
	return ottl2.NewTransformContext(
		ContextType,
		func(v *SpanContext) runtime.Structure { return v },
		ottl2.WithFunctions[SpanContext](funcs.StandardFuncs()),
	)
}
