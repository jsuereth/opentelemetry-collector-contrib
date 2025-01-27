// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import "context"

// A value that can be interpreted given an OTTL context.
// You can extend capabilities using other interfaces.
type Interpretable interface {
	// Eval an Activation to produce an output.
	Eval(ctx context.Context, ec EvalContext) Val
}

// Type contianing literal expression
type literalExpr struct {
	value Val
}

func (le *literalExpr) Eval(ctx context.Context, ec EvalContext) Val {
	return le.value
}

func ValExpr(v Val) Interpretable {
	return &literalExpr{v}
}

func BooleanExpr(v bool) Interpretable {
	if v {
		return ValExpr(trueVal)
	}
	return ValExpr(falseVal)
}

func NilExpr() Interpretable {
	return ValExpr(theNilVal)
}

func IntExpr(v int64) Interpretable {
	return ValExpr(newIntVal(v))
}

func FloatExpr(v float64) Interpretable {
	return ValExpr(newFloatVal(v))
}

func StringExpr(v string) Interpretable {
	return ValExpr(newStringVal(v))
}

func ByteSliceExpr(v []byte) Interpretable {
	return ValExpr(newByteSliceVal(v))
}

type listExpr struct {
	items []Interpretable
}

func (l *listExpr) Eval(ctx context.Context, ec EvalContext) Val {
	list := make([]any, len(l.items))
	for i, v := range l.items {
		r := v.Eval(ctx, ec)
		list[i] = r.Value()
	}
	return newListVal(list)
}

func ListExpr(items []Interpretable) Interpretable {
	// TODO - check the type of the list.
	return &listExpr{items}
}

type mapExpr struct {
	items map[string]Interpretable
}

func (m *mapExpr) Eval(ctx context.Context, ec EvalContext) Val {
	result := map[string]any{}
	for k, v := range m.items {
		r := v.Eval(ctx, ec).Value()
		result[k] = r
	}
	return newMapVal(result)
}

func MapExpr(items map[string]Interpretable) Interpretable {
	// TODO - check the type of the map, and choose a better casting strategy
	return &mapExpr{items}
}
