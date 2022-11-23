package evaluator

import (
	"fmt"
	"regexp"
	"strings"

	parser "github.com/zain-bahsarat/rule_egine/parser"
)

// bind native funcitons
func init() {
	bindNativeFns(ListFN, list)
}

func Eval(node parser.Node, env *Environment) Object {

	switch node := node.(type) {

	case *parser.Rule:
		return Eval(node.Statement, env)

	case *parser.ExpressionStatement:
		return Eval(node.Expression, env)

	case *parser.NumberLiteral:
		return &Number{Value: node.Value}

	case *parser.StringLiteral:
		return &String{Value: node.Value}

	case *parser.BooleanLiteral:
		return &Boolean{Value: node.Value}

	case *parser.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: " + node.Value)
		}
		return val

	case *parser.Regex:
		return &Regex{Value: node.Value}

	case *parser.ListName:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("missing list: " + node.Value)
		}

		return val

	case *parser.CallExpression:
		fn, ok := nativeFns[node.Function.String()]
		if !ok {
			return newError("undefined function: " + node.Function.String())
		}

		val, err := fn(node.Arguments)
		if err != nil {
			return newError(err.Error())
		}

		return toObject(val)

	case *parser.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *parser.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	default:
		return newError("unknown: %q", node.String())
	}

}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != NumberObject {
		return nil
	}

	value := right.(*Number).Value
	return &Number{Value: -value}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == NumberObject && right.Type() == NumberObject:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == BooleanObject && right.Type() == BooleanObject:
		return evalLogicalInfixExpression(operator, left, right)
	case (left.Type() == StringObject && right.Type() == StringObject):
		return evalStringInfixExpression(operator, left, right)
	case (left.Type() == StringObject && right.Type() == RegexObject):
		return evalRegexInfixExpression(operator, left, right)
	case (left.Type() == StringObject && right.Type() == RegexListObject):
		return evalRegexListInfixExpression(operator, left, right)
	default:
		return newError("invalid expression %q %q %q ", left, operator, right)
	}
}

func evalLogicalInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Boolean).Value
	rightVal := right.(*Boolean).Value

	switch operator {

	case "and":
		return &Boolean{Value: leftVal && rightVal}
	case "or":
		return &Boolean{Value: leftVal || rightVal}
	default:
		return newError("invalid operator: %q", operator)
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Number).Value
	rightVal := right.(*Number).Value

	switch operator {

	case "<":
		return &Boolean{Value: leftVal < rightVal}
	case ">":
		return &Boolean{Value: leftVal > rightVal}
	case "==":
		return &Boolean{Value: leftVal == rightVal}
	case "!=":
		return &Boolean{Value: leftVal != rightVal}
	case "+":
		return &Number{Value: leftVal + rightVal}
	case "-":
		return &Number{Value: leftVal - rightVal}
	case "*":
		return &Number{Value: leftVal * rightVal}
	case "/":
		return &Number{Value: leftVal / rightVal}
	default:
		return newError("invalid operator: %q", operator)
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch strings.ToLower(operator) {

	case "contains":
		return &Boolean{Value: strings.Contains(leftVal, rightVal)}
	case "not_contains":
		return &Boolean{Value: !strings.Contains(leftVal, rightVal)}
	case "==":
		return &Boolean{Value: leftVal == rightVal}
	case "!=":
		return &Boolean{Value: leftVal != rightVal}
	default:
		return newError("invalid operator: %q", operator)
	}
}

func evalRegexInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*Regex).Value

	re, err := regexp.Compile(rightVal)
	if err != nil {
		return newError("invalid regex: %q", rightVal)
	}

	switch strings.ToLower(operator) {

	case "contains":
		return &Boolean{Value: len(re.Find([]byte(leftVal))) > 0}
	case "not_contains":
		return &Boolean{Value: len(re.Find([]byte(leftVal))) == 0}
	default:
		return newError("invalid operator: %q", operator)
	}
}

func evalRegexListInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*RegexList).Value

	switch strings.ToLower(operator) {

	case "contains":
		for _, re := range rightVal {
			if len(re.Find([]byte(leftVal))) > 0 {
				return &Boolean{Value: true}
			}
		}

		return &Boolean{Value: false}
	case "not_contains":
		for _, re := range rightVal {
			if len(re.Find([]byte(leftVal))) == 0 {
				return &Boolean{Value: true}
			}
		}

		return &Boolean{Value: false}
	default:
		return newError("invalid operator: %q", operator)
	}
}
