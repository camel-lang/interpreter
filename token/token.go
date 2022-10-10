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
	IDENT = "IDENT"

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

var keywords = map[string]TokenType { 
	"fn" : FUNCTION , 
	"let" : LET ,
}

func LookUpIdent(ident string) TokenType { 

	if keyword , ok := keywords[ident] ; ok { 
		return keyword  
	} 
	return IDENT 	
} 
