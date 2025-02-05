// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var (
	traceID  = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	traceID2 = [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	spanID   = [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	spanID2  = [8]byte{8, 7, 6, 5, 4, 3, 2, 1}
)

func TestSpanFields(t *testing.T) {
	refSpan := createSpan()
	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")

	tests := []struct {
		name     string
		path     []testPath
		orig     any
		newVal   types.Val
		expected any
		expect   func(*testing.T, ptrace.Span)
	}{
		{
			name:     "trace_id",
			path:     []testPath{fieldPath("trace_id")},
			orig:     pcommon.TraceID(traceID),
			newVal:   NewTraceIdVal(pcommon.TraceID(traceID2)),
			expected: pcommon.TraceID(traceID2),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, s.TraceID(), pcommon.TraceID(traceID2))
			},
		},
		{
			name:     "span_id",
			path:     []testPath{fieldPath("span_id")},
			orig:     pcommon.SpanID(spanID),
			newVal:   NewSpanIDVal(pcommon.SpanID(spanID2)),
			expected: pcommon.SpanID(spanID2),
		},
		{
			name:     "trace_id string",
			path:     []testPath{fieldPath("trace_id"), fieldPath("string")},
			orig:     hex.EncodeToString(traceID[:]),
			newVal:   NewStringVal(hex.EncodeToString(traceID2[:])),
			expected: hex.EncodeToString(traceID2[:]),
		},
		{
			name:   "span_id string",
			path:   []testPath{fieldPath("span_id"), fieldPath("string")},
			orig:   hex.EncodeToString(spanID[:]),
			newVal: NewStringVal(hex.EncodeToString(spanID2[:])),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, hex.EncodeToString(spanID2[:]), s.SpanID().String())
			},
		},
		{
			name:   "trace_state",
			path:   []testPath{fieldPath("trace_state")},
			orig:   "key1=val1,key2=val2",
			newVal: NewStringVal("key=newVal"),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, "key=newVal", s.TraceState().AsRaw())
			},
		},
		{
			name: "trace_state key",
			path: []testPath{
				fieldPath("trace_state"),
				keyPath("key1"),
			},
			orig:   "val1",
			newVal: NewStringVal("newVal"),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, "key1=newVal,key2=val2", s.TraceState().AsRaw())
			},
		},
		{
			name:   "parent_span_id",
			path:   []testPath{fieldPath("parent_span_id")},
			orig:   pcommon.SpanID(spanID2),
			newVal: NewSpanIDVal(pcommon.SpanID(spanID)),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, pcommon.SpanID(spanID), s.ParentSpanID())
			},
		},
		{
			name:   "parent_span_id string",
			path:   []testPath{fieldPath("parent_span_id"), fieldPath("string")},
			orig:   hex.EncodeToString(spanID2[:]),
			newVal: NewStringVal(hex.EncodeToString(spanID[:])),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, hex.EncodeToString(spanID[:]), s.ParentSpanID().String())
			},
		},
		{
			name:     "name",
			path:     []testPath{fieldPath("name")},
			orig:     "bear",
			newVal:   NewStringVal("cat"),
			expected: "cat",
		},
		// {
		// 	name:     "kind",
		// 	path:     []testPath{"kind"},
		// 	orig:     int64(2),
		// 	newVal:   NewIntVal(int64(3)),
		// 	expected: ptrace.SpanKindClient,
		// },
		// {
		// 	name:     "string kind",
		// 	path:     []testPath{"kind", "string"},
		// 	orig:     "Server",
		// 	newVal:   NewStringVal("Client"),
		// 	expected: ptrace.SpanKindClient,
		// },
		// {
		// 	name:     "deprecated string kind",
		// 	path:     []testPath{"kind", "deprecated_string"},
		// 	orig:     "SPAN_KIND_SERVER",
		// 	newVal:   NewStringVal("SPAN_KIND_CLIENT"),
		// 	expected: ptrace.SpanKindClient,
		// },
		{
			name:   "start_time_unix_nano",
			path:   []testPath{fieldPath("start_time_unix_nano")},
			orig:   int64(100_000_000),
			newVal: NewIntVal(int64(200_000_000)),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, time.Unix(0, 200_000_000).UTC(), s.StartTimestamp().AsTime())
			},
		},
		{
			name:   "end_time_unix_nano",
			path:   []testPath{fieldPath("end_time_unix_nano")},
			orig:   int64(500_000_000),
			newVal: NewIntVal(int64(200_000_000)),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, time.Unix(0, 200_000_000).UTC(), s.EndTimestamp().AsTime())
			},
		},
		{
			name:   "attributes",
			path:   []testPath{fieldPath("attributes")},
			orig:   refSpan.Attributes(),
			newVal: NewPmapVar(newAttrs),
			expect: func(t *testing.T, s ptrace.Span) {
				// TODO - define what to check.
			},
		},
		{
			name: "attributes string",
			path: []testPath{
				fieldPath("attributes"),
				keyPath("str"),
			},
			orig:   "val",
			newVal: NewStringVal("newVal"),
			expect: func(t *testing.T, s ptrace.Span) {
				v, ok := s.Attributes().Get("str")
				assert.True(t, ok)
				assert.Equal(t, "newVal", v.AsString())
			},
		},
		{
			name: "attributes bool",
			path: []testPath{
				fieldPath("attributes"),
				keyPath("bool"),
			},
			orig:   true,
			newVal: NewBoolVal(false),
			expect: func(t *testing.T, s ptrace.Span) {
				v, ok := s.Attributes().Get("bool")
				assert.True(t, ok)
				assert.False(t, v.Bool())
			},
		},
		{
			name: "attributes int",
			path: []testPath{
				fieldPath("attributes"),
				keyPath("int"),
			},
			orig:   int64(10),
			newVal: NewIntVal(20),
			expect: func(t *testing.T, s ptrace.Span) {
				v, ok := s.Attributes().Get("int")
				assert.True(t, ok)
				assert.Equal(t, int64(20), v.Int())
			},
		},
		{
			name: "attributes float",
			path: []testPath{
				fieldPath("attributes"),
				keyPath("double"),
			},
			orig:   float64(1.2),
			newVal: NewFloatVal(2.4),
			expect: func(t *testing.T, s ptrace.Span) {
				v, ok := s.Attributes().Get("double")
				assert.True(t, ok)
				assert.Equal(t, float64(2.4), v.Double())
			},
		},
		{
			name: "attributes bytes",
			path: []testPath{
				fieldPath("attributes"),
				keyPath("bytes"),
			},
			orig: func() any {
				v := pcommon.NewByteSlice()
				v.FromRaw([]byte{1, 3, 2})
				return v
			}(),
			newVal: NewByteSliceVal([]byte{2, 3, 4}),
			expect: func(t *testing.T, s ptrace.Span) {
				v, ok := s.Attributes().Get("bytes")
				assert.True(t, ok)
				assert.Equal(t, []byte{2, 3, 4}, v.Bytes().AsRaw())
			},
		},
		// {
		// 	name: "attributes array empty",
		// 	path: []testPath{
		// 		fieldPath("attributes"),
		// 		keyPath("arr_empty"),
		// 	},
		// 	orig:   pcommon.NewByteSlice(),
		// 	newVal: NewByteSliceVal([]byte{}),
		// 	expect: func(t *testing.T, s ptrace.Span) {
		// 		// no-op ?
		// 	},
		// },
		// {
		// 	name: "attributes array string",
		// 	path: []testPath{
		// 		fieldPath("attributes"),
		// 		keyPath("arr_str"),
		// 	},
		// 	orig:   pcommon.NewStringSlice(),
		// 	newVal: []string{"new"},
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("arr_str").AppendEmpty().SetStr("new")
		// 	},
		// },
		// {
		// 	name: "attributes array bool",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("arr_bool"),
		// 			},
		// 		},
		// 	},
		// 	orig: func() pcommon.Slice {
		// 		val, _ := refSpan.Attributes().Get("arr_bool")
		// 		return val.Slice()
		// 	}(),
		// 	newVal: []bool{false},
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("arr_bool").AppendEmpty().SetBool(false)
		// 	},
		// },
		// {
		// 	name: "attributes array int",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("arr_int"),
		// 			},
		// 		},
		// 	},
		// 	orig: func() pcommon.Slice {
		// 		val, _ := refSpan.Attributes().Get("arr_int")
		// 		return val.Slice()
		// 	}(),
		// 	newVal: []int64{20},
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("arr_int").AppendEmpty().SetInt(20)
		// 	},
		// },
		// {
		// 	name: "attributes array float",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("arr_float"),
		// 			},
		// 		},
		// 	},
		// 	orig: func() pcommon.Slice {
		// 		val, _ := refSpan.Attributes().Get("arr_float")
		// 		return val.Slice()
		// 	}(),
		// 	newVal: []float64{2.0},
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("arr_float").AppendEmpty().SetDouble(2.0)
		// 	},
		// },
		// {
		// 	name: "attributes array bytes",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("arr_bytes"),
		// 			},
		// 		},
		// 	},
		// 	orig: func() pcommon.Slice {
		// 		val, _ := refSpan.Attributes().Get("arr_bytes")
		// 		return val.Slice()
		// 	}(),
		// 	newVal: [][]byte{{9, 6, 4}},
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("arr_bytes").AppendEmpty().SetEmptyBytes().FromRaw([]byte{9, 6, 4})
		// 	},
		// },
		// {
		// 	name: "attributes nested",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("slice"),
		// 			},
		// 			&TestKey[*spanContext]{
		// 				I: ottltest.Intp(0),
		// 			},
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("map"),
		// 			},
		// 		},
		// 	},
		// 	orig: func() string {
		// 		val, _ := refSpan.Attributes().Get("slice")
		// 		val, _ = val.Slice().At(0).Map().Get("map")
		// 		return val.Str()
		// 	}(),
		// 	newVal: "new",
		// 	modified: func(span ptrace.Span) {
		// 		span.Attributes().PutEmptySlice("slice").AppendEmpty().SetEmptyMap().PutStr("map", "new")
		// 	},
		// },
		// {
		// 	name: "attributes nested new values",
		// 	path: &TestPath[*spanContext]{
		// 		N: "attributes",
		// 		KeySlice: []ottl.Key[*spanContext]{
		// 			&TestKey[*spanContext]{
		// 				S: ottltest.Strp("new"),
		// 			},
		// 			&TestKey[*spanContext]{
		// 				I: ottltest.Intp(2),
		// 			},
		// 			&TestKey[*spanContext]{
		// 				I: ottltest.Intp(0),
		// 			},
		// 		},
		// 	},
		// 	orig: func() any {
		// 		return nil
		// 	}(),
		// 	newVal: "new",
		// 	modified: func(span ptrace.Span) {
		// 		s := span.Attributes().PutEmptySlice("new")
		// 		s.AppendEmpty()
		// 		s.AppendEmpty()
		// 		s.AppendEmpty().SetEmptySlice().AppendEmpty().SetStr("new")
		// 	},
		// },
		{
			name:     "dropped_attributes_count",
			path:     []testPath{fieldPath("dropped_attributes_count")},
			orig:     int64(10),
			newVal:   NewIntVal(20),
			expected: int64(20),
		},
		// {
		// 	name:   "events",
		// 	path:   []string{"events"},
		// 	orig:   refSpan.Events(),
		// 	newVal: newEvents,
		// 	modified: func(span ptrace.Span) {
		// 		span.Events().RemoveIf(func(_ ptrace.SpanEvent) bool {
		// 			return true
		// 		})
		// 		newEvents.CopyTo(span.Events())
		// 	},
		// },
		{
			name:     "dropped_events_count",
			path:     []testPath{fieldPath("dropped_events_count")},
			orig:     int64(20),
			newVal:   NewIntVal(10),
			expected: int64(10),
		},
		// {
		// 	name:   "links",
		// 	path:   []string{"links"},
		// 	orig:   refSpan.Links(),
		// 	newVal: newLinks,
		// 	modified: func(span ptrace.Span) {
		// 		span.Links().RemoveIf(func(_ ptrace.SpanLink) bool {
		// 			return true
		// 		})
		// 		newLinks.CopyTo(span.Links())
		// 	},
		// },
		{
			name:     "dropped_links_count",
			path:     []testPath{fieldPath("dropped_links_count")},
			orig:     int64(30),
			newVal:   NewIntVal(40),
			expected: int64(40),
		},
		{
			name: "status",
			path: []testPath{fieldPath("status")},
			orig: refSpan.Status(),
			newVal: NewSpanStatusVar(func() ptrace.Status {
				s := ptrace.NewStatus()
				s.SetMessage("new status")
				return s
			}, func(ptrace.Status) {
				// ignore
			}),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, "new status", s.Status().Message())
			},
		},
		{
			name:   "status code",
			path:   []testPath{fieldPath("status"), fieldPath("code")},
			orig:   int64(ptrace.StatusCodeOk),
			newVal: NewStatusCodeVal(ptrace.StatusCodeError),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, ptrace.StatusCodeError, s.Status().Code())
			},
		},
		{
			name:   "status message",
			path:   []testPath{fieldPath("status"), fieldPath("message")},
			orig:   "good span",
			newVal: NewStringVal("bad span"),
			expect: func(t *testing.T, s ptrace.Span) {
				assert.Equal(t, "bad span", s.Status().Message())
			},
		},
		{
			name:     "start_time",
			path:     []testPath{fieldPath("start_time")},
			orig:     time.Date(1970, 1, 1, 0, 0, 0, 100000000, time.UTC),
			newVal:   NewTimeVal(time.Date(1970, 1, 1, 0, 0, 0, 200000000, time.UTC)),
			expected: time.UnixMilli(200).UTC(),
		},
		{
			name:     "end_time",
			path:     []testPath{fieldPath("end_time")},
			orig:     time.Date(1970, 1, 1, 0, 0, 0, 500000000, time.UTC),
			newVal:   NewTimeVal(time.Date(1970, 1, 1, 0, 0, 0, 200000000, time.UTC)),
			expected: time.UnixMilli(200).UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span := createSpan()
			v := lookupTestPath(NewSpanVal(span), tt.path)
			result := v.Value()
			assert.Equal(t, tt.orig, result)
			setter, ok := v.(types.Var)
			if !ok {
				assert.Fail(t, "path %v is not mutable", pathString(tt.path))
			}
			err := setter.SetValue(tt.newVal)
			assert.Nil(t, err)
			if tt.expected != nil {
				assert.Equal(t, tt.expected, v.Value(), "did not find expected %v on %v", tt.expected, v)
			}
			if tt.expect != nil {
				tt.expect(t, span)
			}
		})
	}
}

func createSpan() ptrace.Span {
	span := ptrace.NewSpan()
	span.SetTraceID(traceID)
	span.SetSpanID(spanID)
	span.TraceState().FromRaw("key1=val1,key2=val2")
	span.SetParentSpanID(spanID2)
	span.SetName("bear")
	span.SetKind(ptrace.SpanKindServer)
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(100)))
	span.SetEndTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(500)))
	span.Attributes().PutStr("str", "val")
	span.Attributes().PutBool("bool", true)
	span.Attributes().PutInt("int", 10)
	span.Attributes().PutDouble("double", 1.2)
	span.Attributes().PutEmptyBytes("bytes").FromRaw([]byte{1, 3, 2})

	span.Attributes().PutEmptySlice("arr_empty")

	arrStr := span.Attributes().PutEmptySlice("arr_str")
	arrStr.AppendEmpty().SetStr("one")
	arrStr.AppendEmpty().SetStr("two")

	arrBool := span.Attributes().PutEmptySlice("arr_bool")
	arrBool.AppendEmpty().SetBool(true)
	arrBool.AppendEmpty().SetBool(false)

	arrInt := span.Attributes().PutEmptySlice("arr_int")
	arrInt.AppendEmpty().SetInt(2)
	arrInt.AppendEmpty().SetInt(3)

	arrFloat := span.Attributes().PutEmptySlice("arr_float")
	arrFloat.AppendEmpty().SetDouble(1.0)
	arrFloat.AppendEmpty().SetDouble(2.0)

	arrBytes := span.Attributes().PutEmptySlice("arr_bytes")
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{1, 2, 3})
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{2, 3, 4})

	s := span.Attributes().PutEmptySlice("slice")
	s.AppendEmpty().SetEmptyMap().PutStr("map", "pass")

	span.SetDroppedAttributesCount(10)

	span.Events().AppendEmpty().SetName("event")
	span.SetDroppedEventsCount(20)

	span.Links().AppendEmpty().SetTraceID(traceID)
	span.SetDroppedLinksCount(30)

	span.Status().SetCode(ptrace.StatusCodeOk)
	span.Status().SetMessage("good span")

	return span
}
