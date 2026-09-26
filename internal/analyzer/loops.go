package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) analyze_while(while *ast.While) *tast.While {
	diagnostics.Dbg("Analyzing while...")
	/* Analyze the condition */
	typed_condition := a.analyze_expression(while.Condition)
	stmts := a.analyze_statement_block(while.Stmts)
	/* Assert the condition type to be boolean */
	ok := a.assert_type(semantic.BooleanType{}, typed_condition)
	if !ok {
		return nil
	}

	return &tast.While{
		NodeInfo:  while.NodeInfo,
		Condition: typed_condition,
		Stmts:     stmts,
	}
}

func (a *Analyzer) analyze_for(f *ast.For) *tast.For {
	/* Check that start is a variable */
	asig := a.analyze_assignment(f.Start)

	/* Check that the step and end are numeric */
	until := a.analyze_expression(f.Until)

	if asig == nil || until == nil {
		return nil
	}

	var step tast.TypedExpr

	/* If step is not nil, then analyze it */
	if f.Step != nil {
		step = a.analyze_expression(f.Step)

		if step == nil {
			return nil
		}
	}

	/* Everything in order! */

	/* Check GOOD */
	variable := asig.Target.(*tast.VariableExpr)

	if variable == nil {
		a.Report(f.Start.Target.NodeSpan(), "expected variable")
	}

	if !semantic.IsNumeric(until.Type()) {
		a.Report(until.NodeSpan(), "expression must be numeric")
	}

	if step != nil && !semantic.IsNumeric(step.Type()) {
		a.Report(step.NodeSpan(), "expression must be numeric")
	}

	stmts := a.analyze_statement_block(f.Stmts)

	return &tast.For{
		NodeInfo: f.NodeInfo,
		Var:      variable,
		Start:    asig.Value,
		Until:    until,
		Step:     step,
		Stmts:    stmts,
	}
}
