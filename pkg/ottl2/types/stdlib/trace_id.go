// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var TraceIDType = types.NewStructureType("TraceID",
	map[string]types.Type{
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
func (t traceIDVar) SetValue(v types.Val) error {
	id, ok := v.Value().(pcommon.TraceID)
	if ok {
		t.setter(id)
	}
	return nil
}

// Type implements Var.
func (t traceIDVar) Type() types.Type {
	return TraceIDType
}

// Value implements Var.
func (t traceIDVar) Value() any {
	return t.getter()
}

// traceIdVal is StructureAccessible
func (t traceIDVar) GetField(field string) types.Val {
	if field == "string" {
		return NewStringVar(
			func() string {
				return t.getter().String()
			},
			func(s string) {
				id, err := parseTraceID(s)
				if err == nil {
					t.setter(id)
				}
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field: %s", field))
}

func NewTraceIdVar(getter func() pcommon.TraceID,
	setter func(pcommon.TraceID)) types.Var {
	return traceIDVar{getter, setter}
}

func NewTraceIdVal(id pcommon.TraceID) types.Val {
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
