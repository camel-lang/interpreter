package token 


type TokenType string 


type Token struct{ 
	Type TokenType 
	Literal string 
}

const ( 

	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	INT = "INT"
	INDENT = "INDENT"

	PLUS = "+" 
	ASSIGN = "="
	
	COMMA = "," 
	SEMICOLON = ";" 
	
	LPAREN = "(" 
	RPAREN = ")" 
	LBRACE = "{" 
	RBRACE = "}"

	FUNCTION = "FUNCTION" 
	LET = "LET" 

)
