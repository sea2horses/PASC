package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) AnalyzeDimension(d *ast.Dimension) *tast.Dimension {
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
