package evaluator

import "github.com/AzraelSec/cube/pkg/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(o ...object.Object) object.Object {
			if len(o) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(o))
			}

			switch arg := o[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
}
