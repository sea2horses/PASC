package parser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
)

func (p *Parser) parse_operand() (ast.Expr, error) {
	diagnostics.Dbg("Parsing operand!")
	start := p.mark()
	// Try parsing a unary operator
	op := p.parse_operator()
	if op == nil {
		return p.parse_primary()
	} else {
		rhs, err := p.parse_operand()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOperation{Op: *op, Expr: rhs, NodeInfo: p.infoFrom(start)}, nil
	}
}

func (p *Parser) parse_primary() (ast.Expr, error) {
	diagnostics.Dbg("Parsing primary!")
	token, err := p.get()
	if err != nil {
		return nil, err
	}

	start := p.mark()

	switch token.Type {
	case lexer.STRING_LITERAL:
		val, _ := p.eat_token(lexer.STRING_LITERAL)
		return &ast.StringLiteral{Content: val, NodeInfo: p.infoFrom(start)}, nil
	case lexer.NUMBER_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		num, _ := lexer.MapToNumber([]rune(val))
		var frac uint64 = 0

		/* Parse decimal */
		if tok, err := p.get(); err == nil && tok.Type == lexer.DOT {
			p.eat_token(lexer.DOT)
			val, err := p.eat_token(lexer.NUMBER_LITERAL)
			if err != nil {
				p.Report(p.currentSpan(), "expected decimal part")
				return nil, err
			}
			frac, _ = lexer.MapToNumber([]rune(val))
		}
		return &ast.NumberLiteral{Int: num, Frac: frac, NodeInfo: p.infoFrom(start)}, nil
	case lexer.BOOLEAN_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		vbool, _ := lexer.MapToBool([]rune(val))
		return &ast.BoolLiteral{Value: vbool, NodeInfo: p.infoFrom(start)}, nil
	case lexer.IDENTIFIER:
		val, _ := p.eat_token(lexer.IDENTIFIER)
		return &ast.Variable{Name: val, NodeInfo: p.infoFrom(start)}, nil
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

		expr.Info().JoinSpan(
			p.spanFrom(start),
		)
		return expr, nil
	}

	return nil, ErrExpectedOperand
}
