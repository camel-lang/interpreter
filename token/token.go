package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INT   = "INT"
	IDENT = "IDENT"

	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	ASSIGN   = "="
	EQ 		 = "==" 
	NOT_EQ   = "!=" 

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "fn"
	LET      = "let"

	IF     = "if"
	ELSE   = "else"
	TRUE   = "true"
	FALSE  = "false"
	RETURN = "return"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func LookUpIdent(ident string) TokenType {

	if keyword, ok := keywords[ident]; ok {
		return keyword
	}
	return IDENT
}
