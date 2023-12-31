package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/AzraelSec/cube/pkg/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(1, o...); err != nil {
				return err
			}
			switch arg := o[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(1, o...); err != nil {
				return err
			}

			switch arg := o[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			case *object.String:
				if len(arg.Value) == 0 {
					return &object.String{Value: ""}
				}
				return &object.String{Value: arg.Value[0:1]}
			default:
				return newError("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(1, o...); err != nil {
				return err
			}
			switch arg := o[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			case *object.String:
				if len(arg.Value) == 0 {
					return &object.String{Value: ""}
				}
				return &object.String{Value: arg.Value[len(arg.Value)-1:]}
			default:
				return newError("argument to `last` not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(1, o...); err != nil {
				return err
			}
			switch arg := o[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return &object.Array{Elements: []object.Object{}}
				}
				return &object.Array{Elements: arg.Elements[1:]}
			default:
				return newError("argument to `rest` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(2, o...); err != nil {
				return err
			}

			switch arg := o[0].(type) {
			case *object.Array:
				arr := make([]object.Object, len(arg.Elements)+1, len(arg.Elements)+1)
				copy(arr, arg.Elements)
				arr[len(arr)-1] = o[1]
				return &object.Array{Elements: arr}
			default:
				return newError("argument to `push` not supported, got %s", arg.Type())
			}
		},
	},
	"print": {
		Fn: func(o ...object.Object) object.Object {
			for _, arg := range o {
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
			return NULL
		},
	},
	"read": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(0); err != nil {
				return err
			}

			reader := bufio.NewReader(os.Stdin)
			str, err := reader.ReadString('\n')
			if err != nil {
				return newError("impossible to read from stdin")
			}
			return &object.String{Value: str[:len(str)-1]}
		},
	},
	"int": {
		Fn: func(o ...object.Object) object.Object {
			if err := checkBuiltinsLenParams(1, o...); err != nil {
				return err
			}

			switch arg := o[0].(type) {
			case *object.String:
				res, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("value %s cannot be converted to int", arg.Value)
				}
				return &object.Integer{Value: res}
			case *object.Integer:
				return arg
			case *object.Boolean:
				if arg.Value == true {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}
			default:
				return newError("argument to `int` not supported, got %s", arg.Type())
			}
		},
	},
}

func checkBuiltinsLenParams(expected int, o ...object.Object) *object.Error {
	if len(o) != expected {
		return newError("wrong number of arguments. got=%d, want=1", len(o))
	}
	return nil
}
