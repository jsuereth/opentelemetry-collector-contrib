// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpcel

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// Helper to simplify defining new custom type in our type provider.
type Field[S any] struct {
	Name    string
	Type    *types.Type
	IsSet   func(target S) bool
	GetFrom func(target S) (any, error)
}

// Helper implementing our custom type provider.
type StructTypeProvider interface {
	Name() string
	Type() *types.Type
	FindFieldNames() []string
	FindFieldType(fieldName string) (*types.FieldType, bool)
	Construct(map[string]ref.Val) ref.Val
}

// Stores "generic" version of structure fields and constructor.
type realStructType struct {
	name   string
	tpe    *types.Type
	fields map[string]*types.FieldType
	// Construct from pieces
	constructor func(map[string]ref.Val) ref.Val
}

// Construct implements StructTypeProvider.
func (r realStructType) Construct(args map[string]ref.Val) ref.Val {
	return r.constructor(args)
}

// FindFieldNames implements StructTypeProvider.
func (r realStructType) FindFieldNames() []string {
	// TODO: memoize?
	return mapKeys(r.fields)
}

// FindFieldType implements StructTypeProvider.
func (r realStructType) FindFieldType(fieldName string) (*types.FieldType, bool) {
	for n, f := range r.fields {
		if n == fieldName {
			return f, true
		}
	}
	return nil, false
}

// Name implements StructTypeProvider.
func (r realStructType) Name() string {
	return r.name
}

// Type implements StructTypeProvider.
func (r realStructType) Type() *types.Type {
	return r.tpe
}

func NewStructType[S any](
	name string,
	fields []Field[S],
	constructor func(map[string]any) (S, error),
	toVal func(S) ref.Val,
	traits ...int,
) StructTypeProvider {
	fs := map[string]*types.FieldType{}
	for _, f := range fields {
		fs[f.Name] = &types.FieldType{
			Type: f.Type,
			IsSet: func(target any) bool {
				t, ok := target.(S)
				if !ok {
					return false
				}
				return f.IsSet(t)
			},
			GetFrom: func(target any) (any, error) {
				t, ok := target.(S)
				if !ok {
					return nil, fmt.Errorf("cannot find field %s on type %v", f.Name, reflect.TypeFor[S]().Name())
				}
				return f.GetFrom(t)
			},
		}
	}
	t := types.NewObjectType(name, traits...)
	return realStructType{
		name:   name,
		tpe:    t,
		fields: fs,
		constructor: func(m map[string]ref.Val) ref.Val {
			args := map[string]any{}
			for k, v := range m {
				args[k] = v.Value()
			}
			result, err := constructor(args)
			if err != nil {
				return types.NewErrFromString(err.Error())
			}
			return toVal(result)
		},
	}
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	result := make([]K, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i += 1
	}
	return result
}
