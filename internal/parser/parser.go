package parser

import (
	"fmt"
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
	"slices"
)

type Parser struct {
	Tokens      []lexer.Token
	diagnostics []diagnostics.Diagnostic
	position    uint32
}

func (p *Parser) Parse() (ast.Node, []diagnostics.Diagnostic) {
	p.position = 0
	p.diagnostics = nil

	program, err := p.parse_main_function()
	if err != nil {
		p.report(p.currentSpan(), "%s", err)
	}

	return program, slices.Clone(p.diagnostics)
}

func (p *Parser) LastAt() (*lexer.Token, error) {
	return p.get()
}

func (p *Parser) eof() bool {
	return p.position >= uint32(len(p.Tokens))
}

func (p *Parser) currentSpan() models.Span {
	if token, err := p.get(); err == nil {
		return token.Span
	}

	if len(p.Tokens) == 0 {
		return models.Span{}
	}

	last := p.Tokens[len(p.Tokens)-1].Span

	// Empty span immediately after the final token.
	return models.Span{
		Start: last.End,
		End:   last.End,
	}
}

func (p *Parser) diagnosis(span models.Span, level diagnostics.DiagnosticLevel, format string, args ...any) {
	diagnostics.Dbg("Filing diagnosis of type: ", level, ". span: ", span)
	p.diagnostics = append(p.diagnostics, diagnostics.Diagnostic{
		Span:  span,
		Msg:   fmt.Sprintf(format, args...),
		Level: level,
	})
}

func (p *Parser) report(span models.Span, format string, args ...any) {
	p.diagnosis(span, diagnostics.ERROR, format, args...)
}

func (p *Parser) warn(span models.Span, format string, args ...any) {
	p.diagnosis(span, diagnostics.WARNING, format, args...)
}

func (p *Parser) info(span models.Span, format string, args ...any) {
	p.diagnosis(span, diagnostics.INFO, format, args...)
}

func (p *Parser) at(position uint32) (*lexer.Token, error) {
	if position >= uint32(len(p.Tokens)) {
		return nil, ErrOutOfBounds
	}

	return &p.Tokens[position], nil
}

func (p *Parser) mark() uint32 {
	return p.position
}

func (p *Parser) spanFrom(start uint32) models.Span {
	// No tokens consumed.
	if start >= p.position {
		return models.Span{}
	}

	if start >= uint32(len(p.Tokens)) {
		return models.Span{}
	}

	endIndex := p.position - 1
	if endIndex >= uint32(len(p.Tokens)) {
		endIndex = uint32(len(p.Tokens)) - 1
	}

	first := p.Tokens[start].Span
	last := p.Tokens[endIndex].Span

	return models.Span{
		Start: first.Start,
		End:   last.End,
	}
}

func (p *Parser) infoFrom(start uint32) ast.NodeInfo {
	return ast.NodeInfo{
		Span: p.spanFrom(start),
	}
}

/* TODO: When golang 1.27 is out, change this to a generic method */
func parserTry[T any](p *Parser, fn func() (T, error)) (T, bool) {
	init_position := p.position
	val, err := fn()
	if err != nil {
		p.position = init_position
	}
	return val, err == nil
}

func (p *Parser) get() (*lexer.Token, error) {
	return p.at(p.position)
}

/* Runs a predicate against every token, if predicate was successful, returns true, if it got to EOF, returns false */
func (p *Parser) until(predicate func(lexer.Token) bool) bool {
	token, err := p.get()

	for err == nil {
		if predicate(*token) {
			break
		}

		p.position++
		token, err = p.get()
	}

	return err == nil
}

/* Skips to token, return true if token type was found */
func (p *Parser) skip_to_token(token_type lexer.TokenType) bool {
	return p.until(
		func(token lexer.Token) bool { return token.Type == token_type },
	)
}

func (p *Parser) skip_to_keyword(keyword lexer.Keyword) bool {
	return p.until(
		func(token lexer.Token) bool {
			kw, ok := lexer.MapToKeyword([]rune(token.Value))
			return token.Type == lexer.KEYWORD && ok && kw == keyword
		},
	)
}

func (p *Parser) eat_token(token_type lexer.TokenType) (string, error) {
	diagnostics.Dbg("Ate token: ", token_type)
	token, err := p.get()

	if err != nil {
		return "", err
	}

	if token.Type != token_type {
		diagnostics.Dbg("Didn't go well!")
		return "", ErrExpectedToken{Got: token.Type, Expected: token_type}
	}

	p.position++

	return token.Value, nil
}

func (p *Parser) eat_keyword(expected lexer.Keyword) (*lexer.Keyword, error) {
	diagnostics.Dbg("Ate keyword: ", expected)
	value, err := p.eat_token(lexer.KEYWORD)
	if err != nil {
		return nil, ErrExpectedKeyword{Expected: expected}
	}

	keyword, ok := lexer.MapToKeyword([]rune(value))
	if !ok {
		return nil, ErrUnrecognizedKeyword
	}

	if keyword != expected {
		return nil, ErrExpectedKeyword{Expected: expected}
	}

	return &keyword, nil
}
