// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

//revive:disable:var-naming The methods in this interface are defined by pdata types.
type SchemaURLItem interface {
	SchemaUrl() string
	SetSchemaUrl(v string)
}

//revive:enable:var-naming

// Allows interacting with Resource{signal} types.
type ResourceContext interface {
	GetResource() pcommon.Resource
	GetResourceSchemaURLItem() SchemaURLItem
}

var ResourceType types.Type = types.NewStructureType("pcommon.Resource", map[string]types.Type{
	"attributes":               PmapType,
	"dropped_attributes_count": IntType,
	"schema_url":               StringType,
})

type resourceVal struct {
	ctx ResourceContext
}

// ConvertTo implements types.Val.
func (r resourceVal) ConvertTo(t types.Type) (any, error) {
	if t == ResourceType {
		return r.ctx.GetResource(), nil
	}
	return nil, fmt.Errorf("cannot convert pcommon.Resource to %s", t.Name())
}

// Type implements types.Val.
func (r resourceVal) Type() types.Type {
	return ResourceType
}

// Value implements types.Val.
func (r resourceVal) Value() any {
	return r.ctx.GetResource()
}

// GetField implements StructureAccessible
func (r resourceVal) GetField(field string) types.Val {
	switch field {
	case "attributes":
		return NewPmapVar(r.ctx.GetResource().Attributes())
	case "dropped_attributes_count":
		return NewIntVar(
			func() int64 {
				return int64(r.ctx.GetResource().DroppedAttributesCount())
			},
			func(i int64) {
				r.ctx.GetResource().SetDroppedAttributesCount(uint32(i))
			},
		)
	case "schema_url":
		return NewStringVar(
			func() string {
				return r.ctx.GetResourceSchemaURLItem().SchemaUrl()
			},
			func(s string) {
				r.ctx.GetResourceSchemaURLItem().SetSchemaUrl(s)
			},
		)
	}
	return NewErrorVal(fmt.Errorf("unknown field %s on pcommon.Resource", field))
}

// TODO - should we take in schema urls separately?
func NewResourceVal(v ResourceContext) types.Val {
	return &resourceVal{v}
}
