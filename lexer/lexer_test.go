package lexer 

import ( 
	"testing"
	"xxx/token"
	"xxx/lexer"
) 


func TestNextToken(t *testing.T) { 

	input := "+=();,"
	tests := []struct { 
		expectedType token.TokenType 
		expectedLiteral string
	}{ 
		{token.PLUS , "+"} , 
		{token.ASSIGN , "="} , 
		{token.LPAREN , "("} , 
		{token.RPAREN , ")" } , 
		{token.SEMICOLON , ";"},
		{token.COMMA , ","}, }

	lex := lexer.New(input) 

	for idx , testCase := range(tests) { 
	
		tok := lex.NextToken()
		
		if testCase.expectedType != tok.Type  { 
			t.Fatalf("tests[%d] - wrong token type expected [%q] : got [%q]" , 
			idx , testCase.expectedType , tok.Type) 
		}  

		if testCase.expectedLiteral != tok.Literal { 
			t.Fatalf("test [%d] faild - wrong token literal expected [%q] : got [%q]" , 
			idx , testCase.expectedLiteral , tok.Literal)  
		} 

	} 
} 
