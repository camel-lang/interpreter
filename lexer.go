package lexer

import "token"


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
	
	if lexer.readPosition >= lexer.len(input) { 
		lexer.char = 0 
	} else { 
		lexer.char = lexer.input[lexer.readPosition]
	} 
	lexer.position = lexer.readPosition 
	lexer.readPosition += 1 
}	

func (lexer *Lexer) NextToken() token.Token {

	var tok token.Token 
	lexer.readChar() 

	switch lexer.char { 


		case "=" : 
			tok = newToken(token.ASSIGN , "=")
		case "+" : 
			tok = newToken(token.PLUS , "+") 
		case ";" : 
			tok = newToken(token.SEMICOLON , ";") 
		case "(" : 
			tok = newToken(token.LPARAN , "(") 
		case ")" :
			tok = newToken(token.RPARAN , ")") 
		case "{" : 
			tok = newToken(token.LBRACE , "{") 
		case "}" : 
			tok = newToken(token.RBRACE , "}") 
		case "," : 
			tok = newToken(token.COLON , ",") 
		case 0 : 
			tok = newToken(token.EOF , "") 	

		return tok 

	}
} 

func newToken(tokenType token.TokenType , ch byte) token.Token { 
	return token.Token{TokenType: tokenType , Literal: string(ch)}
}

