package parser

import "strings"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	//Special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	DOLLAR  = "$"
	ATSIGN  = "@"

	// Operators
	EQUALS      = "=="
	GT          = ">"
	LT          = "<"
	GTE         = ">="
	LTE         = "<="
	ASTERIK     = "*"
	PLUS        = "+"
	MINUS       = "-"
	FSLASH      = "/"
	MODULUS     = "%"
	AND         = "AND"
	OR          = "OR"
	NOTEQUAL    = "!="
	CONTAINS    = "CONTAINS"
	NOTCONTAINS = "NOT_CONTAINS"

	LPAREN      = "("
	RPAREN      = ")"
	DOUBLEQUOTE = "\""
	COMMA       = ","

	// Identifiers & literals
	FN       = "FN"
	IDENT    = "IDENT"
	STRING   = "STRING"
	NUMBER   = "NUMBER"
	REGEX    = "REGEX"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	LISTNAME = "LISTNAME"
)

var keywords = map[string]TokenType{
	"and":          AND,
	"or":           OR,
	"contains":     CONTAINS,
	"not_contains": NOTCONTAINS,
	"true":         TRUE,
	"false":        FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}

	return IDENT
}
