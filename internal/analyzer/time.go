package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) analyze_timeout(to *ast.TimeOut) *tast.TimeOut {
	/* Check expr is numeric */
	amount := a.analyze_expression(to.Amount)
	if !semantic.IsNumeric(amount.Type()) {
		a.Report(to.Amount.NodeSpan(), "amount must be numeric")
		return nil
	}

	return &tast.TimeOut{
		NodeInfo: to.NodeInfo,
		Amount:   amount,
		Unit:     to.Unit,
	}
}
