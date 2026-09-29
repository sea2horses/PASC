package neoparser

import (
	"pseint-compiled/internal/lexer"
	"slices"
)

type RecoveryStrategy interface {
	Recover(*Context)
}

func Recover[T any](
	p Pattern[T],
	strategy RecoveryStrategy,
) Pattern[T] {
	return PatternFunc[T](func(ctx *Context) Match[T] {
		result := p.Match(ctx)

		if result.Kind == Failed {
			strategy.Recover(ctx)
		}

		return result
	})
}

type UntilNewline struct{}

func (UntilNewline) Recover(ctx *Context) {
	for !ctx.EOF() {
		if ctx.Peek().Type == lexer.NEWLINE {
			return
		}

		ctx.Advance()
	}
}

type UntilKeywords struct {
	Keywords []lexer.Keyword
}

func (uk UntilKeywords) Recover(ctx *Context) {
	for !ctx.EOF() {
		tok := ctx.Peek()

		kw, ok := lexer.MapToKeyword([]rune(tok.Value))
		if tok.Type == lexer.KEYWORD && ok && slices.Contains(uk.Keywords, kw) {
			return
		}

		ctx.Advance()
	}
}
