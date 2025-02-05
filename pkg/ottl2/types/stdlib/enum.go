// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
)

type enumType string

func (e enumType) Name() string {
	return string(e)
}

func EnumType(e types.EnumProvider) types.Type {
	return enumType(e.TypeName())
}

type enumVal struct {
	v        int64
	provider types.EnumProvider
}

// ConvertTo implements Val.
func (e enumVal) ConvertTo(t types.Type) (any, error) {
	switch t {
	case enumType(t.Name()):
		return e.v, nil
	case IntType:
		return e.v, nil
	}
	return nil, fmt.Errorf("cannot convert enum %s to %s", e.provider.TypeName(), t.Name())
}

// Type implements Val.
func (e enumVal) Type() types.Type {
	return EnumType(e.provider)
}

// Value implements Val.
func (e enumVal) Value() any {
	return e.v
}

func NewEnumVal(id int64, e types.EnumProvider) types.Val {
	return enumVal{id, e}
}
