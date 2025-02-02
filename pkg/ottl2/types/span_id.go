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

var SpanIDType = NewStructureType("SpanID",
	map[string]Type{
		"string": StringType,
	})

type spanIDVar struct {
	getter func() pcommon.SpanID
	setter func(pcommon.SpanID)
}

// ConvertTo implements Var.
func (s spanIDVar) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements Var.
func (s spanIDVar) SetValue(v Val) error {
	id, ok := v.Value().(pcommon.SpanID)
	if ok {
		s.setter(id)
		return nil
	}
	return fmt.Errorf("unable to set span id to %v", v)
}

// Type implements Var.
func (s spanIDVar) Type() Type {
	return SpanIDType
}

// Value implements Var.
func (s spanIDVar) Value() any {
	return s.getter()
}

// SpanID is StructureAccessible.
func (t spanIDVar) GetField(field string) Val {
	if field == "string" {
		return NewStringVar(
			func() string {
				return t.getter().String()
			},
			func(s string) {
				id, err := parseSpanID(s)
				if err == nil {
					t.setter(id)
				}
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field: %s", field))
}

func NewSpanIDVar(getter func() pcommon.SpanID,
	setter func(pcommon.SpanID)) Var {
	return spanIDVar{getter, setter}
}

func NewSpanIDVal(id pcommon.SpanID) Val {
	v := &id
	return NewSpanIDVar(
		func() pcommon.SpanID {
			return *v
		},
		func(s pcommon.SpanID) {
			*v = s
		},
	)
}

func parseSpanID(spanIDStr string) (pcommon.SpanID, error) {
	var id pcommon.SpanID
	if hex.DecodedLen(len(spanIDStr)) != len(id) {
		return pcommon.SpanID{}, errors.New("span ids must be 16 hex characters")
	}
	_, err := hex.Decode(id[:], []byte(spanIDStr))
	if err != nil {
		return pcommon.SpanID{}, err
	}
	return id, nil
}
