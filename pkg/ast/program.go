package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var buff bytes.Buffer

	for _, stm := range p.Statements {
		buff.WriteString(stm.String())
	}

	return buff.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
