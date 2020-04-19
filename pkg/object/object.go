package object

import (
	"bytes"
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ    = "INTEGER"
	BOOLEAN_OBJ    = "BOOLEAN"
	NULL_OBJ       = "NULL"
	RETURN_VAL_OBJ = "RETURN_VAL"
	FUNCTION_OBJ   = "FUNCTION"
	ERROR_OBJ      = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VAL_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fun *Function) Type() ObjectType { return FUNCTION_OBJ }
func (fun *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fun.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fun.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return err.Message }
