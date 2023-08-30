package token

const (
	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifiers
	IDENT = "IDENT"
	INT   = "INT"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	EQ = "=="
	NE = "!="

	LT = "<"
	GT = ">"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// keywords
	FUNCTION = "fn"
	LET      = "let"
	IF       = "if"
	ELSE     = "else"
	TRUE     = "true"
	FALSE    = "false"
	RETURN   = "return"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(literal),
	}
}

func LookupIdent(s string) TokenType {
	if v, ok := keywords[s]; ok {
		return v
	}
	return IDENT
}
