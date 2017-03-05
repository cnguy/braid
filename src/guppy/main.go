package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type ValueType int

const (
	STRING = iota
	INT
	FLOAT
	BOOL
	CHAR
	CONTAINER
	NIL
)

type BasicAst struct {
	Type        string
	StringValue string
	CharValue   rune
	BoolValue   bool
	IntValue    int
	FloatValue  float64
	ValueType   ValueType
	Subvalues   []Ast
}

type Func struct {
	Arguments []Ast
	ValueType ValueType
	Subvalues []Ast
}

type Ast interface {
	Print(indent int) string
}

func main() {
	input := `
let a = -4 + 55.0 > 99 or "hi" ++ 'm';
let b = "cheese" ++ "ham";
let c = 5.0;
34 + 5;
let cheesy = fun item -> {
    item ++ " with cheese";
}
`
	fmt.Println(input)
	r := strings.NewReader(input)
	result, err := ParseReader("", r)
	ast := result.(BasicAst)
	fmt.Println("=", ast.Print(0))
	fmt.Println(err)
}

func (a BasicAst) String() string {
	switch a.ValueType {
	case STRING:
		return fmt.Sprintf("\"%s\"", a.StringValue)
	case CHAR:
		return fmt.Sprintf("'%s'", string(a.CharValue))
	case INT:
		return fmt.Sprintf("%d", a.IntValue)
	case FLOAT:
		return fmt.Sprintf("%f", a.FloatValue)
	}
	return "()"
}

func (a BasicAst) Print(indent int) string {
	str := ""

	for i := 0; i < indent; i++ {
		str += "  "
	}
	str += fmt.Sprintf("%s %s:\n", a.Type, a)
	for _, el := range a.Subvalues {
		str += el.Print(indent + 1)
	}
	return str
}

func (a Func) String() string {
	return "fun"
}

func (a Func) Print(indent int) string {
	str := ""

	for i := 0; i < indent; i++ {
		str += "  "
	}
	str += "fun"
	if len(a.Arguments) > 0 {
		str += " ("
		for _, el := range a.Arguments {
			str += el.Print(indent)
		}
		str += ") "
	}
	for _, el := range a.Subvalues {
		str += el.Print(indent + 1)
	}
	return str
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Module",
			pos:  position{line: 115, col: 1, offset: 2018},
			expr: &actionExpr{
				pos: position{line: 115, col: 10, offset: 2027},
				run: (*parser).callonModule1,
				expr: &seqExpr{
					pos: position{line: 115, col: 10, offset: 2027},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 115, col: 10, offset: 2027},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 115, col: 12, offset: 2029},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 115, col: 17, offset: 2034},
								name: "Statement",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 115, col: 27, offset: 2044},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 115, col: 29, offset: 2046},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 115, col: 34, offset: 2051},
								expr: &ruleRefExpr{
									pos:  position{line: 115, col: 35, offset: 2052},
									name: "Statement",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 115, col: 47, offset: 2064},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 115, col: 49, offset: 2066},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 128, col: 1, offset: 2486},
			expr: &choiceExpr{
				pos: position{line: 128, col: 13, offset: 2498},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 128, col: 13, offset: 2498},
						run: (*parser).callonStatement2,
						expr: &seqExpr{
							pos: position{line: 128, col: 13, offset: 2498},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 128, col: 13, offset: 2498},
									val:        "let",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 19, offset: 2504},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 128, col: 22, offset: 2507},
									label: "i",
									expr: &ruleRefExpr{
										pos:  position{line: 128, col: 24, offset: 2509},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 35, offset: 2520},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 128, col: 37, offset: 2522},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 41, offset: 2526},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 128, col: 43, offset: 2528},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 128, col: 48, offset: 2533},
										name: "Expr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 128, col: 53, offset: 2538},
									name: "_",
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 131, col: 5, offset: 2712},
						name: "Expr",
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 133, col: 1, offset: 2718},
			expr: &choiceExpr{
				pos: position{line: 133, col: 8, offset: 2725},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 133, col: 8, offset: 2725},
						name: "FuncDefn",
					},
					&ruleRefExpr{
						pos:  position{line: 133, col: 19, offset: 2736},
						name: "BinExpr",
					},
				},
			},
		},
		{
			name: "FuncDefn",
			pos:  position{line: 135, col: 1, offset: 2745},
			expr: &actionExpr{
				pos: position{line: 135, col: 12, offset: 2756},
				run: (*parser).callonFuncDefn1,
				expr: &seqExpr{
					pos: position{line: 135, col: 12, offset: 2756},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 135, col: 12, offset: 2756},
							val:        "fun",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 18, offset: 2762},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 135, col: 21, offset: 2765},
							label: "ids",
							expr: &zeroOrMoreExpr{
								pos: position{line: 135, col: 25, offset: 2769},
								expr: &seqExpr{
									pos: position{line: 135, col: 26, offset: 2770},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 135, col: 26, offset: 2770},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 135, col: 37, offset: 2781},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 42, offset: 2786},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 135, col: 44, offset: 2788},
							val:        "->",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 49, offset: 2793},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 135, col: 51, offset: 2795},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 55, offset: 2799},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 135, col: 57, offset: 2801},
							label: "statements",
							expr: &oneOrMoreExpr{
								pos: position{line: 135, col: 68, offset: 2812},
								expr: &ruleRefExpr{
									pos:  position{line: 135, col: 69, offset: 2813},
									name: "Statement",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 81, offset: 2825},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 135, col: 83, offset: 2827},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 87, offset: 2831},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "BinExpr",
			pos:  position{line: 158, col: 1, offset: 3490},
			expr: &actionExpr{
				pos: position{line: 158, col: 11, offset: 3500},
				run: (*parser).callonBinExpr1,
				expr: &seqExpr{
					pos: position{line: 158, col: 11, offset: 3500},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 158, col: 11, offset: 3500},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 158, col: 13, offset: 3502},
							label: "op",
							expr: &ruleRefExpr{
								pos:  position{line: 158, col: 16, offset: 3505},
								name: "BinOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 158, col: 22, offset: 3511},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 158, col: 27, offset: 3516},
								expr: &ruleRefExpr{
									pos:  position{line: 158, col: 28, offset: 3517},
									name: "BinOp",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 36, offset: 3525},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 158, col: 38, offset: 3527},
							val:        ";",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 158, col: 42, offset: 3531},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "BinOp",
			pos:  position{line: 171, col: 1, offset: 3941},
			expr: &choiceExpr{
				pos: position{line: 171, col: 9, offset: 3949},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 171, col: 9, offset: 3949},
						name: "BinOpBool",
					},
					&ruleRefExpr{
						pos:  position{line: 171, col: 21, offset: 3961},
						name: "BinOpEquality",
					},
					&ruleRefExpr{
						pos:  position{line: 171, col: 37, offset: 3977},
						name: "BinOpLow",
					},
					&ruleRefExpr{
						pos:  position{line: 171, col: 48, offset: 3988},
						name: "BinOpHigh",
					},
				},
			},
		},
		{
			name: "BinOpBool",
			pos:  position{line: 173, col: 1, offset: 3999},
			expr: &actionExpr{
				pos: position{line: 173, col: 13, offset: 4011},
				run: (*parser).callonBinOpBool1,
				expr: &seqExpr{
					pos: position{line: 173, col: 13, offset: 4011},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 173, col: 13, offset: 4011},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 173, col: 15, offset: 4013},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 173, col: 21, offset: 4019},
								name: "BinOpEquality",
							},
						},
						&labeledExpr{
							pos:   position{line: 173, col: 35, offset: 4033},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 173, col: 40, offset: 4038},
								expr: &seqExpr{
									pos: position{line: 173, col: 41, offset: 4039},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 173, col: 41, offset: 4039},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 173, col: 44, offset: 4042},
											name: "OperatorBoolean",
										},
										&ruleRefExpr{
											pos:  position{line: 173, col: 60, offset: 4058},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 173, col: 63, offset: 4061},
											name: "BinOpEquality",
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
			name: "BinOpEquality",
			pos:  position{line: 191, col: 1, offset: 4639},
			expr: &actionExpr{
				pos: position{line: 191, col: 17, offset: 4655},
				run: (*parser).callonBinOpEquality1,
				expr: &seqExpr{
					pos: position{line: 191, col: 17, offset: 4655},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 191, col: 17, offset: 4655},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 191, col: 19, offset: 4657},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 191, col: 25, offset: 4663},
								name: "BinOpLow",
							},
						},
						&labeledExpr{
							pos:   position{line: 191, col: 34, offset: 4672},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 191, col: 39, offset: 4677},
								expr: &seqExpr{
									pos: position{line: 191, col: 40, offset: 4678},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 191, col: 40, offset: 4678},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 191, col: 43, offset: 4681},
											name: "OperatorEquality",
										},
										&ruleRefExpr{
											pos:  position{line: 191, col: 60, offset: 4698},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 191, col: 63, offset: 4701},
											name: "BinOpLow",
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
			name: "BinOpLow",
			pos:  position{line: 210, col: 1, offset: 5275},
			expr: &actionExpr{
				pos: position{line: 210, col: 12, offset: 5286},
				run: (*parser).callonBinOpLow1,
				expr: &seqExpr{
					pos: position{line: 210, col: 12, offset: 5286},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 210, col: 12, offset: 5286},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 210, col: 14, offset: 5288},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 210, col: 20, offset: 5294},
								name: "BinOpHigh",
							},
						},
						&labeledExpr{
							pos:   position{line: 210, col: 30, offset: 5304},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 210, col: 35, offset: 5309},
								expr: &seqExpr{
									pos: position{line: 210, col: 36, offset: 5310},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 210, col: 36, offset: 5310},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 210, col: 39, offset: 5313},
											name: "OperatorLow",
										},
										&ruleRefExpr{
											pos:  position{line: 210, col: 51, offset: 5325},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 210, col: 54, offset: 5328},
											name: "BinOpHigh",
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
			name: "BinOpHigh",
			pos:  position{line: 229, col: 1, offset: 5903},
			expr: &actionExpr{
				pos: position{line: 229, col: 13, offset: 5915},
				run: (*parser).callonBinOpHigh1,
				expr: &seqExpr{
					pos: position{line: 229, col: 13, offset: 5915},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 229, col: 13, offset: 5915},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 229, col: 15, offset: 5917},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 229, col: 21, offset: 5923},
								name: "Value",
							},
						},
						&labeledExpr{
							pos:   position{line: 229, col: 27, offset: 5929},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 229, col: 32, offset: 5934},
								expr: &seqExpr{
									pos: position{line: 229, col: 33, offset: 5935},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 229, col: 33, offset: 5935},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 229, col: 36, offset: 5938},
											name: "OperatorHigh",
										},
										&ruleRefExpr{
											pos:  position{line: 229, col: 49, offset: 5951},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 229, col: 52, offset: 5954},
											name: "Value",
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
			name: "Operator",
			pos:  position{line: 247, col: 1, offset: 6524},
			expr: &choiceExpr{
				pos: position{line: 247, col: 12, offset: 6535},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 247, col: 12, offset: 6535},
						name: "OperatorBoolean",
					},
					&ruleRefExpr{
						pos:  position{line: 247, col: 30, offset: 6553},
						name: "OperatorEquality",
					},
					&ruleRefExpr{
						pos:  position{line: 247, col: 49, offset: 6572},
						name: "OperatorHigh",
					},
					&ruleRefExpr{
						pos:  position{line: 247, col: 64, offset: 6587},
						name: "OperatorLow",
					},
				},
			},
		},
		{
			name: "OperatorBoolean",
			pos:  position{line: 249, col: 1, offset: 6600},
			expr: &actionExpr{
				pos: position{line: 249, col: 19, offset: 6618},
				run: (*parser).callonOperatorBoolean1,
				expr: &choiceExpr{
					pos: position{line: 249, col: 21, offset: 6620},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 249, col: 21, offset: 6620},
							val:        "not",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 249, col: 29, offset: 6628},
							val:        "or",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 249, col: 36, offset: 6635},
							val:        "and",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "OperatorEquality",
			pos:  position{line: 253, col: 1, offset: 6734},
			expr: &actionExpr{
				pos: position{line: 253, col: 20, offset: 6753},
				run: (*parser).callonOperatorEquality1,
				expr: &choiceExpr{
					pos: position{line: 253, col: 22, offset: 6755},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 253, col: 22, offset: 6755},
							val:        ">=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 253, col: 29, offset: 6762},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 253, col: 36, offset: 6769},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 253, col: 42, offset: 6775},
							val:        ">",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 253, col: 48, offset: 6781},
							val:        "===",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 253, col: 56, offset: 6789},
							val:        "==",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "OperatorHigh",
			pos:  position{line: 257, col: 1, offset: 6895},
			expr: &choiceExpr{
				pos: position{line: 257, col: 16, offset: 6910},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 257, col: 16, offset: 6910},
						run: (*parser).callonOperatorHigh2,
						expr: &choiceExpr{
							pos: position{line: 257, col: 18, offset: 6912},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 257, col: 18, offset: 6912},
									val:        "/.",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 257, col: 25, offset: 6919},
									val:        "*.",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 260, col: 3, offset: 7025},
						run: (*parser).callonOperatorHigh6,
						expr: &choiceExpr{
							pos: position{line: 260, col: 5, offset: 7027},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 260, col: 5, offset: 7027},
									val:        "*",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 260, col: 11, offset: 7033},
									val:        "/",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 260, col: 17, offset: 7039},
									val:        "^",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 263, col: 3, offset: 7142},
						run: (*parser).callonOperatorHigh11,
						expr: &litMatcher{
							pos:        position{line: 263, col: 3, offset: 7142},
							val:        "++",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "OperatorLow",
			pos:  position{line: 267, col: 1, offset: 7246},
			expr: &choiceExpr{
				pos: position{line: 267, col: 15, offset: 7260},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 267, col: 15, offset: 7260},
						run: (*parser).callonOperatorLow2,
						expr: &choiceExpr{
							pos: position{line: 267, col: 17, offset: 7262},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 267, col: 17, offset: 7262},
									val:        "+.",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 267, col: 24, offset: 7269},
									val:        "-.",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 270, col: 3, offset: 7375},
						run: (*parser).callonOperatorLow6,
						expr: &choiceExpr{
							pos: position{line: 270, col: 5, offset: 7377},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 270, col: 5, offset: 7377},
									val:        "+",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 270, col: 11, offset: 7383},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 274, col: 1, offset: 7485},
			expr: &choiceExpr{
				pos: position{line: 274, col: 9, offset: 7493},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 274, col: 9, offset: 7493},
						name: "Identifier",
					},
					&actionExpr{
						pos: position{line: 274, col: 22, offset: 7506},
						run: (*parser).callonValue3,
						expr: &labeledExpr{
							pos:   position{line: 274, col: 22, offset: 7506},
							label: "v",
							expr: &ruleRefExpr{
								pos:  position{line: 274, col: 24, offset: 7508},
								name: "Const",
							},
						},
					},
					&actionExpr{
						pos: position{line: 277, col: 3, offset: 7549},
						run: (*parser).callonValue6,
						expr: &seqExpr{
							pos: position{line: 277, col: 3, offset: 7549},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 277, col: 3, offset: 7549},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 277, col: 7, offset: 7553},
									label: "expr",
									expr: &ruleRefExpr{
										pos:  position{line: 277, col: 12, offset: 7558},
										name: "Expr",
									},
								},
								&litMatcher{
									pos:        position{line: 277, col: 17, offset: 7563},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 281, col: 1, offset: 7604},
			expr: &actionExpr{
				pos: position{line: 281, col: 14, offset: 7617},
				run: (*parser).callonIdentifier1,
				expr: &choiceExpr{
					pos: position{line: 281, col: 15, offset: 7618},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 281, col: 15, offset: 7618},
							exprs: []interface{}{
								&oneOrMoreExpr{
									pos: position{line: 281, col: 15, offset: 7618},
									expr: &charClassMatcher{
										pos:        position{line: 281, col: 15, offset: 7618},
										val:        "[a-z]",
										ranges:     []rune{'a', 'z'},
										ignoreCase: false,
										inverted:   false,
									},
								},
								&zeroOrMoreExpr{
									pos: position{line: 281, col: 22, offset: 7625},
									expr: &charClassMatcher{
										pos:        position{line: 281, col: 22, offset: 7625},
										val:        "[a-zA-Z0-9_]",
										chars:      []rune{'_'},
										ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 281, col: 38, offset: 7641},
							val:        "_",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 285, col: 1, offset: 7742},
			expr: &choiceExpr{
				pos: position{line: 285, col: 9, offset: 7750},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 285, col: 9, offset: 7750},
						run: (*parser).callonConst2,
						expr: &seqExpr{
							pos: position{line: 285, col: 9, offset: 7750},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 285, col: 9, offset: 7750},
									expr: &litMatcher{
										pos:        position{line: 285, col: 9, offset: 7750},
										val:        "-",
										ignoreCase: false,
									},
								},
								&oneOrMoreExpr{
									pos: position{line: 285, col: 14, offset: 7755},
									expr: &charClassMatcher{
										pos:        position{line: 285, col: 14, offset: 7755},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
								&notExpr{
									pos: position{line: 285, col: 21, offset: 7762},
									expr: &litMatcher{
										pos:        position{line: 285, col: 22, offset: 7763},
										val:        ".",
										ignoreCase: false,
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 292, col: 3, offset: 7939},
						run: (*parser).callonConst10,
						expr: &seqExpr{
							pos: position{line: 292, col: 3, offset: 7939},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 292, col: 3, offset: 7939},
									expr: &litMatcher{
										pos:        position{line: 292, col: 3, offset: 7939},
										val:        "-",
										ignoreCase: false,
									},
								},
								&oneOrMoreExpr{
									pos: position{line: 292, col: 8, offset: 7944},
									expr: &charClassMatcher{
										pos:        position{line: 292, col: 8, offset: 7944},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
								&litMatcher{
									pos:        position{line: 292, col: 15, offset: 7951},
									val:        ".",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 292, col: 19, offset: 7955},
									expr: &charClassMatcher{
										pos:        position{line: 292, col: 19, offset: 7955},
										val:        "[0-9]",
										ranges:     []rune{'0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 299, col: 3, offset: 8145},
						val:        "True",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 299, col: 12, offset: 8154},
						run: (*parser).callonConst20,
						expr: &litMatcher{
							pos:        position{line: 299, col: 12, offset: 8154},
							val:        "False",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 305, col: 3, offset: 8355},
						run: (*parser).callonConst22,
						expr: &litMatcher{
							pos:        position{line: 305, col: 3, offset: 8355},
							val:        "()",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 308, col: 3, offset: 8418},
						run: (*parser).callonConst24,
						expr: &seqExpr{
							pos: position{line: 308, col: 3, offset: 8418},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 308, col: 3, offset: 8418},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 308, col: 7, offset: 8422},
									expr: &seqExpr{
										pos: position{line: 308, col: 8, offset: 8423},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 308, col: 8, offset: 8423},
												expr: &ruleRefExpr{
													pos:  position{line: 308, col: 9, offset: 8424},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 308, col: 21, offset: 8436,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 308, col: 25, offset: 8440},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 315, col: 3, offset: 8624},
						run: (*parser).callonConst33,
						expr: &seqExpr{
							pos: position{line: 315, col: 3, offset: 8624},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 315, col: 3, offset: 8624},
									val:        "'",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 315, col: 7, offset: 8628},
									label: "val",
									expr: &seqExpr{
										pos: position{line: 315, col: 12, offset: 8633},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 315, col: 12, offset: 8633},
												expr: &ruleRefExpr{
													pos:  position{line: 315, col: 13, offset: 8634},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 315, col: 25, offset: 8646,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 315, col: 28, offset: 8649},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 319, col: 1, offset: 8740},
			expr: &charClassMatcher{
				pos:        position{line: 319, col: 15, offset: 8754},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 321, col: 1, offset: 8770},
			expr: &choiceExpr{
				pos: position{line: 321, col: 18, offset: 8787},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 321, col: 18, offset: 8787},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 321, col: 37, offset: 8806},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 323, col: 1, offset: 8821},
			expr: &charClassMatcher{
				pos:        position{line: 323, col: 20, offset: 8840},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 325, col: 1, offset: 8853},
			expr: &charClassMatcher{
				pos:        position{line: 325, col: 16, offset: 8868},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 327, col: 1, offset: 8875},
			expr: &charClassMatcher{
				pos:        position{line: 327, col: 23, offset: 8897},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 329, col: 1, offset: 8904},
			expr: &charClassMatcher{
				pos:        position{line: 329, col: 12, offset: 8915},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name:        "__",
			displayName: "\"sigwhitespace\"",
			pos:         position{line: 331, col: 1, offset: 8926},
			expr: &oneOrMoreExpr{
				pos: position{line: 331, col: 22, offset: 8947},
				expr: &charClassMatcher{
					pos:        position{line: 331, col: 22, offset: 8947},
					val:        "[ \\r\\n\\t]",
					chars:      []rune{' ', '\r', '\n', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 333, col: 1, offset: 8959},
			expr: &zeroOrMoreExpr{
				pos: position{line: 333, col: 18, offset: 8976},
				expr: &charClassMatcher{
					pos:        position{line: 333, col: 18, offset: 8976},
					val:        "[ \\r\\n\\t]",
					chars:      []rune{' ', '\r', '\n', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 335, col: 1, offset: 8988},
			expr: &notExpr{
				pos: position{line: 335, col: 7, offset: 8994},
				expr: &anyMatcher{
					line: 335, col: 8, offset: 8995,
				},
			},
		},
	},
}

func (c *current) onModule1(expr, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{expr.(BasicAst)}
		for _, el := range vals {
			subvalues = append(subvalues, el.(BasicAst))
		}
		return BasicAst{Type: "Module", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return BasicAst{Type: "Module", Subvalues: []Ast{expr.(BasicAst)}, ValueType: CONTAINER}, nil
	}
}

func (p *parser) callonModule1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onModule1(stack["expr"], stack["rest"])
}

func (c *current) onStatement2(i, expr interface{}) (interface{}, error) {
	fmt.Printf("assignment: %s\n", string(c.text))
	return BasicAst{Type: "Assignment", Subvalues: []Ast{i.(BasicAst), expr.(BasicAst)}, ValueType: CONTAINER}, nil
}

func (p *parser) callonStatement2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement2(stack["i"], stack["expr"])
}

func (c *current) onFuncDefn1(ids, statements interface{}) (interface{}, error) {
	fmt.Println("func", string(c.text))
	subvalues := []Ast{}
	args := []Ast{}
	vals := statements.([]interface{})
	if len(vals) > 0 {
		for _, el := range vals {
			subvalues = append(subvalues, el.(BasicAst))
		}
	}
	vals = ids.([]interface{})
	if len(vals) > 0 {
		restSl := toIfaceSlice(ids)
		for _, v := range restSl {
			// we can get each item in the grammar by index
			restExpr := toIfaceSlice(v)
			v := restExpr[0].(BasicAst)
			args = append(args, v)
		}
	}
	return Func{Arguments: args, Subvalues: subvalues, ValueType: CONTAINER}, nil
}

func (p *parser) callonFuncDefn1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFuncDefn1(stack["ids"], stack["statements"])
}

func (c *current) onBinExpr1(op, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{op.(BasicAst)}
		for _, el := range vals {
			subvalues = append(subvalues, el.(BasicAst))
		}
		return BasicAst{Type: "Expr", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return BasicAst{Type: "Expr", Subvalues: []Ast{op.(BasicAst)}, ValueType: CONTAINER}, nil
	}
}

func (p *parser) callonBinExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinExpr1(stack["op"], stack["rest"])
}

func (c *current) onBinOpBool1(first, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{first.(BasicAst)}
		restSl := toIfaceSlice(rest)
		for _, v := range restSl {
			// we can get each item in the grammar by index
			restExpr := toIfaceSlice(v)
			v := restExpr[3].(BasicAst)
			op := restExpr[1].(BasicAst)
			subvalues = append(subvalues, op, v)
		}
		return BasicAst{Type: "BinOp", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return first.(BasicAst), nil
	}
}

func (p *parser) callonBinOpBool1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOpBool1(stack["first"], stack["rest"])
}

func (c *current) onBinOpEquality1(first, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{first.(BasicAst)}
		restSl := toIfaceSlice(rest)
		for _, v := range restSl {
			// we can get each item in the grammar by index
			restExpr := toIfaceSlice(v)
			v := restExpr[3].(BasicAst)
			op := restExpr[1].(BasicAst)
			subvalues = append(subvalues, op, v)
		}
		return BasicAst{Type: "BinOp", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return first.(BasicAst), nil
	}

}

func (p *parser) callonBinOpEquality1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOpEquality1(stack["first"], stack["rest"])
}

func (c *current) onBinOpLow1(first, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{first.(BasicAst)}
		restSl := toIfaceSlice(rest)
		for _, v := range restSl {
			// we can get each item in the grammar by index
			restExpr := toIfaceSlice(v)
			v := restExpr[3].(BasicAst)
			op := restExpr[1].(BasicAst)
			subvalues = append(subvalues, op, v)
		}
		return BasicAst{Type: "BinOp", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return first.(BasicAst), nil
	}

}

func (p *parser) callonBinOpLow1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOpLow1(stack["first"], stack["rest"])
}

func (c *current) onBinOpHigh1(first, rest interface{}) (interface{}, error) {
	vals := rest.([]interface{})
	if len(vals) > 0 {
		subvalues := []Ast{first.(BasicAst)}
		restSl := toIfaceSlice(rest)
		for _, v := range restSl {
			// we can get each item in the grammar by index
			restExpr := toIfaceSlice(v)
			v := restExpr[3].(BasicAst)
			op := restExpr[1].(BasicAst)
			subvalues = append(subvalues, op, v)
		}
		return BasicAst{Type: "BinOp", Subvalues: subvalues, ValueType: CONTAINER}, nil
	} else {
		return first.(BasicAst), nil
	}
}

func (p *parser) callonBinOpHigh1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOpHigh1(stack["first"], stack["rest"])
}

func (c *current) onOperatorBoolean1() (interface{}, error) {
	return BasicAst{Type: "BoolOp", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorBoolean1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorBoolean1()
}

func (c *current) onOperatorEquality1() (interface{}, error) {
	return BasicAst{Type: "EqualityOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorEquality1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorEquality1()
}

func (c *current) onOperatorHigh2() (interface{}, error) {
	return BasicAst{Type: "FloatOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorHigh2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorHigh2()
}

func (c *current) onOperatorHigh6() (interface{}, error) {
	return BasicAst{Type: "IntOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorHigh6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorHigh6()
}

func (c *current) onOperatorHigh11() (interface{}, error) {
	return BasicAst{Type: "StringOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorHigh11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorHigh11()
}

func (c *current) onOperatorLow2() (interface{}, error) {
	return BasicAst{Type: "FloatOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorLow2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorLow2()
}

func (c *current) onOperatorLow6() (interface{}, error) {
	return BasicAst{Type: "IntOperator", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonOperatorLow6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOperatorLow6()
}

func (c *current) onValue3(v interface{}) (interface{}, error) {
	return v.(BasicAst), nil
}

func (p *parser) callonValue3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onValue3(stack["v"])
}

func (c *current) onValue6(expr interface{}) (interface{}, error) {
	return expr.(BasicAst), nil
}

func (p *parser) callonValue6() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onValue6(stack["expr"])
}

func (c *current) onIdentifier1() (interface{}, error) {
	return BasicAst{Type: "Identifier", StringValue: string(c.text), ValueType: STRING}, nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onConst2() (interface{}, error) {
	val, err := strconv.Atoi(string(c.text))
	if err != nil {
		return nil, err
	}
	return BasicAst{Type: "Integer", IntValue: val, ValueType: INT}, nil
}

func (p *parser) callonConst2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst2()
}

func (c *current) onConst10() (interface{}, error) {
	val, err := strconv.ParseFloat(string(c.text), 64)
	if err != nil {
		return nil, err
	}
	return BasicAst{Type: "Float", FloatValue: val, ValueType: FLOAT}, nil
}

func (p *parser) callonConst10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst10()
}

func (c *current) onConst20() (interface{}, error) {
	if string(c.text) == "True" {
		return BasicAst{Type: "Bool", BoolValue: true, ValueType: BOOL}, nil
	}
	return BasicAst{Type: "Bool", BoolValue: false, ValueType: BOOL}, nil
}

func (p *parser) callonConst20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst20()
}

func (c *current) onConst22() (interface{}, error) {
	return BasicAst{Type: "Nil", ValueType: NIL}, nil
}

func (p *parser) callonConst22() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst22()
}

func (c *current) onConst24() (interface{}, error) {
	val, err := strconv.Unquote(string(c.text))
	if err == nil {
		return BasicAst{Type: "String", StringValue: val, ValueType: STRING}, nil
	}
	return nil, err
}

func (p *parser) callonConst24() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst24()
}

func (c *current) onConst33(val interface{}) (interface{}, error) {
	return BasicAst{Type: "Char", CharValue: rune(c.text[1]), ValueType: CHAR}, nil
}

func (p *parser) callonConst33() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst33(stack["val"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
