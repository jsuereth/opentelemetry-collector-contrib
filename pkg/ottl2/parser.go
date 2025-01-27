// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
) // Parser is responsible for converting an OTTL expression string into an Expr.
type Parser struct{}

// // Parses a value expression
func (p *Parser) ParseValueString(raw string) (Interpretable, error) {
	parsed, err := valueExpressionParser.ParseString("", raw)
	if err != nil {
		return NilExpr(), err
	}
	return p.parseValue(*parsed)
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
	// TODO - use provided enumerations in context.
	panic("Enums not supported")
}

func (p *Parser) parseEditor(e editor) (Interpretable, error) {
	// TODO - Look up editor and verify it.
	panic("Editor not supported")
}
func (p *Parser) parsePath(e path) (Interpretable, error) {
	// TODO - Look up path and verify it.
	// return pathExpr(e), nil
	panic("not implemented")
}

func (p *Parser) parseMathExpression(e mathExpression) (Interpretable, error) {
	// result, err := p.parseAddSubterm(*e.Left)
	// if err != nil {
	// 	return NilExpr(), nil
	// }
	// for _, rhs := range e.Right {
	// 	v, err := p.parseAddSubterm(*rhs.Term)
	// 	if err != nil {
	// 		return NilExpr(), err
	// 	}
	// 	result = mathExpr(result, rhs.Operator, v)
	// }
	// return mapExpr[T](kvs), nil
	panic("not implemented")
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

// // converts a constExpr AST to an evaluatable Expr.
// func (p *Parser) parseConstExpr(ce constExpr) (Interpretable, error) {
// 	switch {
// 	case ce.Boolean != nil:
// 		if *ce.Boolean {
// 			return trueExpr(), nil
// 		} else {
// 			return falseExpr(), nil
// 		}
// 	case ce.Converter != nil:
// 		return p.parseConvertor(*ce.Converter)
// 	default:
// 		return NilExpr(), fmt.Errorf("unhandled boolean operation %v", ce)
// 	}
// }

func (p *Parser) parseConvertor(c converter) (Interpretable, error) {
	// TODO - these operate wierdly in OTTL, using some Getter/Setter abstraction
	// We need to re-implement.
	return NilExpr(), errors.ErrUnsupported
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
