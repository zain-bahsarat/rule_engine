package evaluator

import (
	"errors"
	"strings"

	"github.com/zain-bahsarat/rule_egine/parser"
)

type Rule struct {
	expression string
	parsedRule *parser.Rule
	metadata   map[string]interface{}
}

func NewRule(expression string, metadata map[string]interface{}) (*Rule, error) {

	p := parser.New(parser.NewLexer(expression))
	parsedRule := p.ParseRule()
	if len(p.Errors()) > 0 {
		return nil, errors.New(strings.Join(p.Errors(), "\n"))
	}

	return &Rule{
		expression: expression,
		parsedRule: parsedRule,
		metadata:   metadata,
	}, nil
}

func (r *Rule) Eval(params map[string]interface{}) bool {
	result := Eval(r.parsedRule, NewEnvironment(params))

	res, ok := result.(*Boolean)
	if !ok {
		return false
	}

	return res.Value == true
}

func (r *Rule) Expression() string {
	return r.expression
}

func (r *Rule) AddMetadata(key string, value interface{}) {
	r.metadata[key] = value
}

func (r *Rule) GetMetadata(key string) interface{} {
	return r.metadata[key]
}
