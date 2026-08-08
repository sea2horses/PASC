package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/semantic"
)

func (a *Analyzer) analyze_if(i *ast.If) *tast.If {
	diagnostics.Dbg("Analyzing if...")
	/* Turn the expression into a typed expression */
	typed_condition := a.analyze_expression(i.Condition)
	stmts := a.analyze_statement_block(i.Stmts)
	/* Assert the condition type to be boolean */
	ok := a.assert_type(semantic.BooleanType, typed_condition)
	if !ok {
		return nil
	}

	var else_branch *tast.Else = nil

	if i.Else != nil {
		else_branch = a.analyze_else(i.Else)
	}

	return &tast.If{
		NodeInfo:  i.NodeInfo,
		Condition: typed_condition,
		Stmts:     stmts,
		Else:      else_branch,
	}
}

func (a *Analyzer) analyze_else(e *ast.Else) *tast.Else {
	diagnostics.Dbg("Analyzing else...")
	stmts := a.analyze_statement_block(e.Stmts)
	return &tast.Else{NodeInfo: e.NodeInfo, Stmts: stmts}
}

func (a *Analyzer) analyze_while(while *ast.While) *tast.While {
	diagnostics.Dbg("Analyzing while...")
	/* Analyze the condition */
	typed_condition := a.analyze_expression(while.Condition)
	stmts := a.analyze_statement_block(while.Stmts)
	/* Assert the condition type to be boolean */
	ok := a.assert_type(semantic.BooleanType, typed_condition)
	if !ok {
		return nil
	}

	return &tast.While{
		NodeInfo:  while.NodeInfo,
		Condition: typed_condition,
		Stmts:     stmts,
	}
}
