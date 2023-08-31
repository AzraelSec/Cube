package ast

import (
	"bytes"
	"strings"

	"github.com/AzraelSec/cube/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Expressions
type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (*Identifier) expressionNode()        {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (*IntegerLiteral) expressionNode()         {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token // token.TRUE, token.FALSE
	Value bool
}

func (*Boolean) expressionNode()        {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (*StringLiteral) expressionNode()        {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (*ArrayLiteral) expressionNode()         {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var buff bytes.Buffer

	elems := make([]string, len(al.Elements))
	for idx, e := range al.Elements {
		elems[idx] = e.String()
	}

	buff.WriteString("[")
	buff.WriteString(strings.Join(elems, ", "))
	buff.WriteString("]")

	return buff.String()
}

type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression  // the array
	Index Expression  // the integer (or derived) index
}

func (*IndexExpression) expressionNode()         {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var buff bytes.Buffer

	buff.WriteString("(")
	buff.WriteString(ie.Left.String())
	buff.WriteString("[")
	buff.WriteString(ie.Index.String())
	buff.WriteString("]")
	buff.WriteString(")")

	return buff.String()
}

type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (*IfExpression) expressionNode()         {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var buff bytes.Buffer

	buff.WriteString("if")
	buff.WriteString(ie.Condition.String())
	buff.WriteString(" ")
	buff.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		buff.WriteString("else ")
		buff.WriteString(ie.Alternative.String())
	}

	return buff.String()
}

type PrefixExpression struct {
	Token    token.Token // token.BANG, token.MINUS
	Operator string
	Right    Expression
}

func (*PrefixExpression) expressionNode()         {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var buff bytes.Buffer

	buff.WriteString("(")
	buff.WriteString(pe.Operator)
	buff.WriteString(pe.Right.String())
	buff.WriteString(")")

	return buff.String()
}

type InfixExpression struct {
	Token    token.Token // token.ADD, token.MINUS,...
	Left     Expression
	Operator string
	Right    Expression
}

func (*InfixExpression) expressionNode()         {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var buff bytes.Buffer

	buff.WriteString("(")
	buff.WriteString(ie.Left.String())
	buff.WriteString(" " + ie.Operator + " ")
	buff.WriteString(ie.Right.String())
	buff.WriteString(")")

	return buff.String()
}

type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Parameters []*Identifier
	Body       *BlockStatement
}

func (FunctionLiteral) expressionNode()          {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var buff bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	buff.WriteString(fl.TokenLiteral())
	buff.WriteString("(")
	buff.WriteString(strings.Join(params, ", "))
	buff.WriteString(")")
	buff.WriteString(fl.Body.String())

	return buff.String()
}

type CallExpression struct {
	Token    token.Token // token.LPAREN
	Function Expression  // Identifier || FunctionLiteral
	Args     []Expression
}

func (CallExpression) expressionNode()          {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var buff bytes.Buffer

	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())
	}

	buff.WriteString(ce.Function.String())
	buff.WriteString("(")
	buff.WriteString(strings.Join(args, ", "))
	buff.WriteString(")")

	return buff.String()
}

// Statements
type LetStatement struct {
	Token token.Token // token.TOKEN token
	Name  *Identifier
	Value Expression
}

func (*LetStatement) statementNode()          {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var buff bytes.Buffer

	buff.WriteString(ls.TokenLiteral())
	buff.WriteString(" ")
	buff.WriteString(ls.Name.String())
	buff.WriteString(" = ")

	if ls.Value != nil {
		buff.WriteString(ls.Value.String())
	}

	buff.WriteString(";")
	return buff.String()
}

type ReturnStatement struct {
	Token    token.Token // token.RETURN
	RetValue Expression
}

func (*ReturnStatement) statementNode()          {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var buff bytes.Buffer

	buff.WriteString(rs.TokenLiteral())
	buff.WriteString(" ")

	if rs.RetValue != nil {
		buff.WriteString(rs.RetValue.String())
	}
	buff.WriteString(";")

	return buff.String()
}

// note: needed to handle statement that simply are expressions
// ex: let x = 5; x + 10;
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (*ExpressionStatement) statementNode()          {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (BlockStatement) statementNode()           {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var buff bytes.Buffer
	for _, stm := range bs.Statements {
		buff.WriteString(stm.String())
	}
	return buff.String()
}
