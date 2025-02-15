// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package funcs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/funcs"

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func Test_IsList(t *testing.T) {
	tests := []struct {
		name     string
		value    runtime.Val
		expected bool
	}{
		{
			name:     "map",
			value:    stdlib.NewPmapVar(pcommon.NewMap()),
			expected: false,
		},
		{
			name:     "ValueTypeMap",
			value:    stdlib.NewPvalVar(pcommon.NewValueMap()),
			expected: false,
		},
		{
			name:     "not map",
			value:    stdlib.NewStringVal("not a map"),
			expected: false,
		},
		{
			name:     "ValueTypeSlice",
			value:    stdlib.NewPvalVar(pcommon.NewValueSlice()),
			expected: true,
		},
		{
			name:     "nil",
			value:    stdlib.NilVal,
			expected: false,
		},
		// TODO - support other list types from OTTL
		// {
		// 	name:     "plog.LogRecordSlice",
		// 	value:    plog.NewLogRecordSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "plog.ResourceLogsSlice",
		// 	value:    plog.NewResourceLogsSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "plog.ScopeLogsSlice",
		// 	value:    plog.NewScopeLogsSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.ExemplarSlice",
		// 	value:    pmetric.NewExemplarSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.ExponentialHistogramDataPointSlice",
		// 	value:    pmetric.NewExponentialHistogramDataPointSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.HistogramDataPointSlice",
		// 	value:    pmetric.NewHistogramDataPointSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.MetricSlice",
		// 	value:    pmetric.NewMetricSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.NumberDataPointSlice",
		// 	value:    pmetric.NewNumberDataPointSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.ResourceMetricsSlice",
		// 	value:    pmetric.NewResourceMetricsSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.ScopeMetricsSlice",
		// 	value:    pmetric.NewScopeMetricsSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.SummaryDataPointSlice",
		// 	value:    pmetric.NewSummaryDataPointSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "pmetric.SummaryDataPointValueAtQuantileSlice",
		// 	value:    pmetric.NewSummaryDataPointValueAtQuantileSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "ptrace.ResourceSpansSlice",
		// 	value:    ptrace.NewResourceSpansSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "ptrace.ScopeSpansSlice",
		// 	value:    ptrace.NewScopeSpansSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "ptrace.SpanEventSlice",
		// 	value:    ptrace.NewSpanEventSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "ptrace.SpanLinkSlice",
		// 	value:    ptrace.NewSpanLinkSlice(),
		// 	expected: true,
		// },
		// {
		// 	name:     "ptrace.SpanSlice",
		// 	value:    ptrace.NewSpanSlice(),
		// 	expected: true,
		// },
		// TODO - Support primitive maps?
		// {
		// 	name:     "[]string",
		// 	value:    []string{},
		// 	expected: true,
		// },
		// {
		// 	name:     "[]bool",
		// 	value:    []bool{},
		// 	expected: true,
		// },
		// {
		// 	name:     "[]int64",
		// 	value:    []int64{},
		// 	expected: true,
		// },
		// {
		// 	name:     "[]float64",
		// 	value:    []float64{},
		// 	expected: true,
		// },
		// {
		// 	name:     "[][]byte",
		// 	value:    [][]byte{},
		// 	expected: true,
		// },
		// {
		// 	name:     "[]any",
		// 	value:    []any{},
		// 	expected: true,
		// },
	}
	isList := NewIsListFunc()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isList.Call([]runtime.Val{
				tt.value,
			})
			v, err := result.ConvertTo(stdlib.BoolType)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, v.(bool))
		})
	}
}
