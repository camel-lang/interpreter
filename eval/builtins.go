package eval

import "camel/object" 

var builtins = map[string]*object.Builtin{
		"len" : &object.Builtin{
			Fn: func(args ...object.Object) object.Object { 
				if len(args) != 1 { 
					return newError(
				"wrong number of arguments, " + 
				"expected:1, got: %d", len(args),
					) 
				}

				switch arg := args[0].(type) { 
				
				case *object.String : 
					return &object.Integer{Value: int64(len(arg.Value))} 
				case *object.Array : 
					return &object.Integer{Value: int64(len(arg.Elements))} 
				default : 
					return newError( 
				"argument to `len` not supported, got: " + 
				"%s" , arg.Type(),
					) 
				} 
			},
		},
}



