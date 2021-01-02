package lox

import (
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    TokenTypeAnd,
	"class":  TokenTypeClass,
	"else":   TokenTypeElse,
	"false":  TokenTypeFalse,
	"for":    TokenTypeFor,
	"fun":    TokenTypeFun,
	"if":     TokenTypeIf,
	"nil":    TokenTypeNil,
	"or":     TokenTypeOr,
	"print":  TokenTypePrint,
	"return": TokenTypeReturn,
	"super":  TokenTypeSuper,
	"this":   TokenTypeThis,
	"true":   TokenTypeTrue,
	"var":    TokenTypeVar,
	"while":  TokenTypeWhile,
	"break":  TokenTypeBreak,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
	errors  []*SyntaxError
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, []*SyntaxError) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{
		ty:      TokenTypeEOF,
		lexeme:  "(EOF)",
		literal: nil,
		line:    s.line,
	})
	return s.tokens, s.errors
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(TokenTypeLeftParen)
	case ')':
		s.addToken(TokenTypeRightParen)
	case '{':
		s.addToken(TokenTypeLeftBrace)
	case '}':
		s.addToken(TokenTypeRightBrace)
	case ',':
		s.addToken(TokenTypeComma)
	case '.':
		s.addToken(TokenTypeDot)
	case '-':
		s.addToken(TokenTypeMinus)
	case '+':
		s.addToken(TokenTypePlus)
	case ';':
		s.addToken(TokenTypeSemicolon)
	case '*':
		s.addToken(TokenTypeStar)
	case '!':
		if s.match('=') {
			s.addToken(TokenTypeBangEqual)
		} else {
			s.addToken(TokenTypeBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TokenTypeEqualEqual)
		} else {
			s.addToken(TokenTypeEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(TokenTypeLessEqual)
		} else {
			s.addToken(TokenTypeLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenTypeGreaterEqual)
		} else {
			s.addToken(TokenTypeGreater)
		}
	case '?':
		s.addToken(TokenTypeQuestion)
	case ':':
		s.addToken(TokenTypeColon)
	case '/':
		if s.match('/') {
			s.scanLineComment()
		} else if s.match('*') {
			s.scanBlockComment()
		} else {
			s.addToken(TokenTypeSlash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			s.errors = append(s.errors, NewSyntaxError(
				s.line, nil, "unexpected character",
			))
		}
	}
}

func (s *Scanner) scanLineComment() {
	for !s.isAtEnd() && s.peek() != '\n' {
		s.advance()
	}
}

func (s *Scanner) scanBlockComment() {
	for !s.isAtEnd() {
		if s.peek() == '/' && s.peekNext() == '*' {
			// Handle nested block comment... Something something
			// recursion in the compiler, doing this properly is a
			// lot more difficult though so this will do for now.
			s.advance()
			s.advance()
			s.scanBlockComment()
		} else if s.peek() == '*' && s.peekNext() == '/' {
			s.advance()
			s.advance()
			return
		} else {
			s.advance()
		}
	}
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Check for a decimal point with a numeric value after it
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic(err)
	}

	s.addTokenWithLiteral(TokenTypeNumber, num)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	// Check if it's actually a keyword
	text := s.source[s.start:s.current]
	ty, ok := keywords[text]
	if !ok {
		ty = TokenTypeIdentifier
	}

	s.addToken(ty)
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errors = append(s.errors, NewSyntaxError(
			s.line, nil, "unterminated string",
		))
		return
	}

	// Consume closing quote
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TokenTypeString, value)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) addTokenWithLiteral(ty TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		ty:      ty,
		lexeme:  text,
		literal: literal,
		line:    s.line,
	})
}
