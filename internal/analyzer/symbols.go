package analyzer

import (
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/semantic"
	"reflect"
)

func (a *Analyzer) analyze_variable(variable *ast.Variable) tast.TypedExpr {
	diagnostics.Dbg("Analyzing variable: ", variable.Name)
	/* Get variable symbol */
	symbol, ok := a.scope.Lookup(variable.Name)

	if !ok || symbol == nil {
		diagnostics.Dbg("Not found!")
		a.Report(variable.NodeSpan(), "%s", ErrVariableDoesntExist{Name: variable.Name}.Error())
		return &tast.ErrorExpr{NodeInfo: variable.NodeInfo}
	}

	diagnostics.Dbg("Found! Symbol: ", symbol)
	return &tast.VariableExpr{
		NodeInfo: variable.NodeInfo,
		Symbol:   symbol,
	}
}

func (a *Analyzer) analyze_declaration(decl *ast.Declaration) *tast.Declaration {
	diagnostics.Dbg("Analyzing declaration")
	/* Let's declare this bitch */
	/* Check the type exists! */
	t := a.analyze_type_ref(decl.Type)

	ok := true
	/* Let's add it to the symbol table */
	if _, exists := a.scope.Lookup(decl.Name); exists {
		a.Report(decl.Span, "%s", ErrSymbolAlreadyExists{Name: decl.Name}.Error())
		ok = false
	}

	if t == nil {
		a.Report(decl.Type.Span, "%s", "unknown type")
		ok = false
	}

	if !ok {
		return nil
	}

	/* Else, let's go */
	symbol := &semantic.Symbol{
		Name:     decl.Name,
		Kind:     semantic.SymbolVariable,
		Type:     t,
		Declared: decl.Span,
		Mutable:  true,
	}

	a.scope.Declare(symbol)

	return &tast.Declaration{
		NodeInfo: decl.NodeInfo,
		Symbol:   symbol,
		Type:     t,
	}
}

/* TODO: 'strict' assignment rule where implicit assignment is forbidden */
func (a *Analyzer) analyze_assignment(ass *ast.Assignment) *tast.Assignment {
	diagnostics.Dbg("Analyzing assignment")
	/* Check if LHS of assignment is a variable node only */
	variable, ok := ass.Target.(*ast.Variable)
	typed_expr := a.analyze_expression(ass.Content)
	nonvoid := a.assert_nonvoid(typed_expr)
	declarative := false

	diagnostics.Dbg("Target: ", ass.Target, reflect.TypeOf(ass.Target))

	if !nonvoid {
		diagnostics.Dbg("Void expression detected")
		return nil
	}

	diagnostics.Dbg("Passing through...")
	/* Let's check it out */
	if ok {
		diagnostics.Dbg("Direct variable assignment detected")
		if _, exists := a.scope.Lookup(variable.Name); !exists {
			diagnostics.Dbg("Making implicit declaration")
			/* Implicit declaration */
			a.scope.Declare(
				&semantic.Symbol{
					Name:     variable.Name,
					Kind:     semantic.SymbolVariable,
					Type:     typed_expr.Type(),
					Declared: ass.NodeSpan(),
					Mutable:  true,
				},
			)

			declarative = true
		}
	}
	target := a.analyze_expression(ass.Target)
	diagnostics.Dbg("Target <typed>: ", target, reflect.TypeOf(target))

	lvalue, ok := target.(tast.LValue)
	if !ok {
		a.Report(ass.Target.NodeSpan(), ErrLValueError)
		return nil
	}

	if !lvalue.IsMutable() {
		a.Report(ass.NodeSpan(), ErrImmutable)
		a.Info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
		return nil
	}

	/* Try coercing */
	_, ok = coerce_to(typed_expr, target.Type())
	if !ok {
		a.Report(ass.NodeSpan(), "%s", ErrIncorrectType{Expected: target.Type(), Got: typed_expr.Type()}.Error())
		a.Info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
		return nil
	}

	/* Check coherence between target type and content type */
	if !semantic.EqualTypes(typed_expr.Type(), target.Type()) {
		var og_type semantic.Type = typed_expr.Type()
		typed_expr, ok = tryConvert(typed_expr, target.Type())
		if !ok {
			a.Report(ass.NodeSpan(), "%s", ErrIncorrectType{Expected: target.Type(), Got: og_type}.Error())
			a.Info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
			return nil
		}
	}

	return &tast.Assignment{
		NodeInfo:    ass.NodeInfo,
		Target:      lvalue,
		Value:       typed_expr,
		Declarative: declarative,
	}
}
