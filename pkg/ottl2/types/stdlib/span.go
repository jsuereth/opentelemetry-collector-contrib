// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanType = types.NewStructureType("span", map[string]types.Type{
	"name":                     StringType,
	"trace_id":                 TraceIDType,
	"start_time":               TimeType,
	"start_time_unix_nano":     IntType,
	"end_time":                 TimeType,
	"end_time_unix_nano":       IntType,
	"dropped_attributes_count": IntType,
	"dropped_events_count":     IntType,
	"dropped_links_count":      IntType,
})

type spanVal ptrace.Span

// Type implements Val.
func (s spanVal) Type() types.Type {
	return SpanType
}

func (s spanVal) ConvertTo(t types.Type) (any, error) {
	return nil, fmt.Errorf("not implemented for span")
}

func (s spanVal) Value() any {
	return ptrace.Span(s)
}

// Span has members
func (s spanVal) GetField(field string) types.Val {
	switch field {
	case "trace_id":
		return NewTraceIdVar(
			func() pcommon.TraceID {
				return ptrace.Span(s).TraceID()
			},
			func(v pcommon.TraceID) {
				ptrace.Span(s).SetTraceID(v)
			},
		)
	case "span_id":
		return NewSpanIDVar(
			func() pcommon.SpanID {
				return ptrace.Span(s).SpanID()
			},
			func(v pcommon.SpanID) {
				ptrace.Span(s).SetSpanID(v)
			},
		)
	case "trace_state":
	case "parent_span_id":
		return NewSpanIDVar(
			func() pcommon.SpanID {
				return ptrace.Span(s).ParentSpanID()
			},
			func(si pcommon.SpanID) {
				ptrace.Span(s).SetParentSpanID(si)
			},
		)
	case "name":
		return NewStringVar(
			func() string {
				return ptrace.Span(s).Name()
			},
			func(v string) {
				ptrace.Span(s).SetName(v)
			},
		)
	case "kind":
	case "start_time_unix_nano":
		return NewIntVar(
			func() int64 {
				return ptrace.Span(s).StartTimestamp().AsTime().UnixNano()
			},
			func(t int64) {
				ptrace.Span(s).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, t)))
			},
		)
	case "end_time_unix_nano":
		return NewIntVar(
			func() int64 {
				return ptrace.Span(s).EndTimestamp().AsTime().UnixNano()
			},
			func(t int64) {
				ptrace.Span(s).SetEndTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, t)))
			},
		)
	case "start_time":
		return NewTimeVar(
			func() time.Time {
				return ptrace.Span(s).StartTimestamp().AsTime()
			},
			func(t time.Time) {
				ptrace.Span(s).SetStartTimestamp(pcommon.NewTimestampFromTime(t))
			},
		)
	case "end_time":
		return NewTimeVar(
			func() time.Time {
				return ptrace.Span(s).EndTimestamp().AsTime()
			},
			func(t time.Time) {
				ptrace.Span(s).SetEndTimestamp(pcommon.NewTimestampFromTime(t))
			},
		)
	case "attributes":
		return NewPmapVar(ptrace.Span(s).Attributes())
	case "dropped_attributes_count":
		return NewIntVar(
			func() int64 { return int64(ptrace.Span(s).DroppedAttributesCount()) },
			func(v int64) { ptrace.Span(s).SetDroppedAttributesCount(uint32(v)) },
		)
	case "events":
	case "dropped_events_count":
		return NewIntVar(
			func() int64 { return int64(ptrace.Span(s).DroppedEventsCount()) },
			func(v int64) { ptrace.Span(s).SetDroppedEventsCount(uint32(v)) },
		)
	case "links":
	case "dropped_links_count":
		return NewIntVar(
			func() int64 { return int64(ptrace.Span(s).DroppedLinksCount()) },
			func(v int64) { ptrace.Span(s).SetDroppedLinksCount(uint32(v)) },
		)
	case "status":
	}
	return NewErrorVal(fmt.Errorf("unknown field on span: %s", field))
}

func NewSpanVal(s ptrace.Span) types.Val {
	return spanVal(s)
}
