package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
)

func Write() Pattern[ast.Stmt] {
	return MapWithLocation(
		After(
			Kw(lexer.ESCRIBIR),
			Seq2(
				Optional(
					AsAny(
						Seq2(
							Kw(lexer.SIN),
							Kw(lexer.SALTAR),
						),
					),
				),
				SepBy1(
					Expression(),
					Tok(lexer.COMMA),
				),
			),
		),
		func(p Pair[OptionalValue[any], []ast.Expr], span models.Span) ast.Stmt {
			return &ast.Write{
				NodeInfo: ast.NodeInfo{
					Span: span,
				},
				Print:   p.Second,
				Newline: p.First.Some,
			}
		},
	)
}

func Read() Pattern[ast.Stmt] {
	return MapWithLocation(
		After(
			Kw(lexer.LEER),
			SepBy1(
				Expression(),
				Tok(lexer.COMMA),
			),
		),
		func(read []ast.Expr, span models.Span) ast.Stmt {
			return &ast.Read{
				NodeInfo: ast.NodeInfo{
					Span: span,
				},
				Into: read,
			}
		},
	)
}
