package evaluator

import (
	"fmt"
	"regexp"
)

type ObjectType string

const (
	NumberObject     = "Number"
	BooleanObject    = "Boolean"
	StringObject     = "String"
	IdentifierObject = "Identifier"
	ErrorObject      = "Error"
	RegexObject      = "Regex"
	RegexListObject  = "RegexList"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType {
	return NumberObject
}

func (n *Number) Inspect() string {
	return fmt.Sprintf("%f", n.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BooleanObject
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type String struct {
	Value string
}

func (b *String) Type() ObjectType {
	return StringObject
}

func (b *String) Inspect() string {
	return fmt.Sprintf("%s", b.Value)
}

type Regex struct {
	Value string
}

func (b *Regex) Type() ObjectType {
	return RegexObject
}

func (b *Regex) Inspect() string {
	return fmt.Sprintf("%s", b.Value)
}

type RegexList struct {
	Value []*regexp.Regexp
}

func (b *RegexList) Type() ObjectType {
	return RegexListObject
}

func (b *RegexList) Inspect() string {
	return fmt.Sprintf("%q", b.Value)
}

func NewRegexList(values []string) *RegexList {
	rs := make([]*regexp.Regexp, 0)

	for _, v := range values {
		r, err := regexp.Compile(v)
		if err != nil {
			continue
		}

		rs = append(rs, r)
	}

	return &RegexList{Value: rs}
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObject }
func (e *Error) Inspect() string  { return "error: " + e.Message }

func toObject(val interface{}) Object {
	switch val := val.(type) {
	case *RegexList:
		return val
	default:
		return nil
	}
}
