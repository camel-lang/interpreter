package eval

import (
	"camel/ast"
	"camel/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.FunctionLiteral:

		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:

		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolean(node.Value)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}

	}
	return nil
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {

	var objs []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		objs = append(objs, evaluated)
	}

	return objs
}

func applyFunction(
	f object.Object,
	args []object.Object,
) object.Object {

	switch fn := f.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("Invalid function call, %s is not a function", f.Type())

	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}
	return obj
}

func extendFunctionEnv(
	f *object.Function,
	args []object.Object,
) *object.Environment {

	extendedEnv := object.NewEnclosedEnvironment(f.Env)
	for id, p := range f.Parameters {
		extendedEnv.Set(p.Value, args[id])
	}

	return extendedEnv
}

func evalIndexExpression(
	left object.Object,
	index object.Object,
) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ &&
		index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("Invalid Index: index operator not "+
			"supported for type %s", left.Type())
	}
}

func evalArrayIndexExpression(
	array object.Object,
	index object.Object,
) object.Object {

	arrayObj := array.(*object.Array)
	id := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements))

	if id < 0 || id > max {
		return newError("Index out of range")
	}

	return arrayObj.Elements[id]
}

func evalHashIndexExpression(
	hash object.Object,
	index object.Object,
) object.Object {

	hashObj, ok := hash.(*object.Hash)

	hashKey, ok := index.(object.Hashable)
	if !ok {
		return newError("Unhashable type %s "+
			"used as index", index.Type())
	}

	p, ok := hashObj.Pairs[hashKey.HashKey()]
	if !ok {
		return NULL
	}

	return p.Value
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {

	pairs := make(map[object.HashKey]object.HashPair)

	for k, v := range node.Pairs {

		key := Eval(k, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("Object %s not hashable", key.Type())
		}

		value := Eval(v, env)
		if isError(value) {
			return value
		}

		pairs[hashKey.HashKey()] = object.HashPair{Key: key,
			Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	fmt.Println(env)
	return newError("Identifier not found: %s", node.Value)

}
func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {

	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTrue(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTrue(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}
func nativeBoolean(input bool) *object.Boolean {

	if input {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(
	operator string,
	obj object.Object,
) object.Object {

	switch operator {
	case "!":
		return evalBangOperatorExpression(obj)
	case "-":
		return evalMinusPrefixOperator(obj)
	default:
		return newError("Unknown operator: operator %s is not a valid prefix operator", operator)
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object {

	switch obj := obj.(type) {
	case *object.Integer:
		return evalBangInteger(obj)
	case *object.Boolean:
		return evalBangBoolean(obj)
	default:
		return FALSE
	}
}
func evalBangInteger(num *object.Integer) *object.Boolean {
	if num.Value == 0 {
		return TRUE
	} else {
		return FALSE
	}
}
func evalBangBoolean(boolean *object.Boolean) *object.Boolean {
	if boolean.Value {
		return FALSE
	} else {
		return TRUE
	}

}
func evalMinusPrefixOperator(
	obj object.Object,
) object.Object {
	if obj.Type() != object.INTEGER_OBJ {
		return newError("Invalid operator: type %s doesn't support '-' operator", obj.Type())
	}

	val := obj.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	switch {

	case left.Type() == object.INTEGER_OBJ &&
		right.Type() == object.INTEGER_OBJ:
		return parseIntegerInfixExpression(operator, left, right)

	case left.Type() == object.BOOLEAN_OBJ &&
		right.Type() == object.BOOLEAN_OBJ:
		return parseBooleanInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ &&
		right.Type() == object.STRING_OBJ:
		return parseStringInfixExpression(operator, left, right)

	case left.Type() != right.Type():
		return newError("Type mismatch: invalid operator %s for types %s %s",
			operator, left.Type(), right.Type())
	default:
		return newError("Unknown operator: no %s operator registered for %s", operator, left.Type())
	}
}

func parseStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {

	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("Unknown operator: no %s operator registered for Strings", operator)
	}

}

func parseIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {

	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "<":
		return nativeBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolean(leftVal > rightVal)
	case "==":
		return nativeBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolean(leftVal != rightVal)
	default:
		return newError("Unknown operator: no %s operator registered for Integers", operator)
	}
}

func parseBooleanInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {

	switch operator {

	case "==":
		return nativeBoolean(left == right)
	case "!=":
		return nativeBoolean(left != right)
	default:
		return newError("Unknown operator: no %s operator registered for BOOLEAN", operator)

	}
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {

	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ERROR_OBJ ||
				rt == object.RETURN_VALUE_OBJ {
				return result
			}
		}
	}
	return result
}

func evalProgram(
	program *ast.Program,
	env *object.Environment,
) object.Object {

	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {

		case *object.Error:
			return result
		case *object.ReturnValue:
			return result.Value
		}
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {

	if obj.Type() == object.ERROR_OBJ {
		return true
	} else {
		return false
	}
}
