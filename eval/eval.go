package eval

import ( 
	"camel/ast" 
	"camel/object" 
)
var ( 
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object { 

	switch node := node.(type) { 
	
	case *ast.Program : 
		return evalStatements(node.Statements) 

	case *ast.ExpressionStatement : 
		return Eval(node.Expression) 

	case *ast.IntegerLiteral : 
		return &object.Integer{Value: node.Value} 
	
	case *ast.Boolean : 
		return evalBoolean(node.Value) 
	
	}
	return nil 
} 

func evalStatements(stmts []ast.Statement) object.Object { 
	
	var result object.Object 
	for _ , stmt := range stmts { 
		result = Eval(stmt) 
	} 
	return result
} 

func evalBoolean(input bool) *object.Boolean { 

	if input { 
		return TRUE 
	} else { 
		return FALSE 
	}
}
