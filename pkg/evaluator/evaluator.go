package evaluator

import (
	"fmt"
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"github.com/Youssef-Mak/baby-interpreter/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

var builtinMap = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"head": &object.BuiltIn{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.String{Value: string(arg.Value[0])}
			case *object.Array:
				return arg.Elements[0]
			default:
				return newError("argument to `head` not supported, got %s", args[0].Type())
			}
		},
	},
	"tail": &object.BuiltIn{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.String{Value: string(arg.Value[len(arg.Value)-1])}
			case *object.Array:
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `tail` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": &object.BuiltIn{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: string(arg.Value[1:len(arg.Value)])}
				}
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElems := make([]object.Object, length-1, length-1)
					copy(newElems, arg.Elements[1:len(arg.Elements)])
					return &object.Array{Elements: newElems}
				}
			default:
				return newError("argument to `tail` not supported, got %s", args[0].Type())
			}
			return NULL
		},
	},
	"push": &object.BuiltIn{
		Func: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
					len(args), 1)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElems := make([]object.Object, length+1, length+1)
					copy(newElems, arg.Elements)
					newElems[length] = args[1]
					return &object.Array{Elements: newElems}
				}
			default:
				return newError("argument to `tail` not supported, got %s", args[0].Type())
			}
			return NULL
		},
	},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return err
		}
		return evalFunctionCall(function, args)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalInfixExpression(node.Operator, right, left)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return left
		}
		return evalIndexExpression(left, index)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements, error := evalExpressions(node.Elements, env)
		if error != nil {
			return error
		}
		return &object.Array{Elements: elements}
	case *ast.Boolean:
		return boolToBooleanObject(node.Value)
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {

		result = Eval(statement, env)

		// Catch Errors
		if isError(result) {
			return result
		}

		// If statement is return, no need to keep evaluating next statements
		if returnVal, ok := result.(*object.ReturnValue); ok {
			return returnVal.Value
		}
	}

	return result
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	var result []object.Object
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{}, evaluated
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {

		result = Eval(statement, env)

		// Catch Errors
		if isError(result) {
			return result
		}

		/*
			If statement is return, no need to keep evaluating next statements
			Return must be wrapped in order to be caught by parseProgram
		*/
		if result != nil && result.Type() == object.RETURN_VAL_OBJ {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperatorExpression(right)
	case "-":
		return evalNegativeOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, right object.Object, left object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(operator, right, left)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalInfixStringExpression(operator, right, left)
	case operator == "==":
		return boolToBooleanObject(left == right)
	case operator == "!=":
		return boolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInfixStringExpression(operator string, right object.Object, left object.Object) object.Object {
	rightVal := right.(*object.String).Value
	leftVal := left.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return boolToBooleanObject(leftVal == rightVal)
	case "!=":
		return boolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalInfixIntegerExpression(operator string, right object.Object, left object.Object) object.Object {
	rightVal := right.(*object.Integer).Value
	leftVal := left.(*object.Integer).Value

	switch operator {
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "==":
		return boolToBooleanObject(leftVal == rightVal)
	case ">":
		return boolToBooleanObject(leftVal > rightVal)
	case "<":
		return boolToBooleanObject(leftVal < rightVal)
	case "<=":
		return boolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return boolToBooleanObject(leftVal >= rightVal)
	case "!=":
		return boolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalFunctionCall(funcCalled object.Object, args []object.Object) object.Object {

	switch funcCalled := funcCalled.(type) {
	case *object.Function:
		if len(args) != len(funcCalled.Parameters) {
			return newError(
				"Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
				len(funcCalled.Parameters), len(args))
		}
		funcScope := object.NewEnclosedEnvironment(funcCalled.Env)
		for idx, param := range funcCalled.Parameters {
			funcScope.Set(param.Value, args[idx])
		}
		evaluatedRes := Eval(funcCalled.Body, funcScope)
		return unwrapReturnValue(evaluatedRes)
	case *object.BuiltIn:
		return funcCalled.Func(args...)
	default:
		return newError("Is not Callable (not a recognized function): %s", funcCalled.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	array, arrOk := left.(*object.Array)
	if !arrOk {
		return newError("expecting Array Type but got %s", left.Type())
	}
	idx, idxOk := index.(*object.Integer)
	if !idxOk {
		return newError("expecting Integer Type but got %s", index.Type())
	}
	if idx.Value > int64(len(array.Elements)-1) || idx.Value < 0 {
		return NULL
	}
	return array.Elements[idx.Value]
}

func isTruthy(ob object.Object) bool {
	switch ob {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalNotOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegativeOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalFunctionLiteral(fun *ast.FunctionLiteral, env *object.Environment) object.Object {
	res := object.Function{
		Parameters: fun.Parameters,
		Body:       fun.Body,
		Env:        env,
	}
	return &res
}

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(id.Value)
	if ok {
		return val
	}

	builtin, ok := builtinMap[id.Value]
	if ok {
		return builtin
	}

	return newError("Identifier not Found: " + id.Value)
}

func boolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
