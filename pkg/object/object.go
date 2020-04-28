package object

import (
	"bytes"
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"hash/fnv"
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
	ARRAY_OBJ      = "ARRAY"
	HASH_OBJ       = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// TODO: Cache Hash Values for all implementations
type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// TODO: Handle hash collision(open addressing or seperate chaining)
func (s *String) HashKey() HashKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: hash.Sum64()}
}

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashEntry struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashEntry
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairsMsg := []string{}
	for _, entry := range h.Pairs {
		pairsMsg = append(pairsMsg, entry.Key.Inspect()+" : "+entry.Value.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairsMsg, ", "))
	out.WriteString("}")
	return out.String()

}

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
