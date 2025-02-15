// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var SpanLinkSliceType runtime.Type = runtime.NewPrimitiveType("SpanLinkSlice")

type spanLinkSliceVar ptrace.SpanLinkSlice

// ConvertTo implements types.Var.
func (s spanLinkSliceVar) ConvertTo(runtime.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements types.Var.
func (s spanLinkSliceVar) SetValue(v runtime.Val) error {
	if v.Type() != SpanLinkSliceType {
		return fmt.Errorf("unable to write ptrace.SpanLinkSlice from %s", v.Type().Name())
	}
	ps := ptrace.SpanLinkSlice(s)
	os := v.Value().(ptrace.SpanLinkSlice)
	os.CopyTo(ps)
	return nil
}

// Type implements types.Var.
func (s spanLinkSliceVar) Type() runtime.Type {
	return SpanLinkSliceType
}

// Value implements types.Var.
func (s spanLinkSliceVar) Value() any {
	return ptrace.SpanLinkSlice(s)
}

// GetIndex implements traits.Indexable
func (s spanLinkSliceVar) GetIndex(index int64) runtime.Val {
	idx := int(index)
	if idx <= 0 || idx >= ptrace.SpanLinkSlice(s).Len() {
		return NewErrorVal(fmt.Errorf("index %d out of bounds", idx))
	}
	return NewSpanLinkVar(ptrace.SpanLinkSlice(s).At(idx))
}

func NewSpanLinkSliceVar(ls ptrace.SpanLinkSlice) runtime.Var {
	return spanLinkSliceVar(ls)
}
