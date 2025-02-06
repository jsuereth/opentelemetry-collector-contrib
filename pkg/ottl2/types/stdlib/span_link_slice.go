// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanLinkSliceType types.Type = types.NewPrimitiveType("SpanLinkSlice")

type spanLinkSliceVar ptrace.SpanLinkSlice

// ConvertTo implements types.Var.
func (s spanLinkSliceVar) ConvertTo(types.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements types.Var.
func (s spanLinkSliceVar) SetValue(v types.Val) error {
	if v.Type() != SpanLinkSliceType {
		return fmt.Errorf("unable to write ptrace.SpanLinkSlice from %s", v.Type().Name())
	}
	ps := ptrace.SpanLinkSlice(s)
	os := v.Value().(ptrace.SpanLinkSlice)
	os.CopyTo(ps)
	return nil
}

// Type implements types.Var.
func (s spanLinkSliceVar) Type() types.Type {
	return SpanLinkSliceType
}

// Value implements types.Var.
func (s spanLinkSliceVar) Value() any {
	return ptrace.SpanLinkSlice(s)
}

// GetIndex implements traits.Indexable
func (s spanLinkSliceVar) GetIndex(index int64) types.Val {
	idx := int(index)
	if idx <= 0 || idx >= ptrace.SpanLinkSlice(s).Len() {
		return NewErrorVal(fmt.Errorf("index %d out of bounds", idx))
	}
	return NewSpanLinkVar(ptrace.SpanLinkSlice(s).At(idx))
}

func NewSpanLinkSliceVar(ls ptrace.SpanLinkSlice) types.Var {
	return spanLinkSliceVar(ls)
}
