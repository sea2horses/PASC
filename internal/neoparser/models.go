package neoparser

import "pseint-compiled/internal/lexer"

type Context struct{}

type Pattern interface {
	Match(*Context) bool
}

type Captures map[string]any

type MatchKind uint8

const (
	NoMatch MatchKind = iota
	Match
	Failed
)

type Cursor struct {
	Tokens []lexer.Token
	Pos    int
}

type Result[T any] struct {
	Value     T
	Next      int
	MatchKind MatchKind
	// Error    *ParseError
}
