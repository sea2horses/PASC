package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
)

func Statement() Pattern[ast.Stmt] {

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

func Declaration() Pattern[*ast.Declaration] {
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
		func(p Located[Pair[[]string, *ast.TypeRef]]) *ast.Declaration {
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
