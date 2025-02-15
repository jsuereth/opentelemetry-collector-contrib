// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanType = runtime.NewStructureType("span", map[string]runtime.Type{
	"name":                     StringType,
	"span_id":                  SpanIDType,
	"trace_id":                 TraceIDType,
	"parent_span_id":           SpanIDType,
	"trace_state":              TraceStateType,
	"start_time":               TimeType,
	"start_time_unix_nano":     IntType,
	"end_time":                 TimeType,
	"end_time_unix_nano":       IntType,
	"dropped_attributes_count": IntType,
	"dropped_events_count":     IntType,
	"dropped_links_count":      IntType,
	"status":                   SpanStatusType,
	"kind":                     SpanKindType,
	"attributes":               PmapType,
	"links":                    SpanLinkSliceType,
	// "events":
})

type spanVal ptrace.Span

// Type implements Val.
func (s spanVal) Type() runtime.Type {
	return SpanType
}

func (s spanVal) ConvertTo(t runtime.Type) (any, error) {
	return nil, fmt.Errorf("not implemented for span")
}

func (s spanVal) Value() any {
	return ptrace.Span(s)
}

// Span has members
func (s spanVal) GetField(field string) runtime.Val {
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
		return NewTraceStateVar(
			func() string {
				return ptrace.Span(s).TraceState().AsRaw()
			},
			func(ts string) {
				ptrace.Span(s).TraceState().FromRaw(ts)
			},
		)
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
		return NewSpanKindVar(
			func() ptrace.SpanKind {
				return ptrace.Span(s).Kind()
			},
			func(sk ptrace.SpanKind) {
				ptrace.Span(s).SetKind(sk)
			},
		)
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
		return NewSpanLinkSliceVar(ptrace.Span(s).Links())
	case "dropped_links_count":
		return NewIntVar(
			func() int64 { return int64(ptrace.Span(s).DroppedLinksCount()) },
			func(v int64) { ptrace.Span(s).SetDroppedLinksCount(uint32(v)) },
		)
	case "status":
		return NewSpanStatusVar(
			func() ptrace.Status {
				return ptrace.Span(s).Status()
			},
			func(st ptrace.Status) {
				st.CopyTo(ptrace.Span(s).Status())
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field on span: %s", field))
}

func NewSpanVal(s ptrace.Span) runtime.Val {
	return spanVal(s)
}
