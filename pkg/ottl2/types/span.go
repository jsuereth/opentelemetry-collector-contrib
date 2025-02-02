// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanType = NewStructureType("span", map[string]Type{
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
func (s spanVal) Type() Type {
	return SpanType
}

func (s spanVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("not implemented for span")
}

func (s spanVal) Value() any {
	return ptrace.Span(s)
}

// Span has members
func (s spanVal) GetField(field string) Val {
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
		return NewGetterSetterVar(
			StringType,
			func() Val {
				return NewStringVal(ptrace.Span(s).Name())
			},
			func(v Val) error {
				n, err := v.ConvertTo(reflect.TypeFor[string]())
				if err != nil {
					return err
				}
				ptrace.Span(s).SetName(n.(string))
				return nil
			},
		)
	case "kind":
	case "start_time_unix_nano":
		return NewGetterSetterVar(
			IntType,
			func() Val {
				return NewIntVal(ptrace.Span(s).StartTimestamp().AsTime().UnixNano())
			},
			func(v Val) error {
				if t, ok := v.Value().(int64); ok {
					ptrace.Span(s).SetStartTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, t)))
				}
				return nil
			},
		)
	case "end_time_unix_nano":
		return NewGetterSetterVar(
			IntType,
			func() Val {
				return NewIntVal(ptrace.Span(s).EndTimestamp().AsTime().UnixNano())
			},
			func(v Val) error {
				if t, ok := v.Value().(int64); ok {
					ptrace.Span(s).SetEndTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, t)))
				}
				return nil
			},
		)
	case "start_time":
		return NewGetterSetterVar(
			TimeType,
			func() Val {
				return NewTimeVal(ptrace.Span(s).StartTimestamp().AsTime())
			},
			func(v Val) error {
				n := pcommon.NewTimestampFromTime(v.Value().(time.Time))
				ptrace.Span(s).SetStartTimestamp(n)
				return nil
			},
		)
	case "end_time":
		return NewGetterSetterVar(
			TimeType,
			func() Val {
				return NewTimeVal(ptrace.Span(s).EndTimestamp().AsTime())
			},
			func(v Val) error {
				n := pcommon.NewTimestampFromTime(v.Value().(time.Time))
				ptrace.Span(s).SetEndTimestamp(n)
				return nil
			},
		)
	case "attributes":
		return NewPmapVar(ptrace.Span(s).Attributes())
	case "dropped_attributes_count":
		// TODO - Getter/Settter
		return NewIntVal(int64(ptrace.Span(s).DroppedAttributesCount()))
	case "events":
	case "dropped_events_count":
		// TODO - Getter/Settter
		return NewIntVal(int64(ptrace.Span(s).DroppedEventsCount()))
	case "links":
	case "dropped_links_count":
		// TODO - Getter/Settter
		return NewIntVal(int64(ptrace.Span(s).DroppedLinksCount()))
	case "status":
	}
	return NewErrorVal(fmt.Errorf("unknown field on span: %s", field))
}

func NewSpanVal(s ptrace.Span) Val {
	return spanVal(s)
}
