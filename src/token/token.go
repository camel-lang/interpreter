package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INT    = "INT"
	STRING = "STRING"
	IDENT  = "IDENT"

	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	ASSIGN   = "="
	EQ       = "=="
	NOT_EQ   = "!="

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

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
