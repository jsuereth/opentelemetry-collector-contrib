// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottl2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl2"

import (
	"regexp"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottltest"
)

// This is not in ottltest because it depends on a type that's a member of OTTL.
func booleanp(b boolean) *boolean {
	return &b
}

func Test_parse(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		expected  *parsedStatement
	}{
		{
			name:      "editor with string",
			statement: `set("foo")`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								String: ottltest.Strp("foo"),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with float",
			statement: `met(1.2)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "met",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Float: ottltest.Floatp(1.2),
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with int",
			statement: `fff(12)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "fff",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Int: ottltest.Intp(12),
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with map",
			statement: `fff({"stringAttr": "value", "intAttr": 3, "floatAttr": 2.5, "boolAttr": true})`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "fff",
					Arguments: []argument{
						{
							Value: value{
								Map: &mapValue{
									Values: []mapItem{
										{
											Key:   ottltest.Strp("stringAttr"),
											Value: &value{String: ottltest.Strp("value")},
										},
										{
											Key: ottltest.Strp("intAttr"),
											Value: &value{
												Literal: &mathExprLiteral{
													Int: ottltest.Intp(3),
												},
											},
										},
										{
											Key: ottltest.Strp("floatAttr"),
											Value: &value{
												Literal: &mathExprLiteral{
													Float: ottltest.Floatp(2.5),
												},
											},
										},
										{
											Key:   ottltest.Strp("boolAttr"),
											Value: &value{Bool: (*boolean)(ottltest.Boolp(true))},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with empty map",
			statement: `fff({})`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "fff",
					Arguments: []argument{
						{
							Value: value{
								Map: &mapValue{
									Values: nil,
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with converter with a map",
			statement: `fff(GetSomething({"foo":"bar"}))`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "fff",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Converter: &converter{
										Function: "GetSomething",
										Arguments: []argument{
											{
												Value: value{
													Map: &mapValue{
														Values: []mapItem{
															{
																Key:   ottltest.Strp("foo"),
																Value: &value{String: ottltest.Strp("bar")},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with nested map",
			statement: `fff({"mapAttr": {"foo": "bar", "get": bear.honey, "arrayAttr":["foo", "bar"]}})`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "fff",
					Arguments: []argument{
						{
							Value: value{
								Map: &mapValue{
									Values: []mapItem{
										{
											Key: ottltest.Strp("mapAttr"),
											Value: &value{
												Map: &mapValue{
													Values: []mapItem{
														{
															Key:   ottltest.Strp("foo"),
															Value: &value{String: ottltest.Strp("bar")},
														},
														{
															Key: ottltest.Strp("get"),
															Value: &value{
																Literal: &mathExprLiteral{
																	Path: &path{
																		Pos: lexer.Position{
																			Offset: 38,
																			Line:   1,
																			Column: 39,
																		},
																		Context: "bear",
																		Fields: []field{
																			{
																				Name: "honey",
																			},
																		},
																	},
																},
															},
														},
														{
															Key: ottltest.Strp("arrayAttr"),
															Value: &value{
																List: &list{
																	Values: []value{
																		{
																			String: ottltest.Strp("foo"),
																		},
																		{
																			String: ottltest.Strp("bar"),
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "complex editor",
			statement: `set("foo", GetSomething(bear.honey))`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								String: ottltest.Strp("foo"),
							},
						},
						{
							Value: value{
								Literal: &mathExprLiteral{
									Converter: &converter{
										Function: "GetSomething",
										Arguments: []argument{
											{
												Value: value{
													Literal: &mathExprLiteral{
														Path: &path{
															Pos: lexer.Position{
																Offset: 24,
																Line:   1,
																Column: 25,
															},
															Context: "bear",
															Fields: []field{
																{
																	Name: "honey",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "complex path",
			statement: `set(foo.attributes["bar"].cat, "dog")`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Context: "foo",
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bar"),
													},
												},
											},
											{
												Name: "cat",
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("dog"),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "single field segment",
			statement: `set(attributes["bar"], "dog")`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Context: "",
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bar"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("dog"),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter parameters (All Uppercase)",
			statement: `replace_pattern(attributes["message"], "device=*", attributes["device_name"], SHA256)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "replace_pattern",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 16,
											Line:   1,
											Column: 17,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("message"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("device=*"),
							},
						},
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 51,
											Line:   1,
											Column: 52,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("device_name"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								Enum: (*enumSymbol)(ottltest.Strp("SHA256")),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter parameters",
			statement: `replace_pattern(attributes["message"], Sha256)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "replace_pattern",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 16,
											Line:   1,
											Column: 17,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("message"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							FunctionName: ottltest.Strp("Sha256"),
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter parameters (One Uppercase symbol)",
			statement: `replace_pattern(attributes["message"], S)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "replace_pattern",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 16,
											Line:   1,
											Column: 17,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("message"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								Enum: (*enumSymbol)(ottltest.Strp("S")),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "complex path",
			statement: `set(foo.bar["x"]["y"].z, Test()[0]["pass"])`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Context: "foo",
										Fields: []field{
											{
												Name: "bar",
												Keys: []key{
													{
														String: ottltest.Strp("x"),
													},
													{
														String: ottltest.Strp("y"),
													},
												},
											},
											{
												Name: "z",
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								Literal: &mathExprLiteral{
									Converter: &converter{
										Function: "Test",
										Keys: []key{
											{
												Int: ottltest.Intp(0),
											},
											{
												String: ottltest.Strp("pass"),
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "where == clause",
			statement: `set(foo.attributes["bar"].cat, "dog") where name == "fido"`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Context: "foo",
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bar"),
													},
												},
											},
											{
												Name: "cat",
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("dog"),
							},
						},
					},
				},
				WhereClause: &booleanExpression{
					Left: &term{
						Left: &booleanValue{
							Comparison: &comparison{
								Left: value{
									Literal: &mathExprLiteral{
										Path: &path{
											Pos: lexer.Position{
												Offset: 44,
												Line:   1,
												Column: 45,
											},
											Fields: []field{
												{
													Name: "name",
												},
											},
										},
									},
								},
								Op: eq,
								Right: value{
									String: ottltest.Strp("fido"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "where != clause",
			statement: `set(foo.attributes["bar"].cat, "dog") where name != "fido"`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Context: "foo",
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bar"),
													},
												},
											},
											{
												Name: "cat",
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("dog"),
							},
						},
					},
				},
				WhereClause: &booleanExpression{
					Left: &term{
						Left: &booleanValue{
							Comparison: &comparison{
								Left: value{
									Literal: &mathExprLiteral{
										Path: &path{
											Pos: lexer.Position{
												Offset: 44,
												Line:   1,
												Column: 45,
											},
											Fields: []field{
												{
													Name: "name",
												},
											},
										},
									},
								},
								Op: ne,
								Right: value{
									String: ottltest.Strp("fido"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "ignore extra spaces",
			statement: `set  ( foo.attributes[ "bar"].cat,   "dog")   where name=="fido"`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 7,
											Line:   1,
											Column: 8,
										},
										Context: "foo",
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bar"),
													},
												},
											},
											{
												Name: "cat",
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								String: ottltest.Strp("dog"),
							},
						},
					},
				},
				WhereClause: &booleanExpression{
					Left: &term{
						Left: &booleanValue{
							Comparison: &comparison{
								Left: value{
									Literal: &mathExprLiteral{
										Path: &path{
											Pos: lexer.Position{
												Offset: 52,
												Line:   1,
												Column: 53,
											},
											Fields: []field{
												{
													Name: "name",
												},
											},
										},
									},
								},
								Op: eq,
								Right: value{
									String: ottltest.Strp("fido"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "handle quotes",
			statement: `set("fo\"o")`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								String: ottltest.Strp("fo\"o"),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with boolean false",
			statement: `convert_gauge_to_sum("cumulative", false)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "convert_gauge_to_sum",
					Arguments: []argument{
						{
							Value: value{
								String: ottltest.Strp("cumulative"),
							},
						},
						{
							Value: value{
								Bool: (*boolean)(ottltest.Boolp(false)),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with boolean true",
			statement: `convert_gauge_to_sum("cumulative", true)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "convert_gauge_to_sum",
					Arguments: []argument{
						{
							Value: value{
								String: ottltest.Strp("cumulative"),
							},
						},
						{
							Value: value{
								Bool: (*boolean)(ottltest.Boolp(true)),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with bytes",
			statement: `set(attributes["bytes"], 0x0102030405060708)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("bytes"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								Bytes: (*byteSlice)(&[]byte{1, 2, 3, 4, 5, 6, 7, 8}),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with nil",
			statement: `set(attributes["test"], nil)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								IsNil: (*isNil)(ottltest.Boolp(true)),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "editor with Enum",
			statement: `set(attributes["test"], TEST_ENUM)`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								Enum: (*enumSymbol)(ottltest.Strp("TEST_ENUM")),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter with empty list",
			statement: `set(attributes["test"], [])`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								List: &list{
									Values: nil,
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter with single-value list",
			statement: `set(attributes["test"], ["value0"])`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								List: &list{
									Values: []value{
										{
											String: ottltest.Strp("value0"),
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter with multi-value list",
			statement: `set(attributes["test"], ["value1", "value2"])`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								List: &list{
									Values: []value{
										{
											String: ottltest.Strp("value1"),
										},
										{
											String: ottltest.Strp("value2"),
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter with nested heterogeneous types",
			statement: `set(attributes["test"], [Concat(["a", "b"], "+"), ["1", 2, 3.0], nil, attributes["test"]])`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								List: &list{
									Values: []value{
										{
											Literal: &mathExprLiteral{
												Converter: &converter{
													Function: "Concat",
													Arguments: []argument{
														{
															Value: value{
																List: &list{
																	Values: []value{
																		{
																			String: ottltest.Strp("a"),
																		},
																		{
																			String: ottltest.Strp("b"),
																		},
																	},
																},
															},
														},
														{
															Value: value{
																String: ottltest.Strp("+"),
															},
														},
													},
												},
											},
										},
										{
											List: &list{
												Values: []value{
													{
														String: ottltest.Strp("1"),
													},
													{
														Literal: &mathExprLiteral{
															Int: ottltest.Intp(2),
														},
													},
													{
														Literal: &mathExprLiteral{
															Float: ottltest.Floatp(3.0),
														},
													},
												},
											},
										},
										{
											IsNil: (*isNil)(ottltest.Boolp(true)),
										},
										{
											Literal: &mathExprLiteral{
												Path: &path{
													Pos: lexer.Position{
														Offset: 70,
														Line:   1,
														Column: 71,
													},
													Fields: []field{
														{
															Name: "attributes",
															Keys: []key{
																{
																	String: ottltest.Strp("test"),
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
		{
			name:      "Converter math mathExpression",
			statement: `set(attributes["test"], 1000 - 600) where 1 + 1 * 2 == three / One()`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Value: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 4,
											Line:   1,
											Column: 5,
										},
										Fields: []field{
											{
												Name: "attributes",
												Keys: []key{
													{
														String: ottltest.Strp("test"),
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Value: value{
								MathExpression: &mathExpression{
									Left: &addSubTerm{
										Left: &mathValue{
											Literal: &mathExprLiteral{
												Int: ottltest.Intp(1000),
											},
										},
									},
									Right: []*opAddSubTerm{
										{
											Operator: sub,
											Term: &addSubTerm{
												Left: &mathValue{
													Literal: &mathExprLiteral{
														Int: ottltest.Intp(600),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				WhereClause: &booleanExpression{
					Left: &term{
						Left: &booleanValue{
							Comparison: &comparison{
								Left: value{
									MathExpression: &mathExpression{
										Left: &addSubTerm{
											Left: &mathValue{
												Literal: &mathExprLiteral{
													Int: ottltest.Intp(1),
												},
											},
										},
										Right: []*opAddSubTerm{
											{
												Operator: add,
												Term: &addSubTerm{
													Left: &mathValue{
														Literal: &mathExprLiteral{
															Int: ottltest.Intp(1),
														},
													},
													Right: []*opMultDivValue{
														{
															Operator: mult,
															Value: &mathValue{
																Literal: &mathExprLiteral{
																	Int: ottltest.Intp(2),
																},
															},
														},
													},
												},
											},
										},
									},
								},
								Op: eq,
								Right: value{
									MathExpression: &mathExpression{
										Left: &addSubTerm{
											Left: &mathValue{
												Literal: &mathExprLiteral{
													Path: &path{
														Pos: lexer.Position{
															Offset: 55,
															Line:   1,
															Column: 56,
														},
														Fields: []field{
															{
																Name: "three",
															},
														},
													},
												},
											},
											Right: []*opMultDivValue{
												{
													Operator: div,
													Value: &mathValue{
														Literal: &mathExprLiteral{
															Converter: &converter{
																Function: "One",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "editor with named arg",
			statement: `set(name="foo")`,
			expected: &parsedStatement{
				Editor: editor{
					Function: "set",
					Arguments: []argument{
						{
							Name: "name",
							Value: value{
								String: ottltest.Strp("foo"),
							},
						},
					},
				},
				WhereClause: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.statement, func(t *testing.T) {
			parsed, err := parseRawStatement(tt.statement)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, parsed)
		})
	}
}

func Test_parseCondition_full(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		expected  *booleanExpression
	}{
		{
			name:      "where == clause",
			condition: `name == "fido"`,
			expected: &booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Comparison: &comparison{
							Left: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 0,
											Line:   1,
											Column: 1,
										},
										Fields: []field{
											{
												Name: "name",
											},
										},
									},
								},
							},
							Op: eq,
							Right: value{
								String: ottltest.Strp("fido"),
							},
						},
					},
				},
			},
		},
		{
			name:      "where != clause",
			condition: `name != "fido"`,
			expected: &booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Comparison: &comparison{
							Left: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 0,
											Line:   1,
											Column: 1,
										},
										Fields: []field{
											{
												Name: "name",
											},
										},
									},
								},
							},
							Op: ne,
							Right: value{
								String: ottltest.Strp("fido"),
							},
						},
					},
				},
			},
		},
		{
			name:      "Converter math mathExpression",
			condition: `1 + 1 * 2 == three / One()`,
			expected: &booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Comparison: &comparison{
							Left: value{
								MathExpression: &mathExpression{
									Left: &addSubTerm{
										Left: &mathValue{
											Literal: &mathExprLiteral{
												Int: ottltest.Intp(1),
											},
										},
									},
									Right: []*opAddSubTerm{
										{
											Operator: add,
											Term: &addSubTerm{
												Left: &mathValue{
													Literal: &mathExprLiteral{
														Int: ottltest.Intp(1),
													},
												},
												Right: []*opMultDivValue{
													{
														Operator: mult,
														Value: &mathValue{
															Literal: &mathExprLiteral{
																Int: ottltest.Intp(2),
															},
														},
													},
												},
											},
										},
									},
								},
							},
							Op: eq,
							Right: value{
								MathExpression: &mathExpression{
									Left: &addSubTerm{
										Left: &mathValue{
											Literal: &mathExprLiteral{
												Path: &path{
													Pos: lexer.Position{
														Offset: 13,
														Line:   1,
														Column: 14,
													},
													Fields: []field{
														{
															Name: "three",
														},
													},
												},
											},
										},
										Right: []*opMultDivValue{
											{
												Operator: div,
												Value: &mathValue{
													Literal: &mathExprLiteral{
														Converter: &converter{
															Function: "One",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.condition, func(t *testing.T) {
			parsed, err := parseRawCondition(tt.condition)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, parsed)
		})
	}
}

// Helper for test cases where the WHERE clause is all that matters.
// Parse string should start with `set(name, "test") where`...
func setNameTest(b *booleanExpression) *parsedStatement {
	return &parsedStatement{
		Editor: editor{
			Function: "set",
			Arguments: []argument{
				{
					Value: value{
						Literal: &mathExprLiteral{
							Path: &path{
								Pos: lexer.Position{
									Offset: 4,
									Line:   1,
									Column: 5,
								},
								Fields: []field{
									{
										Name: "name",
									},
								},
							},
						},
					},
				},
				{
					Value: value{
						String: ottltest.Strp("test"),
					},
				},
			},
		},
		WhereClause: b,
	}
}

func Test_parseWhere(t *testing.T) {
	tests := []struct {
		statement string
		expected  *parsedStatement
	}{
		{
			statement: `true`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(true),
						},
					},
				},
			}),
		},
		{
			statement: `true and false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(true),
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `true and true and false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(true),
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(true),
								},
							},
						},
						{
							Operator: "and",
							Value: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `true or false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(true),
						},
					},
				},
				Right: []*opOrTerm{
					{
						Operator: "or",
						Term: &term{
							Left: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `false and true or false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(false),
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(true),
								},
							},
						},
					},
				},
				Right: []*opOrTerm{
					{
						Operator: "or",
						Term: &term{
							Left: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `(false and true) or false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						SubExpr: &booleanExpression{
							Left: &term{
								Left: &booleanValue{
									ConstExpr: &constExpr{
										Boolean: booleanp(false),
									},
								},
								Right: []*opAndBooleanValue{
									{
										Operator: "and",
										Value: &booleanValue{
											ConstExpr: &constExpr{
												Boolean: booleanp(true),
											},
										},
									},
								},
							},
						},
					},
				},
				Right: []*opOrTerm{
					{
						Operator: "or",
						Term: &term{
							Left: &booleanValue{
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `false and (true or false)`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(false),
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								SubExpr: &booleanExpression{
									Left: &term{
										Left: &booleanValue{
											ConstExpr: &constExpr{
												Boolean: booleanp(true),
											},
										},
									},
									Right: []*opOrTerm{
										{
											Operator: "or",
											Term: &term{
												Left: &booleanValue{
													ConstExpr: &constExpr{
														Boolean: booleanp(false),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `name != "foo" and name != "bar"`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Comparison: &comparison{
							Left: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 24,
											Line:   1,
											Column: 25,
										},
										Fields: []field{
											{
												Name: "name",
											},
										},
									},
								},
							},
							Op: ne,
							Right: value{
								String: ottltest.Strp("foo"),
							},
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								Comparison: &comparison{
									Left: value{
										Literal: &mathExprLiteral{
											Path: &path{
												Pos: lexer.Position{
													Offset: 42,
													Line:   1,
													Column: 43,
												},
												Fields: []field{
													{
														Name: "name",
													},
												},
											},
										},
									},
									Op: ne,
									Right: value{
										String: ottltest.Strp("bar"),
									},
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `name == "foo" or name == "bar"`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Comparison: &comparison{
							Left: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 24,
											Line:   1,
											Column: 25,
										},
										Fields: []field{
											{
												Name: "name",
											},
										},
									},
								},
							},
							Op: eq,
							Right: value{
								String: ottltest.Strp("foo"),
							},
						},
					},
				},
				Right: []*opOrTerm{
					{
						Operator: "or",
						Term: &term{
							Left: &booleanValue{
								Comparison: &comparison{
									Left: value{
										Literal: &mathExprLiteral{
											Path: &path{
												Pos: lexer.Position{
													Offset: 41,
													Line:   1,
													Column: 42,
												},
												Fields: []field{
													{
														Name: "name",
													},
												},
											},
										},
									},
									Op: eq,
									Right: value{
										String: ottltest.Strp("bar"),
									},
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `true and not false`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Boolean: booleanp(true),
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								Negation: ottltest.Strp("not"),
								ConstExpr: &constExpr{
									Boolean: booleanp(false),
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `not name == "bar"`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Negation: ottltest.Strp("not"),
						Comparison: &comparison{
							Left: value{
								Literal: &mathExprLiteral{
									Path: &path{
										Pos: lexer.Position{
											Offset: 28,
											Line:   1,
											Column: 29,
										},
										Fields: []field{
											{
												Name: "name",
											},
										},
									},
								},
							},
							Op: eq,
							Right: value{
								String: ottltest.Strp("bar"),
							},
						},
					},
				},
			}),
		},
		{
			statement: `not (true or false)`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						Negation: ottltest.Strp("not"),
						SubExpr: &booleanExpression{
							Left: &term{
								Left: &booleanValue{
									ConstExpr: &constExpr{
										Boolean: booleanp(true),
									},
								},
							},
							Right: []*opOrTerm{
								{
									Operator: "or",
									Term: &term{
										Left: &booleanValue{
											ConstExpr: &constExpr{
												Boolean: booleanp(false),
											},
										},
									},
								},
							},
						},
					},
				},
			}),
		},
		{
			statement: `True()`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Converter: &converter{
								Function: "True",
							},
						},
					},
				},
			}),
		},
		{
			statement: `True() and False()`,
			expected: setNameTest(&booleanExpression{
				Left: &term{
					Left: &booleanValue{
						ConstExpr: &constExpr{
							Converter: &converter{
								Function: "True",
							},
						},
					},
					Right: []*opAndBooleanValue{
						{
							Operator: "and",
							Value: &booleanValue{
								ConstExpr: &constExpr{
									Converter: &converter{
										Function: "False",
									},
								},
							},
						},
					},
				},
			}),
		},
	}

	// create a test name that doesn't confuse vscode so we can rerun tests with one click
	pat := regexp.MustCompile("[^a-zA-Z0-9]+")
	for _, tt := range tests {
		name := pat.ReplaceAllString(tt.statement, "_")
		t.Run(name, func(t *testing.T) {
			statement := `set(name, "test") where ` + tt.statement
			parsed, err := parseRawStatement(statement)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, parsed)
		})
	}
}

// This test doesn't validate parser results, simply checks whether the parse succeeds or not.
// It's a fast way to check a large range of possible syntaxes.
func Test_parseStatement(t *testing.T) {
	converterNameErrorPrefix := "converter names must start with an uppercase letter"
	editorWithIndexErrorPrefix := "only paths and converters may be indexed"

	tests := []struct {
		statement         string
		wantErr           bool
		wantErrContaining string
	}{
		{statement: `set(`, wantErr: true},
		{statement: `set("foo)`, wantErr: true},
		{statement: `set(name.)`, wantErr: true},
		{statement: `("foo")`, wantErr: true},
		{statement: `set("foo") where name =||= "fido"`, wantErr: true},
		{statement: `set(span_id, SpanIDWrapper{not a hex string})`, wantErr: true},
		{statement: `set(span_id, SpanIDWrapper{01})`, wantErr: true},
		{statement: `set(span_id, SpanIDWrapper{010203040506070809})`, wantErr: true},
		{statement: `set(trace_id, TraceIDWrapper{not a hex string})`, wantErr: true},
		{statement: `set(trace_id, TraceIDWrapper{0102030405060708090a0b0c0d0e0f})`, wantErr: true},
		{statement: `set(trace_id, TraceIDWrapper{0102030405060708090a0b0c0d0e0f1011})`, wantErr: true},
		{statement: `set("foo") where name = "fido"`, wantErr: true},
		{statement: `set("foo") where name or "fido"`, wantErr: true},
		{statement: `set("foo") where name and "fido"`, wantErr: true},
		{statement: `set("foo") where name and`, wantErr: true},
		{statement: `set("foo") where name or`, wantErr: true},
		{statement: `set("foo") where (`, wantErr: true},
		{statement: `set("foo") where )`, wantErr: true},
		{statement: `set("foo") where (name == "fido"))`, wantErr: true},
		{statement: `set("foo") where ((name == "fido")`, wantErr: true},
		{statement: `Set()`, wantErr: true},
		{statement: `set(int())`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set(1 + int())`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set(int() + 1)`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set(1 * int())`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set(1 * 1 + (2 * int()))`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set() where int() == 1`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set() where 1 == int()`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set() where true and 1 == int() `, wantErrContaining: converterNameErrorPrefix},
		{statement: `set() where false or 1 == int() `, wantErrContaining: converterNameErrorPrefix},
		{statement: `set(foo.attributes["bar"].cat)["key"]`, wantErrContaining: editorWithIndexErrorPrefix},
		{statement: `set(foo.attributes["bar"].cat, "dog")`},
		{statement: `set(set = foo.attributes["animal"], val = "dog") where animal == "cat"`},
		{statement: `test() where service == "pinger" or foo.attributes["endpoint"] == "/x/alive"`},
		{statement: `test() where service == "pinger" or foo.attributes["verb"] == "GET" and foo.attributes["endpoint"] == "/x/alive"`},
		{statement: `test() where animal > "cat"`},
		{statement: `test() where animal >= "cat"`},
		{statement: `test() where animal <= "cat"`},
		{statement: `test() where animal < "cat"`},
		{statement: `test() where animal =< "dog"`, wantErr: true},
		{statement: `test() where animal => "dog"`, wantErr: true},
		{statement: `test() where animal <> "dog"`, wantErr: true},
		{statement: `test() where animal = "dog"`, wantErr: true},
		{statement: `test() where animal`, wantErr: true},
		{statement: `test() where animal ==`, wantErr: true},
		{statement: `test() where ==`, wantErr: true},
		{statement: `test() where == animal`, wantErr: true},
		{statement: `test() where attributes["path"] == "/healthcheck"`},
		{statement: `test() where one() == 1`, wantErr: true},
		{statement: `test(fail())`, wantErrContaining: converterNameErrorPrefix},
		{statement: `Test()`, wantErr: true},
		{statement: `set() where test(foo)["key"] == "bar"`, wantErrContaining: converterNameErrorPrefix},
		{statement: `set() where test(foo)["key"] == "bar"`, wantErrContaining: editorWithIndexErrorPrefix},
	}
	pat := regexp.MustCompile("[^a-zA-Z0-9]+")
	for _, tt := range tests {
		name := pat.ReplaceAllString(tt.statement, "_")
		t.Run(name, func(t *testing.T) {
			ast, err := parseRawStatement(tt.statement)
			if (err != nil) != (tt.wantErr || tt.wantErrContaining != "") {
				t.Errorf("parseStatement(%s) error = %v, wantErr %v, wantErrContaining %v", tt.statement, err, tt.wantErr, tt.wantErrContaining)
				t.Errorf("AST: %+v", ast)
				return
			}
			if tt.wantErrContaining != "" {
				require.ErrorContains(t, err, tt.wantErrContaining)
			}
		})
	}
}

// This test doesn't validate parser results, simply checks whether the parse succeeds or not.
// It's a fast way to check a large range of possible syntaxes.
func Test_parseCondition(t *testing.T) {
	converterNameErrorPrefix := "converter names must start with an uppercase letter"
	editorWithIndexErrorPrefix := "only paths and converters may be indexed"

	tests := []struct {
		condition         string
		wantErr           bool
		wantErrContaining string
	}{
		{condition: `set(`, wantErr: true},
		{condition: `set("foo)`, wantErr: true},
		{condition: `set(name.)`, wantErr: true},
		{condition: `("foo")`, wantErr: true},
		{condition: `name =||= "fido"`, wantErr: true},
		{condition: `name = "fido"`, wantErr: true},
		{condition: `name or "fido"`, wantErr: true},
		{condition: `name and "fido"`, wantErr: true},
		{condition: `name and`, wantErr: true},
		{condition: `name or`, wantErr: true},
		{condition: `(`, wantErr: true},
		{condition: `)`, wantErr: true},
		{condition: `(name == "fido"))`, wantErr: true},
		{condition: `((name == "fido")`, wantErr: true},
		{condition: `set()`, wantErr: true},
		{condition: `Int() == 1`},
		{condition: `1 == Int()`},
		{condition: `true and 1 == Int() `},
		{condition: `false or 1 == Int() `},
		{condition: `service == "pinger" or foo.attributes["endpoint"] == "/x/alive"`},
		{condition: `service == "pinger" or foo.attributes["verb"] == "GET" and foo.attributes["endpoint"] == "/x/alive"`},
		{condition: `animal > "cat"`},
		{condition: `animal >= "cat"`},
		{condition: `animal <= "cat"`},
		{condition: `animal < "cat"`},
		{condition: `animal =< "dog"`, wantErr: true},
		{condition: `animal => "dog"`, wantErr: true},
		{condition: `animal <> "dog"`, wantErr: true},
		{condition: `animal = "dog"`, wantErr: true},
		{condition: `animal`, wantErr: true},
		{condition: `animal ==`, wantErr: true},
		{condition: `==`, wantErr: true},
		{condition: `== animal`, wantErr: true},
		{condition: `attributes["path"] == "/healthcheck"`},
		{condition: `One() == 1`},
		{condition: `test(fail())`, wantErr: true},
		{condition: `Test()`},
		{condition: `"test" == Foo`, wantErr: true},
		{condition: `test(animal) == "dog"`, wantErrContaining: converterNameErrorPrefix},
		{condition: `test(animal)["kind"] == "birds"`, wantErrContaining: converterNameErrorPrefix},
		{condition: `test(animal)["kind"] == "birds"`, wantErrContaining: editorWithIndexErrorPrefix},
	}
	pat := regexp.MustCompile("[^a-zA-Z0-9]+")
	for _, tt := range tests {
		name := pat.ReplaceAllString(tt.condition, "_")
		t.Run(name, func(t *testing.T) {
			ast, err := parseRawCondition(tt.condition)
			if (err != nil) != (tt.wantErr || tt.wantErrContaining != "") {
				t.Errorf("parseCondition(%s) error = %v, wantErr %v", tt.condition, err, tt.wantErr)
				t.Errorf("AST: %+v", ast)
				return
			}
			if tt.wantErrContaining != "" {
				require.ErrorContains(t, err, tt.wantErrContaining)
			}
		})
	}
}

// This test doesn't validate parser results, simply checks whether the parse succeeds or not.
// It's a fast way to check a large range of possible syntaxes.
func Test_parseValueExpression(t *testing.T) {
	converterNameErrorPrefix := "converter names must start with an uppercase letter"
	editorWithIndexErrorPrefix := "only paths and converters may be indexed"

	tests := []struct {
		valueExpression   string
		wantErr           bool
		wantErrContaining string
	}{
		{valueExpression: `time_end - time_end`},
		{valueExpression: `time_end - time_end - attributes["foo"]`},
		{valueExpression: `Test("foo")`},
		{valueExpression: `Test(Test("foo")) - attributes["bar"]`},
		{valueExpression: `Test(Test("foo")) - attributes["bar"]"`, wantErr: true},
		{valueExpression: `test("foo")`, wantErr: true, wantErrContaining: converterNameErrorPrefix},
		{valueExpression: `test(animal)["kind"]`, wantErrContaining: editorWithIndexErrorPrefix},
		{valueExpression: `Test("a"")foo"`, wantErr: true},
		{valueExpression: `Test("a"") == 1"`, wantErr: true},
	}
	pat := regexp.MustCompile("[^a-zA-Z0-9]+")
	for _, tt := range tests {
		name := pat.ReplaceAllString(tt.valueExpression, "_")
		t.Run(name, func(t *testing.T) {
			ast, err := parseRawValue(tt.valueExpression)
			if (err != nil) != (tt.wantErr || tt.wantErrContaining != "") {
				t.Errorf("parseCondition(%s) error = %v, wantErr %v", tt.valueExpression, err, tt.wantErr)
				t.Errorf("AST: %+v", ast)
				return
			}
			if tt.wantErrContaining != "" {
				require.ErrorContains(t, err, tt.wantErrContaining)
			}
		})
	}
}
