package parser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
)

func (p *Parser) parse_statement_block() ([]ast.Stmt, error) {
	diagnostics.Dbg("Parsing statement block...")
	tk, _ := p.get()
	diagnostics.Dbgfmt("On token: %+v\n", tk)
	var stmts []ast.Stmt
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
	return stmts, nil
}

func (p *Parser) parse_statement() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing statement...")
	tk, _ := p.get()
	diagnostics.Dbgfmt("On token: %+v\n", tk)
	token, err := p.get()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.KEYWORD:
		return p.parse_keyword()
	default:
		{
			assignment, ok := parserTry(p, p.parse_assignment)
			if ok {
				return assignment, nil
			}

			return nil, nil
		}
	}
}

func (p *Parser) parse_main_function() (*ast.MainFunction, error) {
	diagnostics.Dbg("Parsing main functions...")
	start := p.mark()
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

	info := p.infoFrom(start)

	stmts, err := p.parse_statement_block()
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.FINALGORITMO)
	if err != nil {
		return nil, err
	}

	return &ast.MainFunction{Name: fn_name, Stmts: stmts, NodeInfo: info}, nil
}

func (p *Parser) parse_assignment() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing assignment...")
	tk, _ := p.get()
	start := p.mark()
	diagnostics.Dbgfmt("On token: %+v\n", tk)
	target, err := p.parse_expression()
	if err != nil {
		return nil, err
	}
	_, err = p.eat_token(lexer.EQUALS)
	if err != nil {
		return nil, err
	}
	content, err := p.parse_expression()
	if err != nil {
		return nil, err
	}
	return &ast.Assignment{Target: target, Content: content, NodeInfo: p.infoFrom(start)}, nil
}

func (p *Parser) parse_keyword() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing keyword...")
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
	case lexer.DEFINIR:
		return p.parse_declaration()
	case lexer.ESCRIBIR:
		return p.parse_write()
	case lexer.LEER:
		return p.parse_read()
	case lexer.SI:
		return p.parse_if()
	case lexer.MIENTRAS:
		return p.parse_while()
	case lexer.BORRAR:
		return p.parse_clear_screen()
	}

	return nil, nil
}

func (p *Parser) parse_declaration() (ast.Stmt, error) {
	start := p.mark()
	/* GO */
	_, err := p.eat_keyword(lexer.DEFINIR)
	if err != nil {
		return nil, err
	}

	/* Var name */
	name, err := p.eat_token(lexer.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.COMO)
	if err != nil {
		return nil, err
	}

	typeref, err := p.parse_type()
	if err != nil {
		return nil, err
	}

	return &ast.Declaration{
		NodeInfo: p.infoFrom(start),
		Name:     name,
		Type:     typeref,
	}, nil
}

