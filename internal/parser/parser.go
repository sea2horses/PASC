package parser

import (
	"fmt"
	"pseint-compiled/internal/lexer"
)

type ParserError string

func (pe ParserError) Error() string {
	return string(pe)
}

const (
	ErrOutOfBounds         = ParserError("index is out of bounds")
	ErrUnrecognizedKeyword = ParserError("unrecognized keyword")
	ErrNotImplemented      = ParserError("not implemented")
	ErrExpectedOperand     = ParserError("expected operand")
)

type ErrExpectedToken struct {
	Got      lexer.TokenType
	Expected lexer.TokenType
}

func (er ErrExpectedToken) Error() string {
	return fmt.Sprintf("expected token type: %s, got: %s", er.Expected, er.Got)
}

type ErrExpectedKeyword struct {
	Got      lexer.Keyword
	Expected lexer.Keyword
}

func (er ErrExpectedKeyword) Error() string {
	return fmt.Sprintf("expected keyword: %s, got: %s", er.Expected, er.Got)
}

type Parser struct {
	Tokens   []lexer.Token
	position uint32
}

func (p *Parser) Parse() (Node, error) {
	return p.parse_main_function()
}

func (p *Parser) at(position uint32) (*lexer.Token, error) {
	if position >= uint32(len(p.Tokens)) {
		return nil, ErrOutOfBounds
	}

	return &p.Tokens[position], nil
}

func (p *Parser) get() (*lexer.Token, error) {
	return p.at(p.position)
}

func (p *Parser) eat_token(token_type lexer.TokenType) (string, error) {
	token, err := p.get()

	if err != nil {
		return "", err
	}

	if token.Type != token_type {
		return "", ErrExpectedToken{Got: token.Type, Expected: token_type}
	}

	p.position++

	return token.Value, nil
}

func (p *Parser) eat_keyword(expected lexer.Keyword) (*lexer.Keyword, error) {
	value, err := p.eat_token(lexer.KEYWORD)
	if err != nil {
		return nil, err
	}

	keyword, ok := lexer.MapToKeyword([]rune(value))
	if !ok {
		return nil, ErrUnrecognizedKeyword
	}

	if keyword != expected {
		return nil, ErrExpectedKeyword{Got: keyword, Expected: expected}
	}

	return &keyword, nil
}

func (p *Parser) parse_primary() (Expr, error) {
	token, err := p.get()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.STRING_LITERAL:
		val, _ := p.eat_token(lexer.STRING_LITERAL)
		return &StringLiteral{Content: val}, nil
	case lexer.NUMBER_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		num, _ := lexer.MapToNumber([]rune(val))
		var frac uint64 = 0

		/* Parse decimal */
		if tok, err := p.get(); err != nil && tok.Type == lexer.DOT {
			p.eat_token(lexer.DOT)
			val, _ := p.eat_token(lexer.NUMBER_LITERAL)
			frac, _ = lexer.MapToNumber([]rune(val))
		}
		return &NumberLiteral{Int: num, Frac: frac}, nil
	case lexer.BOOLEAN_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		vbool, _ := lexer.MapToBool([]rune(val))
		return &BoolLiteral{Value: vbool}, nil
	case lexer.IDENTIFIER:
		val, _ := p.eat_token(lexer.IDENTIFIER)
		return &Variable{Name: val}, nil
	case lexer.L_PARENTHESES:
		p.eat_token(lexer.L_PARENTHESES)
		expr, err := p.parse_expression()
		if err != nil {
			return nil, err
		}
		_, err = p.eat_token(lexer.R_PARENTHESES)
		if err != nil {
			return nil, err
		}
		return expr, nil
	}

	return nil, ErrExpectedOperand
}

func (p *Parser) parse_expression() (Expr, error) {
	return p.parse_primary()
}

func (p *Parser) parse_main_function() (*MainFunction, error) {
	var err error
	/* Set the proper structure */
	_, err = p.eat_keyword(lexer.ALGORITMO)
	if err != nil {
		return nil, err
	}

	fn_name, err := p.eat_token(lexer.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	var stmts []Stmt
	for {
		stmt, err := p.parse_statement()
		if err != nil {
			return nil, err
		}

		if stmt == nil {
			break
		}

		stmts = append(stmts, stmt)
	}

	_, err = p.eat_keyword(lexer.FINALGORITMO)
	if err != nil {
		return nil, err
	}

	return &MainFunction{Name: fn_name, Stmts: stmts}, nil
}

func (p *Parser) parse_statement() (Stmt, error) {
	token, err := p.get()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.KEYWORD:
		return p.parse_keyword()
	default:
		return nil, nil
	}
}

func (p *Parser) parse_keyword() (Stmt, error) {
	token, err := p.get()

	if err != nil {
		return nil, err
	}

	if token.Type != lexer.KEYWORD {
		return nil, ErrExpectedToken{Got: token.Type, Expected: lexer.KEYWORD}
	}

	kw, ok := lexer.MapToKeyword([]rune(token.Value))
	if !ok {
		return nil, ErrUnrecognizedKeyword
	}

	switch kw {
	case lexer.ESCRIBIR:
		return p.parse_write()
	}

	return nil, nil
}

func (p *Parser) parse_write() (Stmt, error) {
	p.eat_keyword(lexer.ESCRIBIR)
	/* Expression to be printed */
	expr, err := p.parse_expression()
	if err != nil {
		return nil, err
	}

	return &Write{Print: expr}, nil
}
