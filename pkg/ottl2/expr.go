// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/stdlib"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/types/traits"
)

// A value that can be interpreted given an OTTL context.
// You can extend capabilities using other interfaces.
type Interpretable interface {
	// Eval an Activation to produce an output.
	Eval(ctx context.Context, ec EvalContext) types.Val
	// Disallow implementations outside this package.
	// unexportedFactoryFunc()
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
		return ValExpr(stdlib.TrueVal)
	}
	return ValExpr(stdlib.FalseVal)
}

func NilExpr() Interpretable {
	return ValExpr(stdlib.NilVal)
}

func IntExpr(v int64) Interpretable {
	return ValExpr(stdlib.NewIntVal(v))
}

func FloatExpr(v float64) Interpretable {
	return ValExpr(stdlib.NewFloatVal(v))
}

func StringExpr(v string) Interpretable {
	return ValExpr(stdlib.NewStringVal(v))
}

func ByteSliceExpr(v []byte) Interpretable {
	return ValExpr(stdlib.NewByteSliceVal(v))
}

type lookUp struct {
	name string
}

func (l *lookUp) Eval(ctx context.Context, ec EvalContext) types.Val {
	if v, ok := ec.ResolveName(l.name); ok {
		return v
	}
	return stdlib.NewErrorVal(fmt.Errorf("unable to find name [%s] on context %v", l.name, ec))
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
	return stdlib.NewListVal(list)
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
	return stdlib.NewMapVal(result)
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

type funcArg struct {
	constant *types.Val
	expr     *Interpretable
}

type funcCall struct {
	f    types.Function
	args []funcArg
}

// Eval implements Interpretable.
func (f *funcCall) Eval(ctx context.Context, ec EvalContext) types.Val {
	args := make([]types.Val, len(f.args))
	for idx, v := range f.args {
		switch {
		case v.constant != nil:
			args[idx] = *v.constant
		case v.expr != nil:
			args[idx] = (*v.expr).Eval(ctx, ec)
		}
	}
	return f.f.Call(args)
}

func FuncCallExpr(f types.Function, args []Interpretable, namedArgs map[string]Interpretable) Interpretable {
	// erase named/default values to ONLY be positional when interpreting.
	names := f.ArgNames()
	defaults := f.DefaultArgs()
	realArgs := make([]funcArg, len(names))
	for i, name := range names {
		if i < len(args) {
			realArgs[i] = funcArg{
				expr: &args[i],
			}
		} else if v, ok := namedArgs[name]; name != "" && ok {
			realArgs[i] = funcArg{
				expr: &v,
			}
		} else if v, ok := defaults[name]; name != "" && ok {
			realArgs[i] = funcArg{
				constant: &v,
			}
		} else {
			// TODO - should this ne parser or typer error?
			realArgs[i] = funcArg{
				constant: &stdlib.NilVal,
			}
		}
	}
	return &funcCall{f, realArgs}
}

type negExpr struct {
	e Interpretable
}

func (n *negExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	orig, err := n.e.Eval(ctx, ec).ConvertTo(stdlib.BoolType)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	return stdlib.NewBoolVal(!orig.(bool))
}

func NotExpr(e Interpretable) Interpretable {
	return &negExpr{e}
}

type andExpr struct {
	lhs Interpretable
	rhs Interpretable
}

func (a *andExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs, err := a.lhs.Eval(ctx, ec).ConvertTo(stdlib.BoolType)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	if !lhs.(bool) {
		return stdlib.NewBoolVal(false)
	}
	rhs, err := a.rhs.Eval(ctx, ec).ConvertTo(stdlib.BoolType)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	return stdlib.NewBoolVal(rhs.(bool))
}

func AndExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &andExpr{lhs, rhs}
}

type orExpr struct {
	lhs Interpretable
	rhs Interpretable
}

func (a *orExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	lhs, err := a.lhs.Eval(ctx, ec).ConvertTo(stdlib.BoolType)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	if lhs.(bool) {
		return stdlib.NewBoolVal(true)
	}
	rhs, err := a.rhs.Eval(ctx, ec).ConvertTo(stdlib.BoolType)
	if err != nil {
		return stdlib.NewErrorVal(err)
	}
	return stdlib.NewBoolVal(rhs.(bool))
}

func OrExpr(lhs Interpretable, rhs Interpretable) Interpretable {
	return &orExpr{lhs, rhs}
}

type binOpFunc func(types.Val, types.Val) types.Val
type binOpExpr struct {
	lhs Interpretable
	rhs Interpretable
	f   binOpFunc
}

func (b *binOpExpr) Eval(ctx context.Context, ec EvalContext) types.Val {
	// TODO - should we check for errors?
	lhs := b.lhs.Eval(ctx, ec)
	rhs := b.rhs.Eval(ctx, ec)
	return b.f(lhs, rhs)
}

func NewBinaryOperation(lhs Interpretable, rhs Interpretable, op func(types.Val, types.Val) types.Val) Interpretable {
	return &binOpExpr{lhs, rhs, op}
}
