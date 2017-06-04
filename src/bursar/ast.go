package main

import (
	"fmt"
)

type ValueType int
type State map[string]interface{}

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
    Type string
    StringValue string
    CharValue rune
    BoolValue bool
    IntValue int
    FloatValue float64
    ValueType ValueType
    Subvalues []Ast
}

type Func struct {
    Arguments []Ast
    ValueType ValueType
    Subvalues []Ast
}

type Call struct {
    Module Ast
    Function Ast
    Arguments []Ast
    ValueType ValueType
}

type If struct {
    Condition Ast
    Then []Ast
    Else []Ast
}

type Assignment struct {
    Left []Ast
    Right Ast
}

type Ast interface {
    Print(indent int) string
    Compile(state State) string
}

func (a BasicAst) String() string {
    switch (a.ValueType){
        case STRING:
            if a.Type == "Comment" {
                return fmt.Sprintf("//%s", a.StringValue)
            } else {
                return fmt.Sprintf("\"%s\"", a.StringValue)
            }
        case CHAR:
            return fmt.Sprintf("'%s'", string(a.CharValue))
        case INT:
            return fmt.Sprintf("%d", a.IntValue)
        case FLOAT:
            return fmt.Sprintf("%f", a.FloatValue)
        case BOOL:
            if a.BoolValue {
                return "true"
            }
            return "false"
        case NIL:
            return "nil"
        case CONTAINER:
            values := ""
            for _, el := range(a.Subvalues){
                values += fmt.Sprint(el)
            }
            return values
    }
    return "()"
}

func (a BasicAst) Print(indent int) string {
    str := ""

    for i := 0; i < indent; i++ {
        str += "  "
    }
    str += fmt.Sprintf("%s %s:\n", a.Type, a)
    for _, el := range(a.Subvalues){
        str += el.Print(indent+1)
    }
    return str
}

func (a Func) String() string {
    return "Func"
}

func (i If) String() string {
    return "If"
}

func (a Func) Print(indent int) string {
    str := ""

    for i := 0; i < indent; i++ {
        str += "  "
    }
    str += "Func"
    if len(a.Arguments) > 0 {
        str += " (\n"
        for _, el := range(a.Arguments){
            str += el.Print(indent + 1)
        }
        for i := 0; i < indent; i++ {
            str += "  "
        }
        str += ")\n"
    }
    for _, el := range(a.Subvalues){
        str += el.Print(indent+1)
    }
    return str
}

func (a Call) Print(indent int) string {
    str := ""

    for i := 0; i < indent; i++ {
        str += "  "
    }
    str += "Call:\n"
    if a.Module != nil {
        str += a.Module.Print(indent + 1)
    }
    str += a.Function.Print(indent + 1)

    if len(a.Arguments) > 0 {
        for i := 0; i < indent; i++ {
            str += "  "
        }
        str += "(\n"
        for _, el := range(a.Arguments){
            str += el.Print(indent + 1)
        }
        for i := 0; i < indent; i++ {
            str += "  "
        }
        str += ")\n"
    }
    return str
}

func (i If) Print(indent int) string {
    str := ""

    for i := 0; i < indent; i++ {
        str += "  "
    }
    str += "If"
    if i.Condition != nil {
        str += ":\n"
        str += i.Condition.Print(indent + 1)

    }
    for _, el := range(i.Then){
        for i := 0; i < indent; i++ {
            str += "  "
        }
        str += "Then:\n"
        str += el.Print(indent+1)
    }
    for _, el := range(i.Else){
        for i := 0; i < indent; i++ {
            str += "  "
        }
        str += "Else:\n"
        str += el.Print(indent+1)

    }
    return str
}

func (a Assignment) Print(indent int) string {
    str := ""

    for i := 0; i < indent; i++ {
        str += "  "
    }
    str += "Assignment:\n"

    for _, el := range(a.Left){
        str += el.Print(indent+1)
    }
    str += a.Right.Print(indent+1)

    return str
}

