package parser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
)

func (p *Parser) parse_write() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing write...")
	start := p.mark()
	_, err := p.eat_keyword(lexer.ESCRIBIR)
	if err != nil {
		return nil, err
	}
	/* Expressions to be printed */
	print, err := p.parse_expression_list()
	if err != nil {
		return nil, err
	}

	return &ast.Write{Print: print, NodeInfo: p.infoFrom(start)}, nil
}

func (p *Parser) parse_read() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing read...")
	start := p.mark()
	_, err := p.eat_keyword(lexer.LEER)
	if err != nil {
		return nil, err
	}
	/* Expressions to be read */
	into, err := p.parse_expression_list()
	if err != nil {
		return nil, err
	}

	return &ast.Read{Into: into, NodeInfo: p.infoFrom(start)}, nil
}

func (p *Parser) parse_clear_screen() (ast.Stmt, error) {
	start := p.mark()
	_, err := p.eat_keyword(lexer.BORRAR)
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.PANTALLA)
	if err != nil {
		return nil, err
	}

	return &ast.ClearScreen{NodeInfo: p.infoFrom(start)}, nil
}
