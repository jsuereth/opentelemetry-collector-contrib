// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

// Val interface defines the functions supported by all expression values.
type Val interface {
	// Understands how to convert this type to others.
	// ConvertTo(typeDesc reflect.Type) (any, error)

	// Equal returns true if the `other` value has the same type and content as the implementing struct.
	// Equal(other Val) Val

	// Value returns the raw value of the instance which may not be directly compatible with the expression
	// language types.
	Value() any
}

// Stores literal boolean values
type boolVal struct {
	value bool
}

func (b *boolVal) Value() any {
	return b.value
}

func newBoolVal(v bool) Val {
	return &boolVal{v}
}

var (
	trueVal  = newBoolVal(true)
	falseVal = newBoolVal(false)
)

type nilVal struct{}

func (b *nilVal) Value() any {
	return nil
}

var theNilVal Val = &nilVal{}

type int64Val struct {
	value int64
}

func (i *int64Val) Value() any {
	return i.value
}

func newIntVal(v int64) Val {
	return &int64Val{v}
}

type float64Val struct {
	value float64
}

func (i *float64Val) Value() any {
	return i.value
}

func newFloatVal(v float64) Val {
	return &float64Val{v}
}

type stringVal struct {
	value string
}

func (i *stringVal) Value() any {
	return i.value
}

func newStringVal(v string) Val {
	return &stringVal{v}
}

type byteSliceVal struct {
	value []byte
}

func (i *byteSliceVal) Value() any {
	return i.value
}

func newByteSliceVal(v []byte) Val {
	return &byteSliceVal{v}
}

type listVal[T any] struct {
	value []T
}

func (i *listVal[T]) Value() any {
	return i.value
}

func newListVal[T any](v []T) Val {
	return &listVal[T]{v}
}

type mapVal[K comparable, V any] struct {
	value map[K]V
}

func (m *mapVal[K, V]) Value() any {
	return m.value
}

func newMapVal[K comparable, V any](v map[K]V) Val {
	return &mapVal[K, V]{v}
}
