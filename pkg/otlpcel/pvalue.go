// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"reflect"

	"github.com/google/cel-go/common/types/ref"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type pvalueType struct{}

// HasTrait implements ref.Type.
func (p *pvalueType) HasTrait(trait int) bool {
	return false
}

// TypeName implements ref.Type.
func (p *pvalueType) TypeName() string {
	return "pcommon.Value"
}

var PValueType ref.Type = &pvalueType{}

type pvalueWrapper pcommon.Value

// ConvertToNative implements ref.Val.
func (p pvalueWrapper) ConvertToNative(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

// ConvertToType implements ref.Val.
func (p pvalueWrapper) ConvertToType(typeValue ref.Type) ref.Val {
	panic("unimplemented")
}

// Equal implements ref.Val.
func (p pvalueWrapper) Equal(other ref.Val) ref.Val {
	panic("unimplemented")
}

// Type implements ref.Val.
func (p pvalueWrapper) Type() ref.Type {
	return PValueType
}

// Value implements ref.Val.
func (p pvalueWrapper) Value() any {
	return pcommon.Value(p)
}
