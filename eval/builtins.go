package eval

import "camel/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments, "+
						"expected:1, got: %d", len(args),
				)
			}

			switch arg := args[0].(type) {

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(
					"argument to `len` not supported, got: "+
						"%s", arg.Type(),
				)
			}
		},
	},
	"peek": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments, "+
						"expected:1, got: %d", len(args),
				)
			}

			arg, ok := args[0].(*object.Array)

			if !ok {
				return newError("argument to peek must be array,"+
					" got %s", args[0].Type())
			}

			arr := arg.Elements
			return arr[len(arr)-1]
		},
	},

	"pop": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(
					"wrong number of arguments, "+
						"expected:1, got: %d", len(args),
				)
			}

			arg, ok := args[0].(*object.Array)

			if !ok {
				return newError("argument to peek must be array,"+
					" got %s", args[0].Type())
			}

			arr := arg.Elements
			return &object.Array{Elements: arr[:len(arr)-1]}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(
					"wrong number of arguments, "+
						"expected:2, got: %d", len(args),
				)
			}

			arg, ok := args[0].(*object.Array)

			if !ok {
				return newError("argument to peek must be array,"+
					" got %s", args[0].Type())
			}

			elem := args[1]
			Elements := arg.Elements
			length := len(Elements)
			newElements := make(
				[]object.Object,
				length+1,
				length+1)

			copy(newElements, Elements)
			newElements[length] = elem
			return &object.Array{Elements: newElements}
		},
	},
}
