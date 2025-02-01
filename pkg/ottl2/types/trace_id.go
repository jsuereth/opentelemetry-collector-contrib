// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

var TraceIDType = NewStructureType("TraceID",
	map[string]Type{
		"string": StringType,
	})

type traceIDVar struct {
	getter func() pcommon.TraceID
	setter func(pcommon.TraceID)
}

// ConvertTo implements Var.
func (t traceIDVar) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements Var.
func (t traceIDVar) SetValue(v Val) error {
	id, ok := v.Value().(pcommon.TraceID)
	if ok {
		t.setter(id)
	}
	return nil
}

// Type implements Var.
func (t traceIDVar) Type() Type {
	return TraceIDType
}

// Value implements Var.
func (t traceIDVar) Value() any {
	return t.getter()
}

// traceIdVal is StructureAccessible
func (t traceIDVar) GetField(field string) Val {
	if field == "string" {
		return NewGetterSetterVar(
			StringType,
			func() Val {
				return NewStringVal(t.getter().String())
			},
			func(v Val) error {
				if str, ok := v.Value().(string); ok {
					id, err := parseTraceID(str)
					if err != nil {
						return err
					}
					t.setter(id)
					return nil
				}
				return fmt.Errorf("{trace_id}.string must take string value: %v", v)
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field: %s", field))
}

func NewTraceIdVar(getter func() pcommon.TraceID,
	setter func(pcommon.TraceID)) Var {
	return traceIDVar{getter, setter}
}

func NewTraceIdVal(id pcommon.TraceID) Val {
	v := &id
	return NewTraceIdVar(
		func() pcommon.TraceID { return *v },
		func(ti pcommon.TraceID) {
			*v = ti
		},
	)
}

func parseTraceID(traceIDStr string) (pcommon.TraceID, error) {
	var id pcommon.TraceID
	if hex.DecodedLen(len(traceIDStr)) != len(id) {
		return pcommon.TraceID{}, errors.New("trace ids must be 32 hex characters")
	}
	_, err := hex.Decode(id[:], []byte(traceIDStr))
	if err != nil {
		return pcommon.TraceID{}, err
	}
	return id, nil
}
