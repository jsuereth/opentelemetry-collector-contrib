// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanType = NewStructureType("span", map[string]Type{
	"name":                     StringType,
	"start_time_unix_nano":     IntType,
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
	// TODO - we want to return *writable* values for some of these.
	switch field {
	case "trace_id":
	case "span_id":
	case "trace_state":
	case "parent_span_id":
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
				return err
			},
		)
	case "kind":
	case "start_time_unix_nano":
		return NewIntVal(ptrace.Span(s).StartTimestamp().AsTime().UnixNano())
	case "end_time_unix_nano":
		return NewIntVal(ptrace.Span(s).EndTimestamp().AsTime().UnixNano())
	case "start_time":
	case "end_time":
	case "attributes":
	case "dropped_attributes_count":
		return NewIntVal(int64(ptrace.Span(s).DroppedAttributesCount()))
	case "events":
	case "dropped_events_count":
		return NewIntVal(int64(ptrace.Span(s).DroppedEventsCount()))
	case "links":
	case "dropped_links_count":
		return NewIntVal(int64(ptrace.Span(s).DroppedLinksCount()))
	case "status":
	}
	return NewErrorVal(fmt.Errorf("unknown field on span: %s", field))
}

func NewSpanVal(s ptrace.Span) Val {
	return spanVal(s)
}
