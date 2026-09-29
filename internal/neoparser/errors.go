package neoparser

import (
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
)

type ParseError struct {
	Position uint32

	Message  string
	Expected []Expectation
	Found    *lexer.Token

	Help   string
	Labels []Label
}

type ExpectationKind uint8

const (
	ExpectToken ExpectationKind = iota
	ExpectKeyword
	ExpectExpression
	ExpectType
	ExpectCustom
)

type Expectation struct {
	Kind ExpectationKind

	Token   lexer.TokenType
	Keyword lexer.Keyword
	Label   string
}

type Label struct {
	Span    models.Span
	Message string
}
