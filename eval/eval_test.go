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

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("obj is not Integer, got: %T, (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("wrong value for result.Value, expected: %d, got: %d",
			expected, result.Value)
		return false
	}
	return true 
}
