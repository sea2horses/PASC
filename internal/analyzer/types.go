package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) analyze_type_ref(typeref *ast.TypeRef) *semantic.Type {
	return a.type_table.Get(typeref.Name)
}

func (a *Analyzer) assert_type(target semantic.Type, exprs ...tast.TypedExpr) bool {
	ok := true
	for _, expr := range exprs {
		if expr.Type() != target {
			a.Report(expr.NodeSpan(), "%s", ErrExpectedType{Type: target}.Error())
			ok = false
		}
	}
	return ok
}

func (a *Analyzer) assert_nontype(target semantic.Type, msg string, exprs ...tast.TypedExpr) bool {
	ok := true
	for _, expr := range exprs {
		if expr.Type() == target {
			a.Report(expr.NodeSpan(), "%s", msg)
			ok = false
		}
	}
	return ok
}

func (a *Analyzer) assert_nonvoid(exprs ...tast.TypedExpr) bool {
	return a.assert_nontype(semantic.VoidType{}, ErrUnexpectedVoid, exprs...)
}
