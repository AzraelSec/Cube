package ast

import (
	"testing"

	"github.com/AzraelSec/cube/pkg/token"
)

func TestString(t *testing.T) {
	p := Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myvar"},
					Value: "myvar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anothervar"},
					Value: "anothervar",
				},
			},
		},
	}

	if p.String() != "let myvar = anothervar;" {
		t.Errorf("program.String() wrong. got=%q", p.String())
	}
}
