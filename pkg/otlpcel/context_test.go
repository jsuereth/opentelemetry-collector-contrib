// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package otlpcel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func Test_celCompilation(t *testing.T) {
	tests := []struct {
		name       string
		span       ptrace.Span
		expression string
		expect     func(t *testing.T, span ptrace.Span, out any)
	}{
		{
			name: "test grab name",
			span: func() ptrace.Span {
				s := ptrace.NewSpan()
				s.SetName("test")
				return s
			}(),
			expression: "span.name",
			expect: func(t *testing.T, span ptrace.Span, out any) {
				assert.Equal(t, "test", out)
			},
		},
		{
			name: "test set name",
			span: func() ptrace.Span {
				s := ptrace.NewSpan()
				s.SetName("test")
				return s
			}(),
			expression: `span.SetName("new-value")`,
			expect: func(t *testing.T, span ptrace.Span, out any) {
				assert.Equal(t, true, out)
				assert.Equal(t, "new-value", span.Name())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := NewSpanEnv()
			assert.NoError(t, err)
			ast, iss := env.Compile(tt.expression)
			assert.NoError(t, iss.Err())
			prg, err := env.Program(ast)
			assert.NoError(t, err)
			span := tt.span
			activation, err := NewSpanActivation(span)
			assert.NoError(t, err)
			// Note: third parameters is cost details so we can evaluate those.
			out, _, err := prg.Eval(activation)
			assert.NoError(t, err)
			tt.expect(t, span, out.Value())
		})
	}
}
