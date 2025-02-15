// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/alecthomas/participle/v2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/stdlib"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2/runtime/traits"
) // Parser is responsible for converting an OTTL expression string into an Expr.
type Parser struct {
	env ParserContext
}

// TODO - don't expose this publicly, only expose Statement?
func NewParser(env ParserContext) Parser {
	return Parser{
		env,
	}
}

// Parses a statement into a Condition and an Expression.
func (p *Parser) parseStatement(s parsedStatement) (Interpretable, Interpretable, error) {
	var condition Interpretable
	var err error
	if s.WhereClause != nil {
		condition, err = p.parseBooleanExpression(*s.WhereClause)
	} else {
		condition = BooleanExpr(true)
	}
	if err != nil {
		return nil, nil, err
	}
	action, err := p.parseEditor(s.Editor)
	if err != nil {
		return nil, nil, err
	}
	return condition, action, nil
}

// For testing.
func (p *Parser) ParseValueString(raw string) (Interpretable, error) {
	parsed, err := valueExpressionParser.ParseString("", raw)
	if err != nil {
		return NilExpr(), err
	}
	return p.parseValue(*parsed)
}

// For testing.
func (p *Parser) ParseConditionString(raw string) (Interpretable, error) {
	parsed, err := conditionParser.ParseString("", raw)
	if err != nil {
		return NilExpr(), err
	}
	return p.parseBooleanExpression(*parsed)
}

func (p *Parser) parseValue(v value) (Interpretable, error) {
	switch {
	case v.IsNil != nil:
		return NilExpr(), nil
	case v.Bool != nil:
		if *v.Bool {
			return BooleanExpr(true), nil
		} else {
			return BooleanExpr(false), nil
		}
	case v.String != nil:
		return StringExpr(*v.String), nil
	case v.Bytes != nil:
		return ByteSliceExpr(*v.Bytes), nil
	case v.Enum != nil:
		return p.parseEnum(*v.Enum)
	case v.Literal != nil:
		return p.parseMathExprLiteral(*v.Literal)
	case v.MathExpression != nil:
		return p.parseMathExpression(*v.MathExpression)
	case v.Map != nil:
		return p.parseMap(*v.Map)
	case v.List != nil:
		return p.parseList(*v.List)
	default:
		return NilExpr(), fmt.Errorf("unable to evaluate value: %s", v.checkForCustomError().Error())
	}
}
func (p *Parser) parseEnum(e enumSymbol) (Interpretable, error) {
	value, ok := p.env.ResolveEnum(string(e))
	if !ok {
		return NilExpr(), fmt.Errorf("no such enum: %s", e)
	}
	return ValExpr(value), nil
}

func (p *Parser) parseEditor(e editor) (Interpretable, error) {
	// There are two options for arguments - by name or by position.
	// We currently don't handle by posiiton.
	if f, ok := p.env.ResolveFunction(e.Function); ok {
		pa, na, err := p.parseArguments(e.Arguments)
		if err != nil {
			return NilExpr(), err
		}
		// TODO - we should use type system to erase named args to positional.
		return FuncCallExpr(f, pa, na), nil
	}
	return NilExpr(), fmt.Errorf("could not find function: %s in %v", e.Function, p.env)
}

func (p *Parser) parseArguments(a []argument) ([]Interpretable, map[string]Interpretable, error) {
	positional := []Interpretable{}
	named := map[string]Interpretable{}
	seenNamed := false
	for _, v := range a {
		next, err := p.parseArgument(v)
		if err != nil {
			return nil, nil, err
		}
		if seenNamed && v.Name != "" {
			return nil, nil, errors.New("unnamed argument used after named argument")
		}
		if v.Name != "" {
			seenNamed = true
			named[v.Name] = next
		} else {
			positional = append(positional, next)
		}
	}
	return positional, named, nil
}

func (p *Parser) parseArgument(a argument) (Interpretable, error) {
	// Example for function name: replace_pattern(attributes["message"], Sha256)
	if a.FunctionName != nil {
		// TODO - return function as a value
		if f, ok := p.env.ResolveFunction(*a.FunctionName); ok {
			panic(fmt.Sprintf("function name arguments unsupported, found %v", f))
		} else {
			return NilExpr(), fmt.Errorf("unable to find function: %s", *a.FunctionName)
		}

	}
	return p.parseValue(a.Value)
}

func (p *Parser) parsePath(e path) (Interpretable, error) {
	// So we MAY have context, or we MAY just have a field.
	// We likely want to attach return types to Interpretable at some point.
	var result Interpretable = nil
	var currentType runtime.Type = stdlib.NilType
	if e.Context != "" {
		t, ok := p.env.ResolveName(e.Context)
		if !ok {
			return NilExpr(), fmt.Errorf("invalid name: %s", e.Context)
		}
		currentType = t
		result = LookupExpr(e.Context)
	}
	for _, field := range e.Fields {
		if field.Name != "" {
			if result == nil {
				t, ok := p.env.ResolveName(field.Name)
				if !ok {
					return NilExpr(), fmt.Errorf("invalid name: %s", field.Name)
				}
				currentType = t
				result = LookupExpr(field.Name)
			} else {
				if !reflect.TypeOf(currentType).Implements(reflect.TypeFor[runtime.StructType]()) {
					return NilExpr(), fmt.Errorf("type %s has no fields, cannot find: %s", currentType.Name(), field.Name)
				}
				t, ok := currentType.(runtime.StructType).GetField(field.Name)
				if !ok {
					return NilExpr(), fmt.Errorf("type %s has no field named: %s", currentType.Name(), field.Name)
				}
				currentType = t
				result = AccessExpr(result, field.Name)
			}
		}
		if result == nil {
			return NilExpr(), fmt.Errorf("not a valid path: %v", e)
		}
		for _, key := range field.Keys {
			switch {
			case key.Int != nil:
				result = IndexExpr(result, *key.Int)
			case key.String != nil:
				result = KeyExpr(result, *key.String)
			}
		}
	}
	return result, nil
}

func (p *Parser) parseMathExpression(e mathExpression) (Interpretable, error) {
	result, err := p.parseAddSubterm(*e.Left)
	if err != nil {
		return NilExpr(), err
	}
	if e.Right != nil {
		for _, rhs := range e.Right {
			v, err := p.parseMathVaue(*rhs.Term.Left)
			if err != nil {
				return NilExpr(), err
			}
			result = mathOpExpr(rhs.Operator, result, v)
		}
	}
	return result, nil
}

func (p *Parser) parseAddSubterm(e addSubTerm) (Interpretable, error) {
	result, err := p.parseMathVaue(*e.Left)
	if err != nil {
		return NilExpr(), err
	}
	if e.Right != nil {
		for _, rhs := range e.Right {
			v, err := p.parseMathVaue(*rhs.Value)
			if err != nil {
				return NilExpr(), err
			}
			result = mathOpExpr(rhs.Operator, result, v)
		}
	}
	return result, nil
}

func mathOpExpr(op mathOp, lhs Interpretable, rhs Interpretable) Interpretable {
	switch op {
	case add:
		return AddExpr(lhs, rhs)
	case sub:
		return SubExpr(lhs, rhs)
	case mult:
		return MultExpr(lhs, rhs)
	case div:
		return DivExpr(lhs, rhs)
	}
	panic(fmt.Sprintf("unknown math operation: %s", op.String()))
}

func (p *Parser) parseMathVaue(e mathValue) (Interpretable, error) {
	switch {
	case e.Literal != nil:
		return p.parseMathExprLiteral(*e.Literal)
	case e.SubExpression != nil:
		return p.parseMathExpression(*e.SubExpression)
	}
	return NilExpr(), fmt.Errorf("unknown math expression %v", e)
}

func (p *Parser) parseMathExprLiteral(l mathExprLiteral) (Interpretable, error) {
	switch {
	case l.Float != nil:
		return FloatExpr(*l.Float), nil
	case l.Int != nil:
		return IntExpr(*l.Int), nil
	case l.Editor != nil:
		return p.parseEditor(*l.Editor)
	case l.Path != nil:
		return p.parsePath(*l.Path)
	case l.Converter != nil:
		return p.parseConvertor(*l.Converter)
	}
	return NilExpr(), errors.ErrUnsupported
}

func (p *Parser) parseMap(m mapValue) (Interpretable, error) {
	result := map[string]Interpretable{}
	for _, i := range m.Values {
		v, err := p.parseValue(*i.Value)
		if err != nil {
			return NilExpr(), err
		}
		result[*i.Key] = v
	}
	return MapExpr(result), nil
}

func (p *Parser) parseList(l list) (Interpretable, error) {
	args, err := p.parseValues(l.Values)
	if err != nil {
		return NilExpr(), nil
	}
	return ListExpr(args), nil
}

func (p *Parser) parseValues(vs []value) ([]Interpretable, error) {
	result := make([]Interpretable, len(vs))
	for i, v := range vs {
		r, err := p.parseValue(v)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}

func (p *Parser) parseBooleanExpression(be booleanExpression) (Interpretable, error) {
	// Expressions are or'd here?
	// Previous impelmentation assumed only ors were allowed.
	result, err := p.parseTerm(*be.Left)
	if err != nil {
		return NilExpr(), err
	}
	for _, v := range be.Right {
		rhs, err := p.parseTerm(*v.Term)
		if err != nil {
			return NilExpr(), err
		}
		result = newBoolOperatorExpr(v.Operator, result, rhs)
	}
	return result, nil
}

func (p *Parser) parseTerm(t term) (Interpretable, error) {
	// Expressions are and'd here?
	result, err := p.parseBooleanValue(*t.Left)
	if err != nil {
		return NilExpr(), err
	}
	for _, v := range t.Right {
		rhs, err := p.parseBooleanValue(*v.Value)
		if err != nil {
			return NilExpr(), err
		}
		result = newBoolOperatorExpr(v.Operator, result, rhs)
	}
	return result, nil
}

func (p *Parser) parseBooleanValue(v booleanValue) (Interpretable, error) {
	var r Interpretable
	var err error
	switch {
	case v.Comparison != nil:
		r, err = p.parseComparison(*v.Comparison)
	case v.ConstExpr != nil:
		r, err = p.parseConstExpr(*v.ConstExpr)
	case v.SubExpr != nil:
		r, err = p.parseBooleanExpression(*v.SubExpr)
	}
	if v.Negation != nil {
		r = NotExpr(r)
	}
	return r, err
}

func newBoolOperatorExpr(op string, lhs Interpretable, rhs Interpretable) Interpretable {
	switch op {
	case "and":
		return AndExpr(lhs, rhs)
	case "or":
		return OrExpr(lhs, rhs)
	}
	panic(fmt.Sprintf("unknown boolean operator: %s", op))
}

func (p *Parser) parseComparison(ce comparison) (Interpretable, error) {
	lhs, err := p.parseValue(ce.Left)
	if err != nil {
		return NilExpr(), err
	}
	rhs, err := p.parseValue(ce.Right)
	if err != nil {
		return NilExpr(), err
	}
	// We should benchmark and update these to be more "inlinable", if needed.
	switch ce.Op {
	case eq:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := l.(traits.Comparable).Equals(r)
			return stdlib.NewBoolVal(result)
		}), nil
	case ne:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := !l.(traits.Comparable).Equals(r)
			return stdlib.NewBoolVal(result)
		}), nil
	case lt:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := l.(traits.Comparable).LessThan(r)
			return stdlib.NewBoolVal(result)
		}), nil
	case lte:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := l.(traits.Comparable).LessThan(r) || l.(traits.Comparable).Equals(r)
			return stdlib.NewBoolVal(result)
		}), nil
	case gte:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := !l.(traits.Comparable).LessThan(r)
			return stdlib.NewBoolVal(result)
		}), nil
	case gt:
		return NewBinaryOperation(lhs, rhs, func(l runtime.Val, r runtime.Val) runtime.Val {
			result := !(l.(traits.Comparable).LessThan(r) || l.(traits.Comparable).Equals(r))
			return stdlib.NewBoolVal(result)
		}), nil
	}
	return NilExpr(), fmt.Errorf("unknown comparison: %v", ce)
}

func (p *Parser) parseConstExpr(ce constExpr) (Interpretable, error) {
	switch {
	case ce.Boolean != nil:
		if *ce.Boolean {
			return BooleanExpr(true), nil
		} else {
			return BooleanExpr(false), nil
		}
	case ce.Converter != nil:
		return p.parseConvertor(*ce.Converter)
	default:
		return NilExpr(), fmt.Errorf("unhandled boolean operation %v", ce)
	}
}

func (p *Parser) parseConvertor(c converter) (Interpretable, error) {
	// This is what previous parser did.  This requires that convertors
	// and editors share the same relevant field names.
	result, err := p.parseEditor(editor(c))
	if err != nil {
		return NilExpr(), err
	}
	return result, err
}

var (
	statementParser       = newParser[parsedStatement]()
	conditionParser       = newParser[booleanExpression]()
	valueExpressionParser = newParser[value]()
)

// newParser returns a parser that can be used to read a string into a parsedStatement. An error will be returned if the string
// is not formatted for the DSL.
func newParser[G any]() *participle.Parser[G] {
	lex := buildLexer()
	parser, err := participle.Build[G](
		participle.Lexer(lex),
		participle.Unquote("String"),
		participle.Elide("whitespace"),
		participle.UseLookahead(participle.MaxLookahead), // Allows negative lookahead to work properly in 'value' for 'mathExprLiteral'.
	)
	if err != nil {
		panic("Unable to initialize parser; this is a programming error in OTTL:" + err.Error())
	}
	return parser
}

func parseRawStatement(raw string) (*parsedStatement, error) {
	parsed, err := statementParser.ParseString("", raw)
	if err != nil {
		return nil, fmt.Errorf("statement has invalid syntax: %w", err)
	}
	err = parsed.checkForCustomError()
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func parseRawCondition(raw string) (*booleanExpression, error) {
	parsed, err := conditionParser.ParseString("", raw)
	if err != nil {
		return nil, fmt.Errorf("condition has invalid syntax: %w", err)
	}
	err = parsed.checkForCustomError()
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func parseRawValue(raw string) (*value, error) {
	parsed, err := valueExpressionParser.ParseString("", raw)
	if err != nil {
		return nil, fmt.Errorf("value has invalid syntax: %w", err)
	}
	err = parsed.checkForCustomError()
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
