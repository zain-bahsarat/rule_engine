package evaluator

import (
	"testing"

	"github.com/zain-bahsarat/rule_egine/parser"
)

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"-10", -10},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input, make(map[string]interface{}))
		testNumberObject(t, evaluated, tt.expected)
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testEval(t *testing.T, input string, bindings map[string]interface{}) Object {
	l := parser.NewLexer(input)
	p := parser.New(l)
	rule := p.ParseRule()
	// fmt.Println(rule.String())
	checkParserErrors(t, p)
	env := NewEnvironment(bindings)

	return Eval(rule, env)
}

func testNumberObject(t *testing.T, obj Object, expected float64) bool {
	result, ok := obj.(*Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{`"a" == "a"`, true},
		{`"a" == "b"`, false},
		{`"a" != "a"`, false},
		{`"a" != "b"`, true},
		{`"a" contains "a"`, true},
		{`"a" not_contains "a"`, false},
		{`"a" not_contains "b"`, true},
		// {`test == 7`, true},
		// {`test1 == 7.5`, true},
		// {`test == 0`, false},
		{`"a" contains r"a*"`, true},
		{`"abc" contains r"ad.*"`, false},
		{`@list contains regex`, false},
		{`regex contains list("a.*", "d")`, true},
		{`a > b and a > b and a > b and a > b and a > b and a > b and a > b`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input, map[string]interface{}{"a": 8, "b": 7.5, "regex": "abc", "list": []string{"abd", "ad"}})
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input, make(map[string]interface{}))
		n := evaluated.(*Number)
		if n.Value != tt.expected {
			t.Errorf("object has wrong value. got=%f, want=%f",
				n.Value, tt.expected)
		}
	}
}
