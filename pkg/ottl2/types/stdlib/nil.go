// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

// TODO - Nil should be a special type.
var NilType = types.NewPrimitiveType("nil")

type nilVal struct{}

// Type implements Val.
func (n *nilVal) Type() types.Type {
	return NilType
}

// ConvertTo implements Val.
func (n *nilVal) ConvertTo(t types.Type) (any, error) {
	switch t {
	case BoolType:
		return false, nil
	}
	// If the type conversion isn't supported return an error.
	return nil, fmt.Errorf("type conversion error from 'nil' to '%v'", t.Name())
}

func (n *nilVal) Value() any {
	return nil
}

var NilVal types.Val = &nilVal{}
