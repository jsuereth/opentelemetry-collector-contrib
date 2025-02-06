// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// TODO - define a special pccomon.Map type
// We want this to denote it can be keyed
var PmapType = types.NewPrimitiveType("pcommon.Map")

type pmapVal pcommon.Map

func (m pmapVal) Type() types.Type {
	return PmapType
}

func (m pmapVal) ConvertTo(t types.Type) (any, error) {
	switch t {
	case PmapType:
		return m.Value(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion from 'pcommon.Map' to %s", t.Name())
}

func (m pmapVal) Value() any {
	return pcommon.Map(m)
}

// PmapVal is KeyAccessible
func (m pmapVal) GetKey(key string) types.Val {
	// Note: pcommon.Map and pcommon.Value are references that mutate the underlying value.
	// so we don't need to use a getter/setter appraoch.
	v, ok := pcommon.Map(m).Get(key)
	if !ok {
		return NewLazyPval(func() pcommon.Value {
			return pcommon.Map(m).PutEmpty(key)
		})
	}
	return NewPvalVar(v)
}

// SetValue implements Var.
func (m pmapVal) SetValue(o types.Val) error {
	other, err := o.ConvertTo(PmapType)
	if err != nil {
		return err
	}
	other.(pcommon.Map).MoveTo(pcommon.Map(m))
	return nil
}

func NewPmapVar(m pcommon.Map) types.Var {
	return pmapVal(m)
}
