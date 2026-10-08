package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
	"pseint-compiled/internal/models"
)

func Statement() Pattern[ast.Stmt] {
	return OneOf(
		Declaration(),
	)
}

func StatementBlock() Pattern[[]ast.Stmt] {
	return SepBy(
		Statement(),
		Tok(lexer.NEWLINE),
	)
}

func MainFunction() Pattern[*ast.MainFunction] {
	return Map(
		Seq2(
			Locate(
				Right(
					Kw(lexer.ALGORITMO),
					Name(),
				),
			),
			Left(
				StatementBlock(),
				Seq2(
					Tok(lexer.NEWLINE),
					Kw(lexer.FINALGORITMO),
				),
			),
		),
		func(p Pair[Located[string], []ast.Stmt]) *ast.MainFunction {
			return &ast.MainFunction{
				NodeInfo: ast.NodeInfo{
					Span: p.First.Span,
				},
				Name:  p.First.Value,
				Stmts: p.Second,
			}
		},
	)
}

func Declaration() Pattern[ast.Stmt] {
	return Map(
		Locate(
			After(
				Kw(lexer.DEFINIR),
				Seq2(
					SepBy1(
						Name(),
						Tok(lexer.COMMA),
					),
					Right(
						Kw(lexer.COMO),
						Type(),
					),
				),
			),
		),
		func(p Located[Pair[[]string, *ast.TypeRef]]) ast.Stmt {
			return &ast.Declaration{
				NodeInfo: ast.NodeInfo{
					Span: p.Span,
				},
				Names: p.Value.First,
				Type:  p.Value.Second,
			}
		},
	)
}

func Assignment() Pattern[ast.Stmt] {
	return MapWithLocation(
		Seq2(
			Expression(),
			Right(
				OneOf(
					AsAny(Tok(lexer.EQUALS)),
					AsAny(
						Seq2(
							/* -> */
							Tok(lexer.MINUS),
							Tok(lexer.R_ANGLE),
						),
					),
				),
				Expression(),
			),
		),
		func(p Pair[ast.Expr, ast.Expr], span models.Span) ast.Stmt {
			return &ast.Assignment{
				NodeInfo: ast.NodeInfo{
					Span: span,
				},
				Target:  p.First,
				Content: p.Second,
			}
		},
	)
}

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
