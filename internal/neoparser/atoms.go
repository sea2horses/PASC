package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
)

func StringLiteral() Pattern[ast.Expr] {
	return Map(
		Tok(lexer.STRING_LITERAL),
		func(tok lexer.Token) ast.Expr {
			return &ast.StringLiteral{
				Content: tok.Value,
				NodeInfo: ast.NodeInfo{
					Span: tok.Span,
				},
			}
		},
	)
}

func RawNumber() Pattern[string] {
	return Map(
		Tok(lexer.NUMBER_LITERAL),
		func(tok lexer.Token) string {
			return tok.Value
		},
	)
}

func NumberLiteral() Pattern[ast.Expr] {
	return MapWithSpan(
		Seq2(
			RawNumber(),
			Optional(
				Right(
					Tok(lexer.DOT),
					RawNumber(),
				),
			),
		),
		func(p Pair[string, OptionalValue[string]], span models.Span) ast.Expr {
			return &ast.NumberLiteral{
				Int:  p.First,
				Frac: p.Second.Or(""),
				NodeInfo: ast.NodeInfo{
					Span: span,
				},
			}
		},
	)
}
