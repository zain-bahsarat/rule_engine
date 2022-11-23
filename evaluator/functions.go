package evaluator

import (
	"strings"

	"github.com/zain-bahsarat/rule_egine/parser"
)

const (
	ListFN = "LIST"
)

type nativeFn func(args interface{}) (interface{}, error)

var (
	nativeFns = make(map[string]nativeFn)
)

type Stringable interface {
	String() string
}

func bindNativeFns(name string, fn nativeFn) {
	name = strings.ToLower(name)
	nativeFns[name] = fn
}

func list(args interface{}) (interface{}, error) {
	list := make([]string, 0)
	for _, arg := range args.([]parser.Expression) {
		list = append(list, arg.TokenLiteral())
	}

	return NewRegexList(list), nil
}
