package parser

import (
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `a == "category is not equal" OR (b == 10 AND c >=20.5) r"a.*"  LOWER(a)  != CONTAINS NOT_CONTAINS @LIST_345324 a BELOW(10) b`

	tests := []struct {
		expected        TokenType
		expectedLiteral string
	}{
		{IDENT, "a"},
		{EQUALS, "=="},

		{STRING, "category is not equal"},
		{OR, "OR"},
		{LPAREN, "("},
		{IDENT, "b"},
		{EQUALS, "=="},
		{NUMBER, "10"},
		{AND, "AND"},
		{IDENT, "c"},
		{GTE, ">="},
		{NUMBER, "20.5"},
		{RPAREN, ")"},
		{REGEX, "a.*"},
		{IDENT, "LOWER"},
		{LPAREN, "("},
		{IDENT, "a"},
		{RPAREN, ")"},
		{NOTEQUAL, "!="},
		{CONTAINS, "CONTAINS"},
		{NOTCONTAINS, "NOT_CONTAINS"},
		{LISTNAME, "LIST_345324"},
		{IDENT, "a"},
		{IDENT, "BELOW"},
		{LPAREN, "("},
		{NUMBER, "10"},
		{RPAREN, ")"},
		{IDENT, "b"},
	}

	lex := NewLexer(input)

	for i, tt := range tests {

		tok := lex.NextToken()
		// fmt.Println(tok)
		if tok.Type != tt.expected {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expected, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
