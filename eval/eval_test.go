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
	env := object.NewEnvironment()
	return Eval(program, env)
}

func TestEvalBooleanExpression(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expectedValue)
	}
}

func TestBangOperator(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"!true", false},
		{"!!false", false},
		{"!!5", true},
		{"!3", false},
		{"!0", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expectedValue)
	}

}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}
	for _, tt := range tests {
		expectedValue := testEval(tt.input)
		testIntegerObject(t, expectedValue, tt.expected)
	}
}
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Type mismatch: invalid operator + for types INTEGER BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Type mismatch: invalid operator + for types INTEGER BOOLEAN",
		},
		{
			"-true",
			"Invalid operator: type BOOLEAN doesn't support '-' operator",
		},
		{
			"true + false;",
			"Unknown operator: no + operator registered for BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Unknown operator: no + operator registered for BOOLEAN",
		},
		{
			"5; true + true - false  + false; 5",
			"Unknown operator: no + operator registered for BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Unknown operator: no + operator registered for BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"Unknown operator: no + operator registered for BOOLEAN",
		},
		{
			"foobar",
			"Identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
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

func testNullObject(t *testing.T, obj object.Object) bool {

	if obj != NULL {
		t.Errorf("obj is not NULL, got: %T (+%v)", obj, obj)
		return false
	}
	return true
}
