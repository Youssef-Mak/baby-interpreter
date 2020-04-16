package evaluator

import (
	"github.com/Youssef-Mak/baby-interpreter/pkg/ast"
	"github.com/Youssef-Mak/baby-interpreter/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return boolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfixExpression(node.Operator, right, left)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalNotOperatorExpression(right)
	case "-":
		return evalNegativeOperatorExpression(right)
	default:
		return NULL
	}
}

func evalInfixExpression(operator string, right object.Object, left object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalInfixIntegerExpression(operator, right, left)
	case operator == "==":
		return boolToBooleanObject(left == right)
	case operator == "!=":
		return boolToBooleanObject(left != right)
	default:
		return NULL
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
		return NULL

	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
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
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func boolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}
