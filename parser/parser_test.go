package parser

import (
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"(a * (b / c))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"(a + b) / c",
			"((a + b) / c)",
		},
		{
			"add(1,3) == true",
			"(add(1, 3) == true)",
		},
		{
			"a == \"category name\" OR true",
			"((a == \"category name\") OR true)",
		},
		{
			"a == r\"category name\" OR true",
			"((a == \"category name\") OR true)",
		},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := New(l)
		rule := p.ParseRule()
		checkParserErrors(t, p)
		actual := rule.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}

}

func TestInfo(t *testing.T) {
	l := NewLexer(`a == "category is not equal" OR (b == 10 AND c >=20.5) OR r"a.*" OR LOWER(a) NOT_CONTAINS @LIST_345324 AND BELOW(a1, b1, 10)`)
	p := New(l)

	info := p.Info()
	if info["lists"].([]string)[0] != "LIST_345324" {
		t.Errorf("expected=%q, got=%q", info["lists"].([]string)[0], "LIST_345324")
	}

}
