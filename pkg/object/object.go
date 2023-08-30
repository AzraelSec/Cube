package object

import (
	"bytes"
	"fmt"
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
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType  { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (*String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }
func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }

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
