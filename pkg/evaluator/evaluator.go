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
var builtinMap map[string]*object.BuiltIn

func init() {
	builtinMap = map[string]*object.BuiltIn{
		"len": { // Return length of array or string
			Func: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
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
		"head": { // Return the first element of array or string
			Func: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
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
		"tail": { // Return the last element of an array or string
			Func: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
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
		"isEmpty": { // Returns string or array omitting the first element
			Func: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
				}
				switch arg := args[0].(type) {
				case *object.String:
					if len(arg.Value) > 0 {
						return FALSE
					} else {
						return TRUE
					}
				case *object.Array:
					length := len(arg.Elements)
					if length > 0 {
						return FALSE
					} else {
						return TRUE
					}
				default:
					return newError("argument to `isEmpty` not supported, got %s", args[0].Type())
				}
			},
		},
		"rest": { // Returns string or array omitting the first element
			Func: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
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
					return newError("argument to `rest` not supported, got %s", args[0].Type())
				}
				return NULL
			},
		},
		"get": { // Returns Element of Array at passed index
			Func: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						2, len(args))
				}

				switch arg := args[0].(type) {
				case *object.Array:
					length := len(arg.Elements)
					index, idOk := args[1].(*object.Integer)
					if length > 0 {
						if !idOk {
							return newError("argument to `get` not supported, got %s", args[1].Type())
						}
						if index.Value > int64(length-1) {
							return newError("Array our bounds. Length of Array is %d, index to be accessed is %d", length, index)
						}
						return arg.Elements[index.Value]
					}
				default:
					return newError("argument to `get` not supported, got %s", args[0].Type())
				}
				return NULL
			},
		},
		"insert": { // Returns Array with element(s) inserted into the index specified
			Func: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						3, len(args))
				}

				switch arg := args[0].(type) {
				case *object.Array:
					length := len(arg.Elements)
					idx, idxOk := args[2].(*object.Integer)
					if !idxOk {
						return newError("arguments to `insert` not supported, expected Integer, got %s", args[2].Type())
					}
					if idx.Value > int64(length-1) {
						return newError("arguments to `insert` not supported, expected Index passed to be equal or less than size of list, got index %d but size of array is %d", idx.Value, length)
					}
					toIns := args[1]
					newElems := make([]object.Object, length, length)
					copy(newElems, arg.Elements)
					newElems[idx.Value] = toIns
					return &object.Array{Elements: newElems}
				default:
					return newError("argument to `insert` not supported, got %s", args[0].Type())
				}
			},
		},
		"append": { // Returns Array with element(s) inserted into the end of passed array
			Func: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						2, len(args))
				}

				switch arg := args[0].(type) {
				case *object.Array:
					length := len(arg.Elements)
					toAdd, isArr := args[1].(*object.Array)
					if isArr {
						newElems := make([]object.Object, length+len(toAdd.Elements), length+len(toAdd.Elements))
						copy(newElems, arg.Elements)
						for i, value := range toAdd.Elements {
							newElems[length+i] = value
						}
						return &object.Array{Elements: newElems}
					} else {
						newElems := make([]object.Object, length+1, length+1)
						copy(newElems, arg.Elements)
						newElems[length] = args[1]
						return &object.Array{Elements: newElems}
					}
				default:
					intHead, intOk := args[0].(*object.Integer)
					arrRest, rarrOk := args[1].(*object.Array)
					// Reverse append(Enqueue edge case)
					if intOk && rarrOk {
						length := len(arrRest.Elements)
						newElems := make([]object.Object, length+1, length+1)
						newElems[0] = intHead
						for i, e := range arrRest.Elements {
							newElems[i+1] = e
						}
						return &object.Array{Elements: newElems}
					}
					return newError("argument to `append` not supported, got %s", args[0].Type())
				}
			},
		},
		"doWhile": { // Calls function returning a boolean until call resolves to false
			Func: func(args ...object.Object) object.Object {
				if len(args) > 1 {
					return newError("Call Arguments and function defined parameters size mismatch.\n Expected %d arguments but got %d parameter(s)",
						1, len(args))
				}
				body, bodyOk := args[0].(*object.Function)
				if !bodyOk {
					return newError("arguments to `doWhile` not supported, expected Function, got %s", args[0].Type())
				}
				ret := TRUE
				for ret == TRUE {
					ret, _ = evalFunctionCall(body, nil).(*object.Boolean)
				}
				if ret != TRUE && ret != FALSE {
					return newError("arguments to `doWhile` not supported, expected Function to return Boolean, got %s", ret.Type())
				}
				return ret
			},
		},
		"print": {
			Func: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.AssignmentStatement:
		deepCopyFlag := node.AssignmentOperator.Literal != "=&"
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val, deepCopyFlag)
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
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.DotExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		attribute := Eval(node.Attribute, env)
		if isError(attribute) {
			return attribute
		}
		return evalDotExpression(left, attribute)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ArrayLiteral:
		elements, error := evalExpressions(node.Elements, env)
		if error != nil {
			return error
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		evalExprMap, error := evalMappedExpressions(node.Pairs, env)
		if error != nil {
			return error
		}
		return evalExprMap
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
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

func evalMappedExpressions(exprs map[ast.Expression]ast.Expression, env *object.Environment) (object.Object, object.Object) {
	evaldMap := &object.Hash{Pairs: map[object.HashKey]object.HashEntry{}}

	for keyExpr, valExpr := range exprs {
		keyEvaled := Eval(keyExpr, env)
		if isError(keyEvaled) {
			return nil, keyEvaled
		}
		hashKey, ok := keyEvaled.(object.Hashable)
		if !ok {
			return evaldMap, newError("This key is not Hashable : %s", keyEvaled.Inspect())
		}
		valEvaled := Eval(valExpr, env)
		if isError(valEvaled) {
			return nil, valEvaled
		}
		hEntry := object.HashEntry{Key: keyEvaled, Value: valEvaled}
		evaldMap.Pairs[hashKey.HashKey()] = hEntry
	}

	return evaldMap, nil
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
	case operator == "=&=":
		return boolToBooleanObject(left == right)
	case operator == "=*=":
		return boolToBooleanObject(left.Inspect() == right.Inspect())
	case operator == "!&=":
		return boolToBooleanObject(left != right)
	case operator == "!*=":
		return boolToBooleanObject(left.Inspect() != right.Inspect())
	case operator == "&":
		return boolToBooleanObject(left == TRUE && right == TRUE)
	case operator == "|":
		return boolToBooleanObject(left == TRUE || right == TRUE)
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
	case "=*=":
		return boolToBooleanObject(leftVal == rightVal)
	case "=&=":
		return boolToBooleanObject(left == right)
	case "!&=":
		return boolToBooleanObject(left != right)
	case "!*=":
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
	case "=*=":
		return boolToBooleanObject(leftVal == rightVal)
	case "=&=":
		return boolToBooleanObject(left == right)
	case "!&=":
		return boolToBooleanObject(left != right)
	case "!*=":
		return boolToBooleanObject(leftVal != rightVal)
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

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	condition := Eval(we.Condition, env)
	ret := object.ReturnValue{Value: NULL}
	for isTruthy(condition) {
		ret.Value = Eval(we.Body, env)
		condition = Eval(we.Condition, env)
	}
	return ret.Value
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
			funcScope.Set(param.Value, args[idx], false)
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

func evalDotExpression(left object.Object, attribute object.Object) object.Object {
	hash, ok := left.(*object.Hash)
	if !ok {
		return newError("expecting Hash Type but got %s", left.Type())
	}
	attr, attrOk := attribute.(object.Hashable)
	if !attrOk {
		return newError("expecting Hashable Type but got %s", attribute.Type())
	}

	pair, found := hash.Pairs[attr.HashKey()]
	if !found {
		return NULL
	}
	return pair.Value
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
		return *val
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
