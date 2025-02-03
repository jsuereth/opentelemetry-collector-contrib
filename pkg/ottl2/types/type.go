// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package types // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"

// A Type we can evaluate.
type Type interface {
	// A name for this type (usable for error messages)
	Name() string

	// TODO - Kind checking?
	// - AnyValue
	// - Map
	// - List
	// - Custom(Name)

	// TODO - trait description, so we can know that we can safely cast a `Val` to a trait interface.
}

// A type that has members (fields).
type StructType interface {
	Type

	// Returns the type of a field if it exists on this structure type.
	GetField(name string) (Type, bool)

	// Returns field names of a structural type
	FieldNames() []string
}

type primitiveType struct {
	name string
}

func (p primitiveType) Name() string {
	return p.name
}

func NewPrimitiveType(name string) Type {
	return primitiveType{name}
}

type structureType struct {
	name   string
	fields map[string]Type
}

// FieldNames implements StructType.
func (s structureType) FieldNames() []string {
	result := []string{}
	for n, _ := range s.fields {
		result = append(result, n)
	}
	return result
}

func (s structureType) Name() string {
	return s.name
}

func (s structureType) GetField(name string) (Type, bool) {
	for n, t := range s.fields {
		if n == name {
			return t, true
		}
	}
	return nil, false
}

func NewStructureType(name string, fields map[string]Type) StructType {
	return structureType{name, fields}
}
