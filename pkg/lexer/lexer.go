package lexer

import "github.com/AzraelSec/cube/pkg/token"

const nul = 0

type Lexer struct {
	input        string
	position     int  // last index of input already tokenized
	readPosition int  // index of input to read
	ch           byte // input[readPosition]
}

func New(s string) *Lexer {
	l := &Lexer{input: s}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = nul
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// todo: add special chars handling + other stuff
func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tkn token.Token

	l.skipWhiteSpaces()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = token.New(token.ASSIGN, string(l.ch))
		}
	case '+':
		tkn = token.New(token.PLUS, string(l.ch))
	case '-':
		tkn = token.New(token.MINUS, string(l.ch))
	case '*':
		tkn = token.New(token.ASTERISK, string(l.ch))
	case '/':
		tkn = token.New(token.SLASH, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tkn = token.Token{Type: token.NE, Literal: string(ch) + string(l.ch)}
		} else {
			tkn = token.New(token.BANG, string(l.ch))
		}
	case '<':
		tkn = token.New(token.LT, string(l.ch))
	case '>':
		tkn = token.New(token.GT, string(l.ch))
	case '(':
		tkn = token.New(token.LPAREN, string(l.ch))
	case ')':
		tkn = token.New(token.RPAREN, string(l.ch))
	case '{':
		tkn = token.New(token.LBRACE, string(l.ch))
	case '}':
		tkn = token.New(token.RBRACE, string(l.ch))
	case ',':
		tkn = token.New(token.COMMA, string(l.ch))
	case ';':
		tkn = token.New(token.SEMICOLON, string(l.ch))
	case '"':
		tkn = token.New(token.STRING, l.readString())
	case '[':
		tkn = token.New(token.LBRACKET, string(l.ch))
	case ']':
		tkn = token.New(token.RBRACKET, string(l.ch))
	case nul:
		tkn.Literal = ""
		tkn.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tkn.Literal = l.readIdentifier()
			tkn.Type = token.LookupIdent(tkn.Literal)
			return tkn // note: we don't want to call `readChar` again
		}
		if isDigit(l.ch) {
			tkn.Type = token.INT
			tkn.Literal = l.readNumber()
			return tkn
		}
		tkn = token.New(token.ILLEGAL, string(l.ch))
	}

	l.readChar()
	return tkn
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
