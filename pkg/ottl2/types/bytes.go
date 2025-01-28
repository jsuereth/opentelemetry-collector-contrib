// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

import "reflect"

// TODO - treat bytes as an array type?
var ByteSliceType = NewPrimitiveType("bytes")

type byteSliceVal struct {
	value []byte
}

// Type implements Val.
func (i *byteSliceVal) Type() Type {
	return ByteSliceType
}

// ConvertTo implements Val.
func (i *byteSliceVal) ConvertTo(typeDesc reflect.Type) (any, error) {
	panic("unimplemented")
}

func (i *byteSliceVal) Value() any {
	return ([]byte)(i.value)
}

func NewByteSliceVal(v []byte) Val {
	return &byteSliceVal{v}
}
