package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) analyze_dimension(d *ast.Dimension) *tast.Dimension {
	ok := true

	inf := a.type_table.MakeInference()

	var symbol *semantic.Symbol = &semantic.Symbol{
		Type:     semantic.ArrayType{Elem: inf, Rank: len(d.Dimensions)},
		Declared: d.NodeSpan(),
		Kind:     semantic.SymbolVariable,
		Mutable:  false,
	}

	/* Check that symbol doesn't exist */
	if symbol, ok := a.scope.Lookup(d.Name); ok {
		a.Report(d.NodeSpan(), "'%s' is already defined in this scope", d.Name)
		a.Info(symbol.Declared, "declared here")
		ok = false
	} else {
		a.scope.Declare(symbol)
	}

	var dims []tast.TypedExpr
	for _, dim := range d.Dimensions {
		/* Try to parse all expressions */
		expr := a.analyze_expression(dim)

		if expr == nil {
			ok = false
			continue
		}

		dims = append(dims, expr)
	}

	if !ok {
		return nil
	}

	return &tast.Dimension{
		NodeInfo:   d.NodeInfo,
		Symbol:     symbol,
		Dimensions: dims,
	}
}

func (a *Analyzer) analyze_indexing(d *ast.Index) *tast.Index {
	/* Right now you can only index arrays */

	/* Analyze expression to be indexed */
	typed_expr := a.analyze_expression(d.Target)
	if typed_expr == nil || semantic.IsInvalid(typed_expr.Type()) {
		return nil
	}

	lv, ok := typed_expr.(tast.LValue)
	if !ok {
		a.Report(d.NodeSpan(), "cannot index non-lvalue")
		return nil
	}

	if !semantic.IsArray(typed_expr.Type()) {
		a.Report(d.NodeSpan(), "cannot index non-array type '%s'", typed_expr.Type())
		return nil
	}

	underlying_type := (typed_expr.Type().(semantic.ArrayType)).Elem
	array_rank := (typed_expr.Type().(semantic.ArrayType)).Rank

	if len(d.Indexes) != array_rank {
		a.Report(d.NodeSpan(), "index count mismatch: expected %d, got %d", array_rank, len(d.Indexes))
		return nil
	}

	ok = true
	indexes := make([]tast.TypedExpr, len(d.Indexes))
	for _, index := range d.Indexes {
		curr := a.analyze_expression(index)
		if curr == nil || semantic.IsInvalid(curr.Type()) {
			ok = false
			continue
		}
		indexes = append(indexes, curr)
	}

	if !ok {
		return nil
	}

	return &tast.Index{
		NodeInfo:  d.NodeInfo,
		Target:    lv,
		Indexes:   indexes,
		IndexType: underlying_type,
	}
}
