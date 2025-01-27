// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
) // Parser is responsible for converting an OTTL expression string into an Expr.
type Parser struct {
	env ParserContext
}

func NewParser(env ParserContext) Parser {
	return Parser{
		env,
	}
}

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
	// There are two options for arguments - by name or by position.
	// We currently don't handle by posiiton.
	if f, ok := p.env.ResolveFunction(e.Function); ok {
		args, err := p.parseArguments(e.Arguments)
		if err != nil {
			return NilExpr(), err
		}
		return FuncCallExpr(f, args), nil
	}
	return NilExpr(), fmt.Errorf("could not find function: %s in %v", e.Function, p.env)
}

func (p *Parser) parseArguments(a []argument) ([]Interpretable, error) {
	result := make([]Interpretable, len(a))
	for i, v := range a {
		next, err := p.parseArgument(v)
		if err != nil {
			return []Interpretable{}, err
		}
		result[i] = next
	}
	return result, nil
}

func (p *Parser) parseArgument(a argument) (Interpretable, error) {
	if a.Name != "" || a.FunctionName != nil {
		panic(fmt.Sprintf("named arguments unsupported, found %v", a))
	}
	return p.parseValue(a.Value)
}

func (p *Parser) parsePath(e path) (Interpretable, error) {
	// So we MAY have context, or we MAY just have a field...
	// TODO - Verify the path exists via types, if we can.
	var result Interpretable = nil
	if e.Context != "" {
		if !p.env.HasName(e.Context) {
			return NilExpr(), fmt.Errorf("invalid name: %s", e.Context)
		}
		result = LookupExpr(e.Context)
	}
	for _, field := range e.Fields {
		if field.Name != "" {
			if result == nil {
				if !p.env.HasName(field.Name) {
					return NilExpr(), fmt.Errorf("invalid name: %s", field.Name)
				}
				result = LookupExpr(field.Name)
			} else {
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

// TODO _ maybe convert these into function calls against "Add", "Mult" etc. functions.
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
	panic("Unknown math operation")
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
