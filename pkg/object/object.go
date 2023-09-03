package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/AzraelSec/cube/pkg/ast"
)

type ObjectType string
type BuiltinFunction func(...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	// note: introduce caching on return values of the implementations
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType  { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type String struct {
	Value string
}

func (*String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) HashKey() HashKey {
	var value uint64 = 1
	if i.Value == false {
		value = 0
	}
	return HashKey{Type: i.Type(), Value: value}
}

type Null struct{}

func (*Null) Type() ObjectType { return NULL_OBJ }
func (*Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (*ReturnValue) Type() ObjectType   { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Error struct {
	Msg string
}

func (*Error) Type() ObjectType  { return ERROR_OBJ }
func (e *Error) Inspect() string { return fmt.Sprintf("Error: %s", e.Msg) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (*Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var buff bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	buff.WriteString("fn(")
	buff.WriteString(strings.Join(params, ", "))
	buff.WriteString(") {\n")
	buff.WriteString(f.Body.String())
	buff.WriteString("\n}")

	return buff.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (*Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (*Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var buff bytes.Buffer

	elems := make([]string, len(a.Elements))
	for i := 0; i < len(a.Elements); i++ {
		elems[i] = a.Elements[i].Inspect()
	}

	buff.WriteString("[")
	buff.WriteString(strings.Join(elems, ", "))
	buff.WriteString("]")

	return buff.String()
}

type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (*Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var buff bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	buff.WriteString("{")
	buff.WriteString(strings.Join(pairs, ", "))
	buff.WriteString("}")

	return buff.String()
}
