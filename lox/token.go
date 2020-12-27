package lox

import (
	"fmt"
)

type TokenType int

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual
	Identifier
	String
	Number
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
	EOF
)

var tokenTypeStringMap = map[TokenType]string{
	LeftParen:    "LeftParen",
	RightParen:   "RightParen",
	LeftBrace:    "LeftBrace",
	RightBrace:   "RightBrace",
	Comma:        "Comma",
	Dot:          "Dot",
	Minus:        "Minus",
	Plus:         "Plus",
	Semicolon:    "Semicolon",
	Slash:        "Slash",
	Star:         "Star",
	Bang:         "Bang",
	BangEqual:    "BangEqual",
	Equal:        "Equal",
	EqualEqual:   "EqualEqual",
	Greater:      "Greater",
	GreaterEqual: "GreaterEqual",
	Less:         "Less",
	LessEqual:    "LessEqual",
	Identifier:   "Identifier",
	String:       "String",
	Number:       "Number",
	And:          "And",
	Class:        "Class",
	Else:         "Else",
	False:        "False",
	Fun:          "Fun",
	For:          "For",
	If:           "If",
	Nil:          "Nil",
	Or:           "Or",
	Print:        "Print",
	Return:       "Return",
	Super:        "Super",
	This:         "This",
	True:         "True",
	Var:          "Var",
	While:        "While",
	EOF:          "EOF",
}

func (ty TokenType) String() string {
	return tokenTypeStringMap[ty]
}

type Token struct {
	ty      TokenType
	lexeme  string
	literal interface{}
	line    int
}

func (t Token) String() string {
	literalStr := ""
	if t.literal != nil {
		literalStr = fmt.Sprintf("literal: %v, ", t.literal)
	}
	return fmt.Sprintf(
		"Token{ty: %v, lexeme: %q, %sline: %d}",
		t.ty,
		t.lexeme,
		literalStr,
		t.line,
	)
}
