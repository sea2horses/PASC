package neoparser

import "pseint-compiled/internal/lexer"

type Pattern[T any] interface {
	Match(*Context) Match[T]
}

type PatternFunc[T any] func(*Context) Match[T]

func (f PatternFunc[T]) Match(c *Context) Match[T] {
	return f(c)
}

/* Parse a token of the given type */
func Tok(expected lexer.TokenType) Pattern[lexer.Token] {
	return PatternFunc[lexer.Token](func(ctx *Context) Match[lexer.Token] {
		start := ctx.Mark()
		tok := ctx.Peek()

		if tok == nil {
			return Match[lexer.Token]{
				Kind:  NoMatch,
				Start: start,
				End:   start,
			}
		}

		if tok.Type != expected {
			return Match[lexer.Token]{
				Kind:  NoMatch,
				Start: tok.Span.Start.Offset,
				End:   tok.Span.End.Offset,
			}
		}

		ctx.Advance()

		return Match[lexer.Token]{
			Kind:  Matched,
			Value: *tok,
			Start: tok.Span.Start.Offset,
			End:   tok.Span.End.Offset,
		}
	})
}

func Kw(expected lexer.Keyword) Pattern[lexer.Keyword] {
	return PatternFunc[lexer.Keyword](func(ctx *Context) Match[lexer.Keyword] {
		start := ctx.Mark()
		tok := ctx.Peek()

		if tok == nil {
			return Match[lexer.Keyword]{
				Kind:  NoMatch,
				Start: start,
				End:   start,
			}
		}

		if tok.Type != lexer.KEYWORD {
			return Match[lexer.Keyword]{
				Kind:  NoMatch,
				Start: tok.Span.Start.Offset,
				End:   tok.Span.End.Offset,
			}
		}

		kw, ok := lexer.MapToKeyword([]rune(tok.Value))

		if !ok || kw != expected {
			return Match[lexer.Keyword]{
				Kind:  NoMatch,
				Start: tok.Span.Start.Offset,
				End:   tok.Span.End.Offset,
			}
		}

		return Match[lexer.Keyword]{
			Kind:  Matched,
			Value: kw,
			Start: tok.Span.Start.Offset,
			End:   tok.Span.End.Offset,
		}
	})
}
