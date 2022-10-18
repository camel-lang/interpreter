package eval

import ( 
	"fmt"
	"camel/ast" 
	"camel/object" 
)
var ( 
	NULL = &object.Null{} 
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object { 

	switch node := node.(type) { 
	
	case *ast.Program : 
		return evalProgram(node) 
	
	case *ast.BlockStatement : 
		return evalBlockStatement(node)  

	case *ast.ExpressionStatement : 
		return Eval(node.Expression) 

	case *ast.PrefixExpression : 
		right := Eval(node.Right) 
		return evalPrefixExpression(node.Operator, right) 
	
	case *ast.InfixExpression : 
		left := Eval(node.Left) 
		right := Eval(node.Right) 
		return evalInfixExpression(node.Operator, left, right) 

	case *ast.IfExpression : 
		return evalIfExpression(node) 

	case *ast.IntegerLiteral : 
		return &object.Integer{Value: node.Value} 
	
	case *ast.Boolean : 
		return nativeBoolean(node.Value) 

	case *ast.ReturnStatement : 
		val := Eval(node.ReturnValue) 
		return &object.ReturnValue{Value: val}
	
	}
	return nil 
} 

func evalIfExpression(ie *ast.IfExpression) object.Object { 

	condition := Eval(ie.Condition)

	if isTrue(condition) { 
		return Eval(ie.Consequence) 
	} else if ie.Alternative != nil{ 
		return Eval(ie.Alternative)
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
func evalStatements(stmts []ast.Statement) object.Object { 
	
	var result object.Object 
	for _ , stmt := range stmts { 
		result = Eval(stmt) 
		
		if returnValue, ok := result.(*object.ReturnValue); ok { 
			return returnValue.Value
		}
	} 
	return result
} 

func nativeBoolean(input bool) *object.Boolean { 

	if input { 
		return TRUE 
	} else { 
		return FALSE 
	}
}

func evalPrefixExpression(operator string, obj object.Object) object.Object { 

	switch operator { 
	case "!": 
		return evalBangOperatorExpression(obj) 
	case "-": 
		return evalMinusPrefixOperator(obj) 
	default : 
		return newError("Unknown operator: operator %s is not a valid prefix operator", operator) 
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object { 
	
	switch obj := obj.(type) { 
		case *object.Integer: 
			return evalBangInteger(obj)  
		case *object.Boolean: 
			return evalBangBoolean(obj)  
		default : 
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
		 right.Type() == object.INTEGER_OBJ :
		return parseIntegerInfixExpression(operator, left, right) 
	
	case left.Type() == object.BOOLEAN_OBJ && 
		 right.Type() == object.BOOLEAN_OBJ : 
		return parseBooleanInfixExpression(operator, left, right) 

	case left.Type() != right.Type() : 
		return newError("Type mismatch: invalid operator %s for types %s %s", 
		operator, left.Type(), right.Type()) 
	default : 
		return newError("Unknown operator: no %s operator registered for %s", operator, left.Type())
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
	default : 
		return newError("Unknown operator: no %s operator registered for Integers", operator)
	}
}

func parseBooleanInfixExpression(
	operator string, 
	left, right object.Object, 
) object.Object { 
	
	switch operator {
		
	case "==" : 
		return nativeBoolean(left == right) 
	case "!=" :
		return nativeBoolean(left != right) 
	default : 
		return newError("Unknown operator: no %s operator registered for BOOLEAN", operator)

	}
}

func evalBlockStatement(block *ast.BlockStatement) object.Object { 

	var result object.Object 
	
	for _, statement := range block.Statements { 
		result = Eval(statement) 
	
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

func evalProgram(program *ast.Program) object.Object { 

	var result object.Object 
	
	for _, statement := range program.Statements { 
		result = Eval(statement) 
	
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

	if obj.Type == object.ERROR_OBJ {
		return true 
	} else { 
		return false 
	} 
}  
