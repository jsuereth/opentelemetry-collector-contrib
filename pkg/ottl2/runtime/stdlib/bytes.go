// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
)

// TODO - treat bytes as an array type?
var ByteSliceType = runtime.NewPrimitiveType("bytes")

type byteSliceVal struct {
	value []byte
}

// Type implements Val.
func (i *byteSliceVal) Type() runtime.Type {
	return ByteSliceType
}

// ConvertTo implements Val.
func (i *byteSliceVal) ConvertTo(t runtime.Type) (any, error) {
	panic("unimplemented")
}

func (i *byteSliceVal) Value() any {
	return ([]byte)(i.value)
}

func NewByteSliceVal(v []byte) runtime.Val {
	return &byteSliceVal{v}
}
