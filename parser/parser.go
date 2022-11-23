package parser

import (
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

var precendences = map[TokenType]int{
	EQUALS:      EQ,
	NOTEQUAL:    EQ,
	CONTAINS:    EQ,
	NOTCONTAINS: EQ,
	LT:          LESSGREATER,
	LTE:         LESSGREATER,
	GT:          LESSGREATER,
	GTE:         LESSGREATER,
	PLUS:        SUM,
	MINUS:       SUM,
	ASTERIK:     PRODUCT,
	FSLASH:      DIVIDE,
	LPAREN:      CALL,
	AND:         LOGICAL,
	OR:          LOGICAL,
}

type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func New(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(NUMBER, p.parseNumberLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(REGEX, p.parseRegex)
	p.registerPrefix(LISTNAME, p.parseList)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(FSLASH, p.parseInfixExpression)
	p.registerInfix(ASTERIK, p.parseInfixExpression)
	p.registerInfix(EQUALS, p.parseInfixExpression)
	p.registerInfix(NOTEQUAL, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(LTE, p.parseInfixExpression)
	p.registerInfix(GTE, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(CONTAINS, p.parseInfixExpression)
	p.registerInfix(NOTCONTAINS, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expect(t TokenType) bool {
	if !p.curTokenIs(t) {
		p.peekError(t)
		return false
	}

	return true
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecendence() int {
	if pr, ok := precendences[p.peekToken.Type]; ok {
		return pr
	}

	return LOWEST
}

func (p *Parser) curPrecendence() int {
	if pr, ok := precendences[p.curToken.Type]; ok {
		return pr
	}

	return LOWEST
}

func (p *Parser) ParseRule() *Rule {
	rule := &Rule{Statement: &ExpressionStatement{}}

	for !p.curTokenIs(EOF) {
		rule.Statement = p.parseExpressionStatement()
		p.nextToken()
	}

	return rule
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))

	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.curTokenIs(EOF) && precedence < p.peekPrecendence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	defer untrace(trace("parseIdentifier"))

	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() Expression {
	defer untrace(trace("parseBooleanLIteral"))

	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TRUE)}
}

func (p *Parser) parseStringLiteral() Expression {
	defer untrace(trace("parseStringLiteral"))

	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseList() Expression {
	defer untrace(trace("parseList"))

	return &ListName{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseRegex() Expression {
	defer untrace(trace("parseRegex"))

	return &Regex{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() Expression {
	defer untrace(trace("parseGroupedExpression"))

	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseNumberLiteral() Expression {
	defer untrace(trace("parseNumberLiteral"))

	lit := &NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() Expression {
	defer untrace(trace("parsePrefixExpression"))

	exp := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(leftExp Expression) Expression {
	defer untrace(trace("parseInfixExpression"))

	exp := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     leftExp,
	}

	precendence := p.curPrecendence()
	p.nextToken()
	exp.Right = p.parseExpression(precendence)

	return exp
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}
	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) traverseNode(node Node, m map[string]interface{}) {
	switch node := node.(type) {
	case *Identifier:
		m["identifiers"] = append(m["identifiers"].([]string), node.Value)

	case *ListName:
		m["lists"] = append(m["lists"].([]string), node.Value)

	case *Regex:
		m["regexs"] = append(m["regexs"].([]string), node.Value)

	case *NumberLiteral:
		m["numbers"] = append(m["numbers"].([]float64), node.Value)

	case *StringLiteral:
		m["strings"] = append(m["strings"].([]string), node.Value)

	case *CallExpression:

		args := []string{}
		for _, a := range node.Arguments {
			args = append(args, a.String())
		}

		info := map[string]interface{}{
			"fn_name": node.Function.String(),
			"args":    args,
		}

		m["calls"] = append(m["calls"].([]map[string]interface{}), info)

	case *ExpressionStatement:
		p.traverseNode(node.Expression, m)

	case *PrefixExpression:
		p.traverseNode(node.Right, m)
		return

	case *InfixExpression:
		p.traverseNode(node.Left, m)
		p.traverseNode(node.Right, m)
		return
	}
}

func (p *Parser) Info() map[string]interface{} {
	m := make(map[string]interface{})
	m["strings"] = make([]string, 0)
	m["regexs"] = make([]string, 0)
	m["numbers"] = make([]float64, 0)
	m["calls"] = make([]map[string]interface{}, 0) // name, arguments
	m["lists"] = make([]string, 0)
	m["identifiers"] = make([]string, 0)

	r := p.ParseRule()
	p.traverseNode(r.Statement, m)

	return m
}
