package eval

import ( 
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
		return evalStatements(node.Statements) 
	
	case *ast.BlockStatement : 
		return evalStatements(node.Statements) 

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
		return nil
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
func evalMinusPrefixOperator(obj object.Object) object.Object { 
	if obj.Type() != object.INTEGER_OBJ { 
		return NULL
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

	default : 
		return NULL
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
		return NULL
	}
}

func parseBooleanInfixExpression(operator string, left, right object.Object) object.Object { 
	
	switch operator {
		
	case "==" : 
		return nativeBoolean(left == right) 
	case "!=" :
		return nativeBoolean(left != right) 
	default : 
		return NULL

	}

} 
