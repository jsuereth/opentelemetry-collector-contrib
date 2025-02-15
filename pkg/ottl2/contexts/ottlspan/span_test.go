// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottlspan // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/contexts/ottlspan"

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func Test_ottlExpressions(t *testing.T) {
	tests := []struct {
		name    string
		ottl    string
		context func() SpanContext
		expect  func(t *testing.T, ctx SpanContext, result any, cond bool, err error)
	}{
		{
			name: "cache",
			ottl: "set(cache[\"key\"], \"test\")",
			context: func() SpanContext {
				return NewSpanContext(
					ptrace.NewSpan(),
					pcommon.NewInstrumentationScope(),
					pcommon.NewResource(),
					ptrace.NewScopeSpans(),
					ptrace.NewResourceSpans(),
				)
			},
			expect: func(t *testing.T, ctx SpanContext, result any, cond bool, err error) {
				assert.Nil(t, err)
				assert.True(t, cond)
				v, ok := ctx.cache.Get("key")
				assert.True(t, ok, "key not found in cache!")
				assert.Equal(t, pcommon.ValueTypeStr, v.Type())
				assert.Equal(t, "test", v.AsString())
			},
		},
		{
			name: "span name",
			ottl: "set(span.name, \"cat\")",
			context: func() SpanContext {
				return NewSpanContext(
					ptrace.NewSpan(),
					pcommon.NewInstrumentationScope(),
					pcommon.NewResource(),
					ptrace.NewScopeSpans(),
					ptrace.NewResourceSpans(),
				)
			},
			expect: func(t *testing.T, ctx SpanContext, result any, cond bool, err error) {
				assert.Nil(t, err)
				assert.True(t, cond)
				assert.Equal(t, "cat", ctx.span.Name())
			},
		},
		{
			name: "resource attribute",
			ottl: "set(resource.attributes[\"animal\"], \"cat\")",
			context: func() SpanContext {
				return NewSpanContext(
					ptrace.NewSpan(),
					pcommon.NewInstrumentationScope(),
					pcommon.NewResource(),
					ptrace.NewScopeSpans(),
					ptrace.NewResourceSpans(),
				)
			},
			expect: func(t *testing.T, ctx SpanContext, result any, cond bool, err error) {
				assert.Nil(t, err)
				assert.True(t, cond)
				v, ok := ctx.resource.Attributes().Get("animal")
				assert.True(t, ok, "resource attribute was not set")
				assert.Equal(t, "cat", v.AsString())
			},
		},
	}
	env := NewSpanTransformContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.context()
			var stmt *ottl2.Statement[SpanContext]
			t.Run("parse", func(t *testing.T) {
				s, err := ottl2.ParseStatement(env, tt.ottl)
				assert.Nil(t, err)
				stmt = &s
			})
			// If compile succesful, run rest of the test.
			if stmt != nil {
				t.Run("eval", func(t *testing.T) {
					result, cond, err := stmt.Execute(context.Background(), &ctx)
					if tt.expect != nil {
						tt.expect(t, ctx, result, cond, err)
					}
				})
			}
		})
	}
}
