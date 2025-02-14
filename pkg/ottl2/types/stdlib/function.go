// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdlib // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
) // This type does not support named args.
type simpleFunc struct {
	name    string
	numArgs int
	f       func([]types.Val) types.Val
}

// Name implements Function.
func (f *simpleFunc) Name() string {
	return f.name
}

// ArgNames implements Function.
func (f *simpleFunc) ArgNames() []string {
	return make([]string, f.numArgs)
}

// DefaultArgs implements Function.
func (f *simpleFunc) DefaultArgs() map[string]types.Val {
	return map[string]types.Val{}
}

func (f *simpleFunc) Call(args []types.Val) types.Val {
	return f.f(args)
}

func NewSimpleFunc(name string, numArgs int, f func([]types.Val) types.Val) types.Function {
	return &simpleFunc{name, numArgs, f}
}

type advancedFunction struct {
	name        string
	argNames    []string
	defaultArgs map[string]types.Val
	f           func([]types.Val) types.Val
}

// Name implements Function.
func (a advancedFunction) Name() string {
	return a.name
}

// ArgNames implements Function.
func (a advancedFunction) ArgNames() []string {
	return a.argNames
}

// Call implements Function.
func (a advancedFunction) Call(args []types.Val) types.Val {
	return a.f(args)
}

// DefaultArgs implements Function.
func (a advancedFunction) DefaultArgs() map[string]types.Val {
	return a.defaultArgs
}

// Constructs a function that can have named or default parameters.
func NewFunc(
	// Name of the function.
	name string,
	// Names of argument values.  Empty strings denote positional only args.
	argNames []string,
	// Default arguments. Must be used with named parameters.
	defaultArgs map[string]types.Val,
	// A positional-only implementation of the function.
	// All named/default arguments will be turned positional using argNames before calling this.
	f func([]types.Val) types.Val,
) types.Function {
	// TODO - verify named arguments exist in argument list.
	return advancedFunction{name, argNames, defaultArgs, f}
}

// A function where the arguments is a structure.
type reflectArgsFunc[T any] struct {
	// Name is the canonical name to be used by the user when invocating
	// the function generated by this Factory.
	name string
	// Initial argument struct, specific to this function.
	args *T
	// The implementation of the function.
	impl func(*T) types.Val
}

// ArgNames implements types.Function.
func (r *reflectArgsFunc[T]) ArgNames() []string {
	argsVal := reflect.ValueOf(r.args).Elem()
	result := []string{}
	for i := 0; i < argsVal.NumField(); i++ {
		result = append(result, argsVal.Type().Field(i).Name)
	}
	return result
}

// Call implements types.Function.
func (r *reflectArgsFunc[T]) Call(args []types.Val) types.Val {
	var rArgs any = reflect.New(reflect.ValueOf(r.args).Elem().Type()).Interface()
	argsVal := reflect.ValueOf(rArgs).Elem()
	for i := 0; i < argsVal.NumField(); i++ {
		field := argsVal.Field(i)
		fieldType := field.Type()
		isOptional := strings.HasPrefix(fieldType.Name(), "Optional")
		if isOptional {
			manager, ok := field.Interface().(optionalManager)
			if !ok {
				return NewErrorVal(errors.New("optional type is not manageable by the OTTL parser. This is an error in the OTTL"))
			}
			manager.set(args[i].Value())
		} else {
			field.Set(reflect.ValueOf(args[i].Value()))
		}
	}
	return r.impl(rArgs.(*T))
}

// DefaultArgs implements types.Function.
func (r *reflectArgsFunc[T]) DefaultArgs() map[string]types.Val {
	// TODO - precalculate this on creating the struct?
	// TODO - Allow other defaults besides optional.
	defaultArgs := map[string]types.Val{}
	argsVal := reflect.ValueOf(r.args).Elem()
	for i := 0; i < argsVal.NumField(); i++ {
		field := argsVal.Field(i)
		isOptional := strings.HasPrefix(field.Type().Name(), "Optional")
		if isOptional {
			defaultArgs[argsVal.Type().Field(i).Name] = NilVal

		}
	}
	return defaultArgs
}

// Name implements types.Function.
func (r *reflectArgsFunc[T]) Name() string {
	return r.name
}

// Constructs a new function that leverages a structure for all arguments.
//
//   - args MUST be a pointer to a structure.
//   - The name of fields in the structure become the name of allowed arguments to
//     the function.
//   - Any field using `ottl2.Optional` will not be required to be provided, and
//     default to an empty optional value.
func NewReflectFunc[T any](
	name string,
	args *T,
	impl func(*T) types.Val,
) types.Function {
	if reflect.TypeOf(args).Kind() != reflect.Pointer {
		// TODO - non-panic error.
		panic(fmt.Sprintf("factory for %s must return pointer to Arguments", name))
	}
	return &reflectArgsFunc[T]{name, args, impl}
}

// optionalManager provides a way for the parser to handle Optional[T] structs
// without needing to know the concrete type of T, which is inaccessible through
// the reflect package.
// Would likely be resolved by https://github.com/golang/go/issues/54393.
type optionalManager interface {
	// set takes a non-reflection value and returns a reflect.Value of
	// an Optional[T] struct with this value set.
	set(val any) reflect.Value
}

type Optional[T any] struct {
	val      T
	hasValue bool
}

// This is called only by reflection.
func (o Optional[T]) set(val any) reflect.Value {
	return reflect.ValueOf(Optional[T]{
		val:      val.(T),
		hasValue: val == nil,
	})
}

func (o Optional[T]) IsEmpty() bool {
	return !o.hasValue
}

func (o Optional[T]) Get() T {
	return o.val
}

// Allows creating an Optional with a value already populated for use in testing
// OTTL functions.
func NewTestingOptional[T any](val T) Optional[T] {
	return Optional[T]{
		val:      val,
		hasValue: true,
	}
}
