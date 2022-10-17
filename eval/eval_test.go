package eval

import (
	"camel/lexer"
	"camel/object"
	"camel/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"5", 5},
		{"1", 1},
		{"-1", -1},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {

		output := testEval(tt.input)
		testIntegerObject(t, output, tt.expectedValue)
	}
}

func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()
	return Eval(program)
}

func TestEvalBooleanExpression(t *testing.T) { 

	tests := []struct{ 
		input string 
		expectedValue bool
	}{ 
		{"true", true}, 
		{"false", false},
	}

	for _ , tt := range tests { 
		evaluated := testEval(tt.input) 
		testBooleanObject(t, evaluated, tt.expectedValue) 
	} 
}

func TestBangOperator(t *testing.T) { 
	
	tests := []struct { 
		input string 
		expectedValue bool
	}{ 
		{"!true", false}, 
		{"!!false", false}, 
		{"!!5", true}, 
		{"!3", false},
		{"!0", true},
	}
	
	for _ , tt := range tests { 
		evaluated := testEval(tt.input) 
		testBooleanObject(t, evaluated, tt.expectedValue) 
	}

}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("obj is not Integer, got: %T, (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("wrong value for result.Value, expected: %v, got: %v",
			expected, result.Value)
		return false
	}
	return true 
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool { 
	
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("obj is not Boolean, got: %T, (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("wrong value for result.Value, expected: %v, got: %v",
			expected, result.Value)
		return false
	}
	return true 

} 
