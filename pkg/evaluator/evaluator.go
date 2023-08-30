package evaluator

import (
	"github.com/AzraelSec/cube/pkg/ast"
	"github.com/AzraelSec/cube/pkg/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBooleanMap(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node.Operator, Eval(node.Right, env))
	case *ast.InfixExpression:
		return evalInfixExpression(node.Operator, Eval(node.Left, env), Eval(node.Right, env))
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.RetValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		return evalFuncLiteral(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	}
	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}
	return val
}

func evalFuncLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}

	evalParams := make([]object.Object, len(node.Args))
	for i, param := range node.Args {
		param := Eval(param, env)
		if isError(param) {
			return param
		}
		evalParams[i] = param
	}

	return applyFunction(function, evalParams)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	if len(args) != len(function.Parameters) {
		return newError("wrong number of arguments for function %s: %d instead of %d", "IDK", len(args), len(function.Parameters))
	}
	extEnv := extendedFunctionEnv(function, args)
	evaluated := Eval(function.Body, extEnv)
	return unwrapReturnValue(evaluated)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}
	return env
}

func unwrapReturnValue(evaluated object.Object) object.Object {
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return evaluated
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}
	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return NULL
}

func evalProgram(stms []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, stm := range stms {
		res = Eval(stm, env)
		switch res := res.(type) {
		// note: early exit if we meet a return statement in top-level loop
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object
	for _, stm := range block.Statements {
		res = Eval(stm, env)

		if res != nil && res.Type() == object.RETURN_VALUE_OBJ || res.Type() == object.ERROR_OBJ {
			return res
		}
	}
	return res
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	if isError(right) {
		return right
	}

	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		switch exp := right.(type) {
		case *object.Integer:
			return nativeBooleanMap(exp.Value == 0)
		default:
			return FALSE
		}
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -1 * value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(op, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalInfixBooleanExpression(op, left, right)
	// todo: remove this case to impement concatenation
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator %s %s %s", left.Type(), op, right.Type())
	}
}

func evalInfixBooleanExpression(op string, left, right object.Object) object.Object {
	leftVal, rightVal := left.(*object.Boolean).Value, right.(*object.Boolean).Value

	switch op {
	case "==":
		return nativeBooleanMap(leftVal == rightVal)
	case "!=":
		return nativeBooleanMap(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalInfixIntegerExpression(op string, left, right object.Object) object.Object {
	leftVal, rightVal := left.(*object.Integer).Value, right.(*object.Integer).Value

	switch op {
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBooleanMap(leftVal > rightVal)
	case "<":
		return nativeBooleanMap(leftVal < rightVal)
	case "!=":
		return nativeBooleanMap(leftVal != rightVal)
	case "==":
		return nativeBooleanMap(leftVal == rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, left.Type())
	}
}