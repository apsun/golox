package lox

import (
	"fmt"
)

type TokenType int

const (
	TokenTypeLeftParen TokenType = iota
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeSlash
	TokenTypeStar
	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeLess
	TokenTypeLessEqual
	TokenTypeQuestion
	TokenTypeColon
	TokenTypeIdentifier
	TokenTypeString
	TokenTypeNumber
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFun
	TokenTypeFor
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile
	TokenTypeEOF
)

var tokenTypeStringMap = map[TokenType]string{
	TokenTypeLeftParen:    "LeftParen",
	TokenTypeRightParen:   "RightParen",
	TokenTypeLeftBrace:    "LeftBrace",
	TokenTypeRightBrace:   "RightBrace",
	TokenTypeComma:        "Comma",
	TokenTypeDot:          "Dot",
	TokenTypeMinus:        "Minus",
	TokenTypePlus:         "Plus",
	TokenTypeSemicolon:    "Semicolon",
	TokenTypeSlash:        "Slash",
	TokenTypeStar:         "Star",
	TokenTypeBang:         "Bang",
	TokenTypeBangEqual:    "BangEqual",
	TokenTypeEqual:        "Equal",
	TokenTypeEqualEqual:   "EqualEqual",
	TokenTypeGreater:      "Greater",
	TokenTypeGreaterEqual: "GreaterEqual",
	TokenTypeLess:         "Less",
	TokenTypeLessEqual:    "LessEqual",
	TokenTypeQuestion:     "Question",
	TokenTypeColon:        "Colon",
	TokenTypeIdentifier:   "Identifier",
	TokenTypeString:       "String",
	TokenTypeNumber:       "Number",
	TokenTypeAnd:          "And",
	TokenTypeClass:        "Class",
	TokenTypeElse:         "Else",
	TokenTypeFalse:        "False",
	TokenTypeFun:          "Fun",
	TokenTypeFor:          "For",
	TokenTypeIf:           "If",
	TokenTypeNil:          "Nil",
	TokenTypeOr:           "Or",
	TokenTypePrint:        "Print",
	TokenTypeReturn:       "Return",
	TokenTypeSuper:        "Super",
	TokenTypeThis:         "This",
	TokenTypeTrue:         "True",
	TokenTypeVar:          "Var",
	TokenTypeWhile:        "While",
	TokenTypeEOF:          "EOF",
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
		literalStr = fmt.Sprintf("literal: %#v, ", t.literal)
	}
	return fmt.Sprintf(
		"Token{ty: %v, lexeme: %q, %sline: %d}",
		t.ty,
		t.lexeme,
		literalStr,
		t.line,
	)
}
