// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/traits"
)

// A value that can be interpreted given an OTTL context.
// You can extend capabilities using other interfaces.
type Interpretable interface {
	// Eval an Activation to produce an output.
	Eval(ctx context.Context, ec EvalContext) types.Val
}

// Type contianing literal expression
type literalExpr struct {
	value types.Val
}

func (le *literalExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	return le.value
}

func ValExpr(v types.Val) Interpretable {
	return &literalExpr{v}
}

func BooleanExpr(v bool) Interpretable {
	if v {
		return ValExpr(types.TrueVal)
	}
	return ValExpr(types.FalseVal)
}

func NilExpr() Interpretable {
	return ValExpr(types.NilVal)
}

func IntExpr(v int64) Interpretable {
	return ValExpr(types.NewIntVal(v))
}

func FloatExpr(v float64) Interpretable {
	return ValExpr(types.NewFloatVal(v))
}

func StringExpr(v string) Interpretable {
	return ValExpr(types.NewStringVal(v))
}

func ByteSliceExpr(v []byte) Interpretable {
	return ValExpr(types.NewByteSliceVal(v))
}

type lookUp struct {
	name string
}

func (l *lookUp) Eval(ctx context.Context, ec EvalContext) types.Val {
	if v, ok := ec.ResolveName(l.name); ok {
		return v
	}
	panic(fmt.Sprintf("Unable to find name [%s] on context %v", l.name, ec))
}

func LookupExpr(n string) Interpretable {
	return &lookUp{n}
}

// Implements {target}.{field} expressions
type accessExpr struct {
	target Interpretable
	field  string
}

// Eval implements Interpretable.
func (a *accessExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	target := a.target.Eval(ctx, ec)
	return target.(traits.StructureAccessible).GetField(a.field)
}

func AccessExpr(target Interpretable, field string) Interpretable {
	return &accessExpr{target, field}
}

type indexExpr struct {
	target Interpretable
	idx    int64
}

// Eval implements Interpretable.
func (i *indexExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	target := i.target.Eval(ctx, ec)
	return target.(traits.Indexable).GetIndex(i.idx)
}

func IndexExpr(target Interpretable, idx int64) Interpretable {
	return &indexExpr{target, idx}
}

type keyExpr struct {
	target Interpretable
	key    string
}

// Eval implements Interpretable.
func (k *keyExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	target := k.target.Eval(ctx, ec)
	return target.(traits.KeyAccessable).GetKey(k.key)
}

func KeyExpr(target Interpretable, key string) Interpretable {
	return &keyExpr{target, key}
}

type listExpr struct {
	items []Interpretable
}

func (l *listExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	list := make([]types.Val, len(l.items))
	for i, v := range l.items {
		r := v.Eval(ctx, ec)
		list[i] = r
	}
	return types.NewListVal(list)
}

func ListExpr(items []Interpretable) Interpretable {
	// TODO - check the type of the list.
	return &listExpr{items}
}

type mapExpr struct {
	items map[string]Interpretable
}

func (m *mapExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	result := map[string]types.Val{}
	for k, v := range m.items {
		r := v.Eval(ctx, ec)
		result[k] = r
	}
	return types.NewMapVal(result)
}

func MapExpr(items map[string]Interpretable) Interpretable {
	// TODO - check the type of the map, and choose a better casting strategy
	return &mapExpr{items}
}

type addOp struct {
	lhs Interpretable
	rhs Interpretable
}

// Eval implements Interpretable.
func (a *addOp) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs := a.lhs.Eval(ctx, ec)
	rhs := a.rhs.Eval(ctx, ec)
	return lhs.(traits.Adder).Add(rhs)
}

func AddExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &addOp{lhs, rhs}
}

type subOp struct {
	lhs Interpretable
	rhs Interpretable
}

func (m *subOp) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs := m.lhs.Eval(ctx, ec)
	rhs := m.rhs.Eval(ctx, ec)
	return lhs.(traits.Subtractor).Subtract(rhs)
}

func SubExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &subOp{lhs, rhs}
}

type multOp struct {
	lhs Interpretable
	rhs Interpretable
}

func (m *multOp) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs := m.lhs.Eval(ctx, ec)
	rhs := m.rhs.Eval(ctx, ec)
	return lhs.(traits.Multiplier).Multiply(rhs)
}

func MultExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &multOp{lhs, rhs}
}

type divOp struct {
	lhs Interpretable
	rhs Interpretable
}

func (m *divOp) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs := m.lhs.Eval(ctx, ec)
	rhs := m.rhs.Eval(ctx, ec)
	return lhs.(traits.Divider).Divide(rhs)
}

func DivExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &divOp{lhs, rhs}
}

// TODO - this should store the raw Function interface, not lookup via name.
type funcCall struct {
	f    types.Function
	args []Interpretable
}

// Eval implements Interpretable.
func (f *funcCall) Eval(ctx context.Context, ec EvalContext) types.Val {
	args := make([]types.Val, len(f.args))
	for idx, v := range f.args {
		args[idx] = v.Eval(ctx, ec)
	}
	result, err := f.f.Call(args)
	if err == nil {
		return result
	}
	// TODO - create error type to return...
	panic(fmt.Sprintf("Invalid function call: %v", f))
}

func FuncCallExpr(f types.Function, args []Interpretable) Interpretable {
	return &funcCall{f, args}
}
