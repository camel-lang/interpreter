package lexer

import (
	"testing"
	"xxx/token"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5 ; 
let ten = 10 ; 
let add = fn(x , y) { 
x + y ; 
}; 
let result = add(five , ten);

< > 5! -* if else true false`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.EOF, ""},
	}

	lex := New(input)

	for idx, testCase := range tests {

		tok := lex.NextToken()

		if testCase.expectedType != tok.Type {
			t.Fatalf("tests[%d] - wrong token type expected [%q] : got [%q]",
				idx, testCase.expectedType, tok.Type)
		}

		if testCase.expectedLiteral != tok.Literal {
			t.Fatalf("test [%d] faild - wrong token literal expected [%q] : got [%q]",
				idx, testCase.expectedLiteral, tok.Literal)
		}

	}
}
