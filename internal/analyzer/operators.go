package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
)

func (a *Analyzer) analyze_operator(op ast.Operator) *tast.Operator {
	return &tast.Operator{
		NodeInfo: op.NodeInfo,
		Op:       op.Type,
	}
}
