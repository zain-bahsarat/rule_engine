package parser

import (
	"bytes"
	"fmt"
	"strings"
)

// precedence
const (
	_ int = iota
	LOWEST
	LOGICAL     // AND OR
	EQ          // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	DIVIDE      // /
	PREFIX      // -
	CALL        // myFunction(X)
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
	expressionNode()
}

// Rule - as a root node
type Rule struct {
	Statement Statement
}

func (r *Rule) TokenLiteral() string {
	if r.Statement != nil {
		return ""
	}

	return r.Statement.TokenLiteral()
}

func (r *Rule) String() string {
	var out bytes.Buffer
	out.WriteString(r.Statement.String())
	return out.String()
}

type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ListName struct {
	Token Token
	Value string
}

func (l *ListName) expressionNode()      {}
func (l *ListName) TokenLiteral() string { return l.Token.Literal }
func (l *ListName) String() string       { return l.Value }

type Regex struct {
	Token Token
	Value string
}

func (r *Regex) expressionNode()      {}
func (r *Regex) TokenLiteral() string { return r.Token.Literal }
func (r *Regex) String() string       { return fmt.Sprintf("\"%s\"", r.Token.Literal) }

type StringLiteral struct {
	Token Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return fmt.Sprintf("\"%s\"", s.Token.Literal) }

type NumberLiteral struct {
	Token Token
	Value float64
}

func (il *NumberLiteral) expressionNode()      {}
func (il *NumberLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *NumberLiteral) String() string       { return il.Token.Literal }

type BooleanLiteral struct {
	Token Token
	Value bool
}

func (il *BooleanLiteral) expressionNode()      {}
func (il *BooleanLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *BooleanLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    Token
	Operator string
	Left     Expression
	Right    Expression
}

func (in *InfixExpression) expressionNode()      {}
func (in *InfixExpression) TokenLiteral() string { return in.Token.Literal }
func (in *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(in.Left.String())
	out.WriteString(" " + in.Operator + " ")
	out.WriteString(in.Right.String())
	out.WriteString(")")

	return out.String()
}

// ExpressionStatement contains the expression
type ExpressionStatement struct {
	Token      Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type CallExpression struct {
	Token     Token      // The '(' token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
