// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - more complex type?
var SliceType = runtime.NewPrimitiveType("pcommon.Slice")

type sliceVal pcommon.Slice

// ConvertTo implements types.Var.
func (s sliceVal) ConvertTo(runtime.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements types.Var.
func (s sliceVal) SetValue(v runtime.Val) error {
	if v.Type() == SliceType {
		return pcommon.Slice(s).FromRaw(v.Value().(pcommon.Slice).AsRaw())
	}
	return fmt.Errorf("unimplemented conversion %v to pcommon.Slice", v.Type())
}

// Type implements types.Var.
func (s sliceVal) Type() runtime.Type {
	return SliceType
}

// Value implements types.Var.
func (s sliceVal) Value() any {
	return pcommon.Slice(s)
}

// GetIndex implements traits.Indexable
func (s sliceVal) GetIndex(index int64) runtime.Val {
	idx := int(index)
	if idx <= 0 || idx >= pcommon.Slice(s).Len() {
		return NewErrorVal(fmt.Errorf("index %d out of bounds", idx))
	}
	return NewPvalVar(pcommon.Slice(s).At(idx))
}

func NewSliceVar(s pcommon.Slice) runtime.Var {
	return sliceVal(s)
}
