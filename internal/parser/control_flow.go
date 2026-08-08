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
		p.Report(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.SI})
	}

	condition, err := p.parse_expression()

	_, err = p.eat_keyword(lexer.ENTONCES)
	if err != nil {
		p.Warn(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.ENTONCES})
	}

	_, err = p.eat_token(lexer.NEWLINE)
	if err != nil {
		p.Report(p.currentSpan(), "expected newline")
		p.skip_to_nextline()
	}

	info := p.infoFrom(start)

	/* Parse if statement block */
	stmts, err := p.parse_statement_block(lexer.FINSI, lexer.SINO)

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
		p.Report(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.FINSI})
		p.Info(info.Span, "%s", "declared here")
	}

	return &ast.If{Condition: condition, Stmts: stmts, Else: else_branch, NodeInfo: info}, nil
}

func (p *Parser) parse_else() (*ast.Else, error) {
	start := p.mark()

	_, err := p.eat_keyword(lexer.SINO)
	if err != nil {
		p.Report(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.SINO})
	}

	_, err = p.eat_token(lexer.NEWLINE)
	if err != nil {
		p.Report(p.currentSpan(), "expected newline")
	}

	info := p.infoFrom(start)

	/* Parse if statement block */
	stmts, err := p.parse_statement_block(lexer.FINSI)
	return &ast.Else{Stmts: stmts, NodeInfo: info}, nil
}

func (p *Parser) parse_while() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing while...")
	start := p.mark()
	_, err := p.eat_keyword(lexer.MIENTRAS)
	if err != nil {
		p.Report(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.MIENTRAS})
	}

	condition, err := p.parse_expression()

	/* TODO: Make it a language option for optional 'Hacer' */
	_, err = p.eat_keyword(lexer.HACER)
	if err != nil {
		p.Warn(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.HACER})
	}

	_, err = p.eat_token(lexer.NEWLINE)
	if err != nil {
		p.Report(p.currentSpan(), "expected newline")
	}

	info := p.infoFrom(start)

	stmts, err := p.parse_statement_block(lexer.FINMIENTRAS)

	_, err = p.eat_keyword(lexer.FINMIENTRAS)
	if err != nil {
		p.Report(p.currentSpan(), "%s", ErrExpectedKeyword{Expected: lexer.FINMIENTRAS})
		p.Info(info.Span, "declared here")
	}

	return &ast.While{Condition: condition, Stmts: stmts, NodeInfo: info}, nil
}
