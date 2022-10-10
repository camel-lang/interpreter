package lexer

import "xxx/token"

type Lexer struct { 

	input string 
	position int 
	readPosition int 
	char byte 

}

func New(input string) *Lexer { 
	return &Lexer{input: input} 
}

func (lexer *Lexer) readChar() { 
	
	if lexer.readPosition >= len(lexer.input) { 
		lexer.char = 0 
	} else { lexer.char = lexer.input[lexer.readPosition] } 
	lexer.position = lexer.readPosition 
	lexer.readPosition += 1 
}	

func (lexer *Lexer) NextToken() token.Token {

	var tok token.Token 
	lexer.readChar() 

	switch lexer.char { 

		case '=' : 
			tok = newToken(token.ASSIGN , lexer.char)
		case '+' : 
			tok = newToken(token.PLUS , lexer.char)
		case ';' : 
			tok = newToken(token.SEMICOLON , lexer.char)
		case '(' : 
			tok = newToken(token.LPAREN , lexer.char)
		case ')' :
			tok = newToken(token.RPAREN , lexer.char)
		case '{' : 
			tok = newToken(token.LBRACE , lexer.char)
		case '}' : 
			tok = newToken(token.RBRACE , lexer.char)
		case ',' : 
			tok = newToken(token.COMMA , lexer.char)
		case 0 : 
			tok.Literal = "" 
			tok.Type = token.EOF
	}

		return tok 
} 

func newToken(tokenType token.TokenType , ch byte) token.Token { 
	return token.Token{Type: tokenType , Literal: string(ch)}
}
