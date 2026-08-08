package parser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
)

func (p *Parser) parse_if() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing if...")
	start := p.mark()
	_, err := p.eat_keyword(lexer.SI)
	if err != nil {
		return nil, err
	}

	condition, err := p.parse_expression()
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.ENTONCES)
	if err != nil {
		return nil, err
	}

	/* Parse if statement block */
	stmts, err := p.parse_statement_block()
	if err != nil {
		return nil, err
	}

	var else_branch *ast.Else = nil

	/* Check for else branch */
	if tok, err := p.get(); err == nil && tok.Type == lexer.KEYWORD {
		if keyword, ok := lexer.MapToKeyword([]rune(tok.Value)); ok && keyword == lexer.SINO {
			else_branch, err = p.parse_else()
			if err != nil {
				return nil, err
			}
		}
	}

	_, err = p.eat_keyword(lexer.FINSI)
	if err != nil {
		return nil, err
	}

	return &ast.If{Condition: condition, Stmts: stmts, Else: else_branch, NodeInfo: p.infoFrom(start)}, nil
}

func (p *Parser) parse_else() (*ast.Else, error) {
	start := p.mark()

	_, err := p.eat_keyword(lexer.SINO)
	if err != nil {
		return nil, err
	}

	/* Parse if statement block */
	stmts, err := p.parse_statement_block()
	if err != nil {
		return nil, err
	}

	return &ast.Else{Stmts: stmts, NodeInfo: p.infoFrom(start)}, nil
}

func (p *Parser) parse_while() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing while...")
	start := p.mark()
	_, err := p.eat_keyword(lexer.MIENTRAS)
	if err != nil {
		return nil, err
	}

	condition, err := p.parse_expression()
	if err != nil {
		return nil, err
	}

	/* TODO: Make it a language option for optional 'Hacer' */
	_, err = p.eat_keyword(lexer.HACER)
	if err != nil {
		return nil, err
	}

	stmts, err := p.parse_statement_block()
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.FINMIENTRAS)
	if err != nil {
		return nil, err
	}

	return &ast.While{Condition: condition, Stmts: stmts, NodeInfo: p.infoFrom(start)}, nil
}
