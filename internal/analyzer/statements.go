package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/diagnostics"
)

/* TODO: Global Scope Totality, PSEINT doesn't have scopes */
func (a *Analyzer) analyze_statement_block(stmts []ast.Stmt) []tast.Stmt {
	a.new_scope()
	statements := []tast.Stmt{}
	for _, statement := range stmts {
		typed := a.analyze_statement(statement)
		if typed != nil {
			statements = append(statements, typed...)
		}
	}
	a.destroy_scope()
	return statements
}

func (a *Analyzer) analyze_statement(stmt ast.Stmt) []tast.Stmt {
	diagnostics.Dbg("Analyzing statement...")
	switch s := stmt.(type) {
	case *ast.Assignment:
		return []tast.Stmt{a.analyze_assignment(s)}
	case *ast.Write:
		return []tast.Stmt{a.analyze_write(s)}
	case *ast.Read:
		reads := a.analyze_read(s)
		statements := make([]tast.Stmt, 0, len(reads))

		for _, read := range reads {
			statements = append(statements, read)
		}
		return statements
	case *ast.While:
		return []tast.Stmt{a.analyze_while(s)}
	case *ast.If:
		return []tast.Stmt{a.analyze_if(s)}
	case *ast.Declaration:
		return []tast.Stmt{a.analyze_declaration(s)}
	case *ast.ClearScreen:
		return []tast.Stmt{a.analyze_clear_screen(s)}
	case *ast.Switch:
		return []tast.Stmt{a.analyze_switch(s)}
	}

	a.Report(stmt.NodeSpan(), "extraneous statement")
	return nil
}
