package neoparser

import (
	"pseint-compiled/internal/ast"
	"pseint-compiled/internal/lexer"
)

func IfStmt() Pattern[ast.Stmt] {
	return Map(
		After(
			Kw(lexer.SI),
			Seq3(
				Locate(
					Left(
						Expression(),
						Optional(
							Kw(lexer.ENTONCES),
						),
					),
				),
				StatementBlock(),
				Left(
					Optional(
						ElseStmt(),
					),
					Seq2(
						Newline(),
						Kw(lexer.SINO),
					),
				),
			),
		),
		func(p Trio[Located[ast.Expr], []ast.Stmt, OptionalValue[*ast.Else]]) ast.Stmt {
			return &ast.If{
				NodeInfo: ast.NodeInfo{
					Span: p.First.Span,
				},
				Condition: p.First.Value,
				Stmts:     p.Second,
				Else:      p.Third.Or(nil),
			}
		},
	)
}

func ElseStmt() Pattern[*ast.Else] {
	return Map(
		Seq2(
			Locate(
				Kw(lexer.SINO),
			),
			StatementBlock(),
		),
		func(p Pair[Located[lexer.Keyword], []ast.Stmt]) *ast.Else {
			return &ast.Else{
				NodeInfo: ast.NodeInfo{
					Span: p.First.Span,
				},
				Stmts: p.Second,
			}
		},
	)
}
