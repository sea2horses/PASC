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
	exp := a.analyze_expression(f.Start)
	variable, ok := exp.(*tast.VariableExpr)
	if !ok {
		a.Report(f.Start.NodeSpan(), "Expected variable")
	}

	start := a.analyze_expression(f.Start)
	end := a.analyze_expression(f.Until)

	targetType := variable.Type()

	start, ok = coerce_to(start, targetType)
	if !ok {
		a.Report(start.NodeSpan(), "Expected antoher typ")
	}
	end, ok = coerce_to(end, targetType)
	if !ok {
		a.Report(end.NodeSpan(), "Expected antoher typ")
	}

	var step tast.TypedExpr

	if f.Step != nil {
		step = a.analyze_expression(f.Step)
		step, ok = coerce_to(step, targetType)
		if !ok {
			a.Report(end.NodeSpan(), "Expected antoher typ")
		}
	} else {
		step = &tast.NumberLiteral{
			Int: "1",
		}

		step, ok = coerce_to(step, targetType)
		if !ok {
			a.Report(end.NodeSpan(), "Expected antoher typ")
		}
	}

	return &tast.For{
		Var:   variable,
		Start: start,
		Until: end,
		Step:  step,
		Stmts: a.analyze_statement_block(f.Stmts),
	}
}
