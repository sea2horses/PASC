package parser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/operators"
)

func (p *Parser) parse_operator_type() operators.OperatorType {
	init_position := p.position
	tok, err := p.get()
	if err != nil {
		return operators.UNRECOGNIZED
	}

	switch tok.Type {
	case lexer.PLUS:
		p.eat_token(lexer.PLUS)
		return operators.ADD
	case lexer.MINUS:
		p.eat_token(lexer.MINUS)
		return operators.SUBTRACT
	case lexer.ASTERISK:
		p.eat_token(lexer.ASTERISK)
		return operators.MULTIPLY
	case lexer.SLASH:
		p.eat_token(lexer.SLASH)
		return operators.DIVIDE
	case lexer.L_ANGLE:
		p.eat_token(lexer.L_ANGLE)

		if tok, err := p.get(); err == nil && tok.Type == lexer.EQUALS {
			p.eat_token(lexer.EQUALS)
			return operators.LESSER_EQ
		}
		return operators.LESSER
	case lexer.R_ANGLE:
		p.eat_token(lexer.R_ANGLE)
		if tok, err := p.get(); err == nil && tok.Type == lexer.EQUALS {
			p.eat_token(lexer.EQUALS)
			return operators.GREATER_EQ
		}
		return operators.GREATER
	case lexer.EQUALS:
		p.eat_token(lexer.EQUALS)
		if tok, err := p.get(); err == nil && tok.Type == lexer.EQUALS {
			p.eat_token(lexer.EQUALS)
			return operators.EQUALS
		}
		/* Restart position */
		p.position = init_position
		return operators.UNRECOGNIZED
	case lexer.EX_MARK:
		p.eat_token(lexer.EX_MARK)
		if tok, err := p.get(); err == nil && tok.Type == lexer.EQUALS {
			p.eat_token(lexer.EQUALS)
			return operators.NOT_EQUALS
		}
		return operators.NOT
	/* This is fucking up everything */
	case lexer.KEYWORD:
		_, err := p.eat_keyword(lexer.MOD)
		if err != nil {
			/* Restart position */
			p.position = init_position
			return operators.UNRECOGNIZED
		}
		return operators.MODULO
	case lexer.AMPERSAND:
		p.eat_token(lexer.AMPERSAND)
		if tok, err := p.get(); err == nil && tok.Type == lexer.AMPERSAND {
			p.eat_token(lexer.AMPERSAND)
			return operators.AND
		}
		/* Restart position */
		p.position = init_position
		return operators.UNRECOGNIZED
	case lexer.PIPE:
		p.eat_token(lexer.PIPE)
		if tok, err := p.get(); err == nil && tok.Type == lexer.PIPE {
			p.eat_token(lexer.PIPE)
			return operators.OR
		}
		/* Restart position */
		p.position = init_position
		return operators.UNRECOGNIZED
	}
	return operators.UNRECOGNIZED
}

func (p *Parser) parse_operator() *ast.Operator {
	diagnostics.Dbg("Parsing operator!")
	start := p.mark()
	op_type := p.parse_operator_type()
	diagnostics.Dbg("Resolved Type: ", op_type)

	if op_type == operators.UNRECOGNIZED {
		return nil
	} else {
		return &ast.Operator{Type: op_type, NodeInfo: p.infoFrom(start)}
	}
}
