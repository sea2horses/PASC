package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
	"slices"
)

func Operand() Pattern[ast.Expr] {
	return Bind(
		Map(
			Seq2(
				Many(
					UnaryOperator(),
				),
				Atom(),
			),
			func(p Pair[[]*ast.Operator, ast.Expr]) ast.Expr {
				new_expr := p.Second
				slices.Reverse(p.First)
				for _, o := range p.First {
					new_expr = &ast.UnaryOperation{
						NodeInfo: ast.NodeInfo{
							Span: models.JoinSpans(new_expr.NodeSpan(), o.Span),
						},
						Op:   *o,
						Expr: new_expr,
					}
				}
				return new_expr
			},
		),
		func(exp ast.Expr) Pattern[ast.Expr] {
			return Map(
				Optional(Postfix(exp)),
				func(o OptionalValue[ast.Expr]) ast.Expr {
					return o.Or(exp)
				},
			)
		},
	)
}

func Postfix(target ast.Expr) Pattern[ast.Expr] {
	return OneOf(
		Index(target),
		Call(target),
	)
}

func Index(target ast.Expr) Pattern[ast.Expr] {
	return MapWithSpan(
		Many1(
			Between(
				Tok(lexer.L_BRACKET),
				Expression(),
				Tok(lexer.R_BRACKET),
			),
		),
		func(indices []ast.Expr, span models.Span) ast.Expr {
			return &ast.Index{
				NodeInfo: ast.NodeInfo{
					Span: models.JoinSpans(target.NodeSpan(), span),
				},
				Indexes: indices,
				Target:  target,
			}
		},
	)
}

func Call(target ast.Expr) Pattern[ast.Expr] {
	return MapWithSpan(
		Between(
			Tok(lexer.L_PARENTHESES),
			SepBy(Expression(), Tok(lexer.COMMA)),
			Tok(lexer.R_PARENTHESES),
		),
		func(arguments []ast.Expr, span models.Span) ast.Expr {
			return &ast.Call{
				NodeInfo: ast.NodeInfo{
					Span: models.JoinSpans(target.NodeSpan(), span),
				},
				Callable:  target,
				Arguments: arguments,
			}
		},
	)
}

func Expression() Pattern[ast.Expr] {}
