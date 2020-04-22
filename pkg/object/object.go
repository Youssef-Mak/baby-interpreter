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
	STRING_OBJ     = "STRING"
	NULL_OBJ       = "NULL"
	RETURN_VAL_OBJ = "RETURN_VAL"
	FUNCTION_OBJ   = "FUNCTION"
	BUILTIN_OBJ    = "BUILTIN"
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

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

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

type BuiltInFunction func(args ...Object) Object
type BuiltIn struct {
	Func BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType { return BUILTIN_OBJ }
func (b *BuiltIn) Inspect() string  { return "BuiltIn Function" }

type Error struct {
	Message string
}

func (err *Error) Type() ObjectType { return ERROR_OBJ }
func (err *Error) Inspect() string  { return err.Message }
