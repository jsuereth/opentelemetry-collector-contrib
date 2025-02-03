// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

// TODO - treat bytes as an array type?
var ByteSliceType = types.NewPrimitiveType("bytes")

type byteSliceVal struct {
	value []byte
}

// Type implements Val.
func (i *byteSliceVal) Type() types.Type {
	return ByteSliceType
}

// ConvertTo implements Val.
func (i *byteSliceVal) ConvertTo(t types.Type) (any, error) {
	panic("unimplemented")
}

func (i *byteSliceVal) Value() any {
	return ([]byte)(i.value)
}

func NewByteSliceVal(v []byte) types.Val {
	return &byteSliceVal{v}
}
