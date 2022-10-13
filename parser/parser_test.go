package parser


import ( 
	"testing"
	"xxx/lexer" 
	"xxx/ast" 
) 


func TestLetStatement(t *testing.T) { 



	input := ` 
let x = 5;
let y = 10; 
let foobar = 23421;`

	l := lexer.New(input) 
	p := New(l) 

	program := p.ParseProgram() 

	if program == nil { 
		t.Fatalf("Program returned nil") 
	} else if len(program.Statements) != 3 { 
		t.Fatalf("Program.Statements doesn't contain 3 elements, got %d instead of 3" , 
		len(program.Statements))
	}
	tests := []struct { 
		expectedIdentifier string
	}{ 

		{"x"} , 	
		{"y"} , 
		{"foobar"} , 

	}
	
	for i , tt := range(tests) { 
		stmt := program.Statements[i] 
		if ! testLetStatement(t , stmt , tt.expectedIdentifier) { 
			return 
		}	
	} 
}

func TestReturnStatement(t *testing.T) { 

	input :=`
return 4; 
return 10; 
return 124124; 
`
	lex := lexer.New(input) 
	p := New(lex) 
	
	program := p.ParseProgram() 
	
	if len(program.Statements) != 3 { 
		t.Fatalf("program.Statements doesn't contain 3 elements, got %d", 
		len(program.Statements)) 
	} 
	for _ , stmt := range(program.Statements) { 
		retStmt , ok := stmt.(*ast.ReturnStatement) 
		if !ok { 
			t.Errorf("stmt not *ast.ReturnStatement, got %T instead", 
			stmt)
			continue 
		} 

		if retStmt.TokenLiteral() != "return" { 
			t.Errorf("token literal wasn't return, got %q instead" , 
			retStmt.TokenLiteral()) 
		} 

	} 
} 

func testLetStatement(t *testing.T , stmt ast.Statement , name string) bool { 

	if stmt.TokenLiteral() != "let" { 
		t.Errorf("TokenLiteral returned wrong value, expected %s , got %s", 
		"let" , stmt.TokenLiteral()) 
		return false 
	} 
	
	letStmt , ok := stmt.(*ast.LetStatement) 

	if !ok { 
		t.Errorf("stmt is not *ast.LetStatement, got %T" , stmt) 
		return false 
	} 	
	
	if letStmt.Name.Value != name { 
		t.Errorf("letStmt.Name.Value not %s, got %s" , name , letStmt.Name.Value)
		return false 
	}
	 
	if letStmt.Name.TokenLiteral() != name { 
		t.Errorf("stmt.Name is not %s, got %s" , name , letStmt.Name.TokenLiteral())
		return false
	}
	return true 
}



