package neoparser

import (
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
)

type Context struct {
	Tokens []lexer.Token
	Pos    uint32
}

func (c *Context) EOF() bool {
	return c.Pos >= uint32(len(c.Tokens))
}

func (c *Context) Peek() *lexer.Token {
	if c.EOF() {
		return nil
	}

	return &c.Tokens[c.Pos]
}

func (c *Context) PeekAt(offset uint32) *lexer.Token {
	pos := c.Pos + offset

	if pos >= uint32(len(c.Tokens)) {
		return nil
	}

	return &c.Tokens[pos]
}

func (c *Context) Advance() {
	if !c.EOF() {
		c.Pos++
	}
}

func (c *Context) Mark() uint32 {
	return c.Pos
}

func (c *Context) Reset(pos uint32) {
	c.Pos = pos
}

func (c *Context) Span(start uint32, end uint32) models.Span {
	/* Get span from other token */
	t1 := c.PeekAt(start)
	t2 := c.PeekAt(end)

	if t1 == nil || t2 == nil {
		return models.Span{}
	}

	return models.JoinSpans(t1.Span, t2.Span)
}
