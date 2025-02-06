// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - more complex type?
var SliceType = types.NewPrimitiveType("pcommon.Slice")

type sliceVal pcommon.Slice

// ConvertTo implements types.Var.
func (s sliceVal) ConvertTo(types.Type) (any, error) {
	panic("unimplemented")
}

// SetValue implements types.Var.
func (s sliceVal) SetValue(v types.Val) error {
	if v.Type() == SliceType {
		return pcommon.Slice(s).FromRaw(v.Value().(pcommon.Slice).AsRaw())
	}
	return fmt.Errorf("unimplemented conversion %v to pcommon.Slice", v.Type())
}

// Type implements types.Var.
func (s sliceVal) Type() types.Type {
	return SliceType
}

// Value implements types.Var.
func (s sliceVal) Value() any {
	return pcommon.Slice(s)
}

// GetIndex implements traits.Indexable
func (s sliceVal) GetIndex(index int64) types.Val {
	idx := int(index)
	if idx <= 0 || idx >= pcommon.Slice(s).Len() {
		return NewErrorVal(fmt.Errorf("index %d out of bounds", idx))
	}
	return NewPvalVar(pcommon.Slice(s).At(idx))
}

func NewSliceVar(s pcommon.Slice) types.Var {
	return sliceVal(s)
}
