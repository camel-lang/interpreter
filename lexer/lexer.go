package lexer

import (
	"camel/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()
	return lex
}

func (lexer *Lexer) readChar() {

	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lex *Lexer) NextToken() token.Token {

	var tok token.Token
	lex.eatSpace()

	switch lex.char {

	case '=':
		if lex.peekChar() == '=' {
			tok = token.Token{Type: token.EQ, Literal: "=="}
			lex.readChar()
		} else {
			tok = newToken(token.ASSIGN, lex.char)
		}
	case '+':
		tok = newToken(token.PLUS, lex.char)
	case '-':
		tok = newToken(token.MINUS, lex.char)
	case '*':
		tok = newToken(token.ASTERISK, lex.char)
	case '/':
		tok = newToken(token.SLASH, lex.char)
	case '<':
		tok = newToken(token.LT, lex.char)
	case '>':
		tok = newToken(token.GT, lex.char)
	case ';':
		tok = newToken(token.SEMICOLON, lex.char)
	case ',':
		tok = newToken(token.COMMA, lex.char)
	case '(':
		tok = newToken(token.LPAREN, lex.char)
	case ')':
		tok = newToken(token.RPAREN, lex.char)
	case '{':
		tok = newToken(token.LBRACE, lex.char)
	case '}':
		tok = newToken(token.RBRACE, lex.char)
	case '!':
		if lex.peekChar() == '=' {
			lex.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, lex.char)
		}
	case '"' : 
		tok.Type = token.STRING 
		tok.Literal = lex.readString() 
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lex.char) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)

			return tok

		} else if isDigit(lex.char) {
			tok.Literal = lex.readNumber()
			tok.Type = token.INT

			return tok

		} else {

			tok.Literal = ""
			tok.Type = token.ILLEGAL
		}
	}

	lex.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (lex *Lexer) readIdentifier() string {

	pos := lex.position
	for isLetter(lex.char) {
		lex.readChar()
	}
	return lex.input[pos:lex.position]

}
func (lex *Lexer) readString() string { 
	
	pos := lex.position + 1 
	for { 
		lex.readChar() 
		if lex.char == '"' || lex.char == 0 { 
			break 
		}
	} 

	return lex.input[pos:lex.position]
}

func (lex *Lexer) readNumber() string {
	pos := lex.position
	for isDigit(lex.char) {
		lex.readChar()
	}
	return lex.input[pos:lex.position]
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
func isLetter(ch byte) bool {

	return 'A' <= ch && ch <= 'Z' ||
		'a' <= ch && ch <= 'z' ||
		ch == '_'
}

func (lex *Lexer) eatSpace() {

	for lex.char == ' ' ||
		lex.char == '\n' ||
		lex.char == '\t' ||
		lex.char == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) peekChar() byte {

	if lex.readPosition >= len(lex.input) {
		return 0
	}
	return lex.input[lex.readPosition]
}
