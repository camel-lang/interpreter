package parser

import (
	"camel/ast"
	"camel/lexer"
	"fmt"
	"testing"
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
		t.Fatalf("Program.Statements doesn't contain 3 elements, got %d instead of 3",
			len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {

	input := `
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
	for _, stmt := range program.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T instead",
				stmt)
			continue
		}

		if retStmt.TokenLiteral() != "return" {
			t.Errorf("token literal wasn't return, got %q instead",
				retStmt.TokenLiteral())
		}

	}
}

func TestIdentifierExpression(t *testing.T) {

	input := "foobar;"

	lex := lexer.New(input)
	p := New(lex)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Expected program length 1, got %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement, got %T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier, got %T",
			stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not foobar, got %q",
			ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not foobar, got %q",
			ident.TokenLiteral())
	}
}

func testIntegerLiteralExpression(t *testing.T) {

	input := "5;"
	lex := lexer.New(input)
	p := New(lex)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("wrong length for program.Statements, got %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("wrong type for stmt, got %T, expected *ast.ExpressionStatement",
			program.Statements[0])
	}

	iliteral, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("wrong type for stmt, got %T, expected *ast.IntegerLiteral",
			program.Statements[0])
	}

	if iliteral.Value != 5 {
		t.Errorf("wrong value for integer literal, got %d, expected 5",
			iliteral.Value)
	}

}

func TestParsingPrefixExpressions(t *testing.T) {

	tests := []struct {
		Input        string
		Operator     string
		IntegerValue int64
	}{
		{"!15", "!", 15},
		{"-5", "-", 5},
	}

	for _, tt := range tests {

		lex := lexer.New(tt.Input)
		p := New(lex)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("wrong length for program.Statements, got %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("wrong type for stmt, got %T, expected *ast.ExpressionStatement",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("wrong type for stmt, got %T, expected *ast.IntegerLiteral",
				program.Statements[0])
		}

		if exp.Operator != tt.Operator {
			t.Fatalf("exp.Operator is not %s, got %s",
				tt.Operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.IntegerValue) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {

	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 / 5", 5, "/", 5},
		{"5 * 5", 5, "*", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {

		lex := lexer.New(tt.input)
		p := New(lex)

		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("wrong value for program.Statements, expected: 1, got: %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("wrong type for stmt, expected: *ast.ExpressionStatement, got: %T",
				program.Statements[0])
		}

		infexp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("wrong type for infexp, expected: *ast.InfixExpression, got: %T",
				stmt.Expression)
		}

		if !testIntegerLiteral(t, infexp.Left, tt.left) {
			return
		}
		if tt.operator != infexp.Operator {
			t.Fatalf("wrong value for exp.Operator, expected: %s, got: %s",
				tt.operator, infexp.Operator)
		}
		if !testIntegerLiteral(t, infexp.Right, tt.right) {
			return
		}

	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{

		{
		"-a * b", 
		"((-a) * b)", 
		}, 
		{
		"a + b + c", 
		"((a + b) + c)",
		},
		{
		"!-a", 
		"(!(-a))",
		},
		{
		"a - b * c + x", 
		"((a - (b * c)) + x)",
		},
		{
		"a < b == c / d", 
		"((a < b) == (c / d))",
		},
		{
		"3 - 4; 6 * 7", 
		"(3 - 4)(6 * 7)",
		},
	}

	for _, tt := range tests {

		lex := lexer.New(tt.input)
		p := New(lex)
		program := p.ParseProgram()

		if program.String() != tt.expected {
			t.Fatalf("expected: %s ,got: %s",
				tt.expected, program.String())
		}

	}

}
func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {

	if stmt.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral returned wrong value, expected %s , got %s",
			"let", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Errorf("stmt is not *ast.LetStatement, got %T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s, got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name is not %s, got %s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {

	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got %d",
			value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() returned %s, expected %d",
			integ.TokenLiteral(), value)
		return false
	}

	return true

}
