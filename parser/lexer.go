package parser

type Lexer struct {
	input        string
	position     int // current caracter position
	readPosition int //(next character in input)
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0x00
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	if l.ch != 0x00 {
		l.readPosition += 1
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	// skip whitespace characters
	l.skipWhitespace()

	switch l.ch {
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Literal = "<="
			tok.Type = LTE
		} else {
			tok = newToken(LT, l.ch)
		}

	case '>':
		tok = newToken(GT, l.ch)
		if l.peekChar() == '=' {
			l.readChar()
			tok.Literal = ">="
			tok.Type = GTE
		} else {
			tok = newToken(GT, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Literal = "=="
			tok.Type = EQUALS
		}
	case '/':
		tok = newToken(FSLASH, l.ch)
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '*':
		tok = newToken(ASTERIK, l.ch)
	case '%':
		tok = newToken(MODULUS, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Literal = "!="
			tok.Type = NOTEQUAL
		}
	case 'r':
		if l.peekChar() == '"' {
			l.readChar()
			l.readChar()
			tok.Literal = l.readString()
			tok.Type = REGEX
			l.readChar()
		} else {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		}
	case '@':
		l.readChar()
		tok.Literal = l.readIdentifier()
		tok.Type = LISTNAME
		return tok
	case '"':
		l.readChar()
		tok.Literal = l.readString()
		tok.Type = STRING
		l.readChar()

		return tok
	case 0x00:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
			return tok
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readString() string {
	pos := l.position

	prevCh := ""
	for (prevCh == "\\" && l.ch == '"') || l.ch != '"' {
		prevCh = string(l.ch)
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_' || ch >= '0' && ch <= '9'
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
