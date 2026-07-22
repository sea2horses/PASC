package parser

import (
	"fmt"
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/operators"
	"slices"
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
	ErrUnexpectedNode      = ParserError("unexpected node in expression")
	ErrInvalidExpression   = ParserError("invalid expression")
)

type ErrExpectedToken struct {
	Got      lexer.TokenType
	Expected lexer.TokenType
}

func (er ErrExpectedToken) Error() string {
	return fmt.Sprintf("expected token type: %s, got: %s", er.Expected, er.Got)
}

type ErrExpectedKeyword struct {
	Expected lexer.Keyword
}

func (er ErrExpectedKeyword) Error() string {
	return fmt.Sprintf("expected keyword: %s", er.Expected)
}

type Parser struct {
	Tokens   []lexer.Token
	position uint32
}

func (p *Parser) Parse() (ast.Node, error) {
	return p.parse_main_function()
}

func (p *Parser) LastAt() (*lexer.Token, error) {
	return p.get()
}

func (p *Parser) at(position uint32) (*lexer.Token, error) {
	if position >= uint32(len(p.Tokens)) {
		return nil, ErrOutOfBounds
	}

	return &p.Tokens[position], nil
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

func (p *Parser) parse_primary() (ast.Expr, error) {
	diagnostics.Dbg("Parsing primary!")
	token, err := p.get()
	if err != nil {
		return nil, err
	}

	switch token.Type {
	case lexer.STRING_LITERAL:
		val, _ := p.eat_token(lexer.STRING_LITERAL)
		return &ast.StringLiteral{Content: val}, nil
	case lexer.NUMBER_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		num, _ := lexer.MapToNumber([]rune(val))
		var frac uint64 = 0

		/* Parse decimal */
		if tok, err := p.get(); err != nil && tok.Type == lexer.DOT {
			p.eat_token(lexer.DOT)
			val, err := p.eat_token(lexer.NUMBER_LITERAL)
			if err != nil {
				return nil, err
			}
			frac, _ = lexer.MapToNumber([]rune(val))
		}
		return &ast.NumberLiteral{Int: num, Frac: frac}, nil
	case lexer.BOOLEAN_LITERAL:
		val, _ := p.eat_token(lexer.NUMBER_LITERAL)
		vbool, _ := lexer.MapToBool([]rune(val))
		return &ast.BoolLiteral{Value: vbool}, nil
	case lexer.IDENTIFIER:
		val, _ := p.eat_token(lexer.IDENTIFIER)
		return &ast.Variable{Name: val}, nil
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
	op_type := p.parse_operator_type()
	diagnostics.Dbg("Resolved Type: ", op_type)

	if op_type == operators.UNRECOGNIZED {
		return nil
	} else {
		return &ast.Operator{Type: op_type}
	}
}

func (p *Parser) parse_operand() (ast.Expr, error) {
	diagnostics.Dbg("Parsing operand!")
	// Try parsing a unary operator
	op := p.parse_operator()
	if op == nil {
		return p.parse_primary()
	} else {
		rhs, err := p.parse_operand()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOperation{Op: *op, Expr: rhs}, nil
	}
}

func (p *Parser) read_infix_to_postfix() ([]ast.Node, error) {
	operator_stack := []ast.Operator{}
	postfix_stack := []ast.Node{}

	for {
		/* Parse operand */
		expr, err := p.parse_operand()
		if err != nil {
			return nil, err
		}

		/* Add it directly to the postfix stack */
		postfix_stack = append(postfix_stack, expr)

		/* Grab operator */
		op := p.parse_operator()
		/* If there are no more operators, break */
		if op == nil {
			break
		}

		/* Else, if the current operator is lesser or equal than the one on top, we pop the stack */
		for len(operator_stack) > 0 &&
			operators.BinaryOperatorPrecedence(operator_stack[len(operator_stack)-1].Type) >= operators.BinaryOperatorPrecedence(op.Type) {
			/* Append top operator to postfix stack */
			postfix_stack = append(postfix_stack, operator_stack[len(operator_stack)-1])
			/* Chop last. element */
			operator_stack = operator_stack[0 : len(operator_stack)-1]
		}

		operator_stack = append(operator_stack, *op)
	}

	/* Unload the remaining operators */
	for len(operator_stack) > 0 {
		/* Append top operator to postfix stack */
		postfix_stack = append(postfix_stack, operator_stack[len(operator_stack)-1])
		/* Chop last. element */
		operator_stack = operator_stack[0 : len(operator_stack)-1]
	}

	return postfix_stack, nil
}

func (p *Parser) postfix_to_expression(chain []ast.Node) (ast.Expr, error) {
	for {
		diagnostics.Dbg("Postfix Iteration, expr: ", chain)
		swapped := false
		/* Evaluate postfix */
		for i := range len(chain) - 2 {
			first, ok := chain[i].(ast.Expr)
			if !ok {
				continue
			}
			second, ok := chain[i+1].(ast.Expr)
			if !ok {
				continue
			}
			third, ok := chain[i+2].(ast.Operator)
			if !ok {
				continue
			}

			swapped = true
			/* Replace with new expression */
			chain = slices.Delete(chain, i, i+3)
			chain = slices.Insert(chain, i, ast.Node(ast.BinaryOperation{LHS: first, RHS: second, Op: third}))
		}

		if len(chain) == 1 {
			diagnostics.Dbg("Expression made: ", chain[0])
			return chain[0], nil
		}

		if !swapped {
			return nil, ErrInvalidExpression
		}
	}
}

func (p *Parser) parse_expression() (ast.Expr, error) {
	diagnostics.Dbg("Parsing expression...")
	tk, _ := p.get()
	diagnostics.Dbgfmt("On token: %+v\n", tk)
	postfix, err := p.read_infix_to_postfix()
	if err != nil {
		return nil, err
	}
	return p.postfix_to_expression(postfix)
}

func (p *Parser) parse_main_function() (*ast.MainFunction, error) {
	diagnostics.Dbg("Parsing main functions...")
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

	stmts, err := p.parse_statement_block()
	if err != nil {
		return nil, err
	}

	_, err = p.eat_keyword(lexer.FINALGORITMO)
	if err != nil {
		return nil, err
	}

	return &ast.MainFunction{Name: fn_name, Stmts: stmts}, nil
}

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

func (p *Parser) parse_assignment() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing assignment...")
	tk, _ := p.get()
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
	return ast.Assignment{Target: target, Content: content}, nil
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
	case lexer.ESCRIBIR:
		return p.parse_write()
	case lexer.MIENTRAS:
		return p.parse_while()
	}

	return nil, nil
}

func (p *Parser) parse_write() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing write...")
	_, err := p.eat_keyword(lexer.ESCRIBIR)
	if err != nil {
		return nil, err
	}
	/* Expression to be printed */
	expr, err := p.parse_expression()
	if err != nil {
		return nil, err
	}

	return &ast.Write{Print: expr}, nil
}

func (p *Parser) parse_while() (ast.Stmt, error) {
	diagnostics.Dbg("Parsing while...")
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

	return ast.While{Condition: condition, Stmts: stmts}, nil
}
