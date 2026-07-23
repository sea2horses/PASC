package analyzer

import (
	"fmt"
	"pseint-compiled/internal/ast"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/diagnostics"
	"pseint-compiled/internal/models"
	"pseint-compiled/internal/operators"
	"pseint-compiled/internal/semantic"
	"reflect"
	"slices"
)

type ErrIncorrectType struct {
	Expected semantic.Type
	Got      semantic.Type
}

func (e ErrIncorrectType) Error() string {
	return fmt.Sprintf("Incorrect type, expected %s, got: %s", e.Expected, e.Got)
}

type ErrUnOperationNotAllowed struct {
	Operation operators.OperatorType
	Type      semantic.Type
}

func (e ErrUnOperationNotAllowed) Error() string {
	return fmt.Sprintf("Unary operation %s is not allowed on type %s", e.Operation, e.Type)
}

type ErrBiOperationNotAllowed struct {
	Operation operators.OperatorType
	LHS       semantic.Type
	RHS       semantic.Type
}

func (e ErrBiOperationNotAllowed) Error() string {
	return fmt.Sprintf("Binary operation %s is not allowed on types %s & %s", e.Operation, e.LHS, e.RHS)
}

type ErrVariableDoesntExist struct {
	Name string
}

func (e ErrVariableDoesntExist) Error() string {
	return fmt.Sprintf("Variable '%s' doesn't exist", e.Name)
}

type ErrSymbolAlreadyExists struct {
	Name string
}

func (e ErrSymbolAlreadyExists) Error() string {
	return fmt.Sprintf("'%s' already exists!", e.Name)
}

type ErrUnsupportedType struct {
	Type semantic.Type
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("Unsupported Type: %s", e.Type)
}

type ErrExpectedType struct {
	Type semantic.Type
}

func (e ErrExpectedType) Error() string {
	return fmt.Sprintf("Expected Type: %s", e.Type)
}

type AnalyzerError string

func (pe AnalyzerError) Error() string {
	return string(pe)
}

const (
	ErrUnexpectedVoid = "Expression doesn't return anything"
	ErrLValueError    = "This expression is not assignable"
	ErrImmutable      = "Expression is immutable"
)

type Analyzer struct {
	scope       *semantic.Scope
	type_table  *semantic.TypeTable
	diagnostics []diagnostics.Diagnostic
}

func New() *Analyzer {
	/* Make base scope */
	return &Analyzer{scope: semantic.NewScope(nil)}
}

func (a *Analyzer) Analyze(program *ast.MainFunction) (*tast.MainFunction, []diagnostics.Diagnostic) {
	a.scope = semantic.NewScope(nil)
	a.type_table = semantic.NewTypeTable()
	a.diagnostics = nil

	a.scope.Declare(
		&semantic.Symbol{
			Name:     program.Name,
			Kind:     semantic.SymbolFunction,
			Type:     semantic.VoidType,
			Mutable:  false,
			Declared: program.Span,
		},
	)

	statements := a.analyze_statement_block(program.Stmts)

	return &tast.MainFunction{
		NodeInfo: program.NodeInfo,
		Name:     program.Name,
		Stmts:    statements,
	}, a.diagnostics
}

func (a *Analyzer) diagnosis(span models.Span, level diagnostics.DiagnosticLevel, format string, args ...any) {
	diagnostics.Dbg("Filing diagnosis of type: ", level, ". span: ", span)
	a.diagnostics = append(a.diagnostics, diagnostics.Diagnostic{
		Span:  span,
		Msg:   fmt.Sprintf(format, args...),
		Level: level,
	})
}

func (a *Analyzer) report(span models.Span, format string, args ...any) {
	a.diagnosis(span, diagnostics.ERROR, format, args...)
}

func (a *Analyzer) warn(span models.Span, format string, args ...any) {
	a.diagnosis(span, diagnostics.WARNING, format, args...)
}

func (a *Analyzer) info(span models.Span, format string, args ...any) {
	a.diagnosis(span, diagnostics.INFO, format, args...)
}

func (a *Analyzer) new_scope() {
	a.scope = semantic.NewScope(a.scope)
}

func (a *Analyzer) destroy_scope() {
	a.scope = a.scope.Parent()
}

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
	}

	return nil
}

func (a *Analyzer) analyze_type_ref(typeref *ast.TypeRef) *semantic.Type {
	return a.type_table.Get(typeref.Name)
}

func (a *Analyzer) analyze_declaration(decl *ast.Declaration) *tast.Declaration {
	diagnostics.Dbg("Analyzing declaration")
	/* Let's declare this bitch */
	/* Check the type exists! */
	t := a.analyze_type_ref(decl.Type)

	ok := true
	/* Let's add it to the symbol table */
	if _, exists := a.scope.Lookup(decl.Name); exists {
		a.report(decl.Span, "%s", ErrSymbolAlreadyExists{Name: decl.Name}.Error())
		ok = false
	}

	if t == nil {
		a.report(decl.Type.Span, "%s", "unknown type")
		ok = false
	}

	if !ok {
		return nil
	}

	/* Else, let's go */
	symbol := &semantic.Symbol{
		Name:     decl.Name,
		Kind:     semantic.SymbolVariable,
		Type:     *t,
		Declared: decl.Span,
		Mutable:  true,
	}

	a.scope.Declare(symbol)

	return &tast.Declaration{
		NodeInfo: decl.NodeInfo,
		Symbol:   symbol,
		Type:     *t,
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
		a.report(ass.Target.NodeSpan(), ErrLValueError)
		return nil
	}

	if !lvalue.IsMutable() {
		a.report(ass.NodeSpan(), ErrImmutable)
		a.info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
		return nil
	}

	/* Check coherence between target type and content type */
	if !semantic.EqualTypes(typed_expr.Type(), target.Type()) {
		var og_type semantic.Type = typed_expr.Type()
		typed_expr, ok = tryConvert(typed_expr, target.Type())
		if !ok {
			a.report(ass.NodeSpan(), "%s", ErrIncorrectType{Expected: target.Type(), Got: og_type}.Error())
			a.info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
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

var allowed_write_types []semantic.Type = []semantic.Type{
	semantic.IntegerType,
	semantic.BooleanType,
	semantic.RealType,
	semantic.StringType,
}

func (a *Analyzer) analyze_write(write *ast.Write) *tast.Write {
	diagnostics.Dbg("Analyzing write...")
	content := []tast.TypedExpr{}
	ok := true

	for _, exp := range write.Print {
		typed := a.analyze_expression(exp)

		if semantic.IsInvalid(typed.Type()) {
			ok = false
			continue
		}
		/* Write cannot have a void value */
		n := a.assert_nonvoid(typed)
		if !n {
			ok = false
			continue
		}
		/* Check that the type is supported */
		if !slices.Contains(allowed_write_types, typed.Type()) {
			a.report(typed.NodeSpan(), "%s", ErrUnsupportedType{Type: typed.Type()}.Error())
			ok = false
			continue
		}

		content = append(content, typed)
	}

	if !ok {
		return nil
	}

	return &tast.Write{NodeInfo: write.NodeInfo, Content: content}
}

func (a *Analyzer) analyze_read(read *ast.Read) []*tast.Read {
	diagnostics.Dbg("analyzing read")
	reads := []*tast.Read{}
	ok := true

	for _, exp := range read.Into {
		/* Check each expression */
		typed := a.analyze_expression(exp)

		if semantic.IsInvalid(typed.Type()) {
			ok = false
			continue
		}

		/* Assert it to lvalue */
		lvalue, good /* Im running out of names */ := typed.(tast.LValue)
		if !good {
			a.report(typed.NodeSpan(), ErrLValueError)
			ok = false
			continue
		}

		/* Check that it is mutable */
		if !lvalue.IsMutable() {
			a.report(typed.NodeSpan(), ErrImmutable)
			a.info(lvalue.AssignmentOrigin().Span, "%s", lvalue.AssignmentOrigin().Message)
			continue
		}

		/* All checks pass, add it */
		reads = append(reads, &tast.Read{NodeInfo: read.NodeInfo, Into: lvalue})
	}

	if !ok {
		return nil
	}

	return reads
}

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

func (a *Analyzer) analyze_expression(expr ast.Expr) tast.TypedExpr {
	diagnostics.Dbg("Analyzing expression")
	switch node := expr.(type) {
	case *ast.NumberLiteral:
		return &tast.NumberLiteral{
			NodeInfo: node.NodeInfo,
			Int:      node.Int,
			Frac:     node.Frac,
		}
	case *ast.StringLiteral:
		return &tast.StringLiteral{
			NodeInfo: node.NodeInfo,
			Value:    node.Content,
		}
	case *ast.BoolLiteral:
		return &tast.BooleanLiteral{
			NodeInfo: node.NodeInfo,
			Value:    node.Value,
		}
	case *ast.Variable:
		v := a.analyze_variable(node)
		return v
	case *ast.UnaryOperation:
		return a.analyze_unary_operation(node)
	case *ast.BinaryOperation:
		return a.analyze_binary_operation(node)
	}

	diagnostics.Dbg("Could not find the expression type")
	a.report(expr.NodeSpan(), "expression is not valid or implemented")
	return tast.ErrorExpr{NodeInfo: expr.Info()}
}

func (a *Analyzer) analyze_binary_operation(bo *ast.BinaryOperation) tast.TypedExpr {
	/* Analyze both inner expressions */
	LHS := a.analyze_expression(bo.LHS)
	RHS := a.analyze_expression(bo.RHS)
	/* Invalid by default */
	var resolved_type semantic.Type = semantic.InvalidType

	ok := a.assert_nonvoid(LHS, RHS)
	if !ok {
		return tast.ErrorExpr{
			NodeInfo: bo.NodeInfo,
		}
	}

	if !semantic.IsInvalid(LHS.Type()) && !semantic.IsInvalid(RHS.Type()) {
		/* Check if we can perform the operation */
		rule, ok := semantic.ResolveBinaryOperation(bo.Op.Type, LHS.Type(), RHS.Type())
		if ok {
			/* Check if we need to do a cast in LHS and RHS */
			if rule.LeftConversion != semantic.ConversionNone {
				LHS = applyConversion(LHS, rule.LeftConversion)
			}
			if rule.RightConversion != semantic.ConversionNone {
				RHS = applyConversion(RHS, rule.RightConversion)
			}
			resolved_type = rule.Result
		} else {
			a.report(bo.NodeInfo.Span, "%s", ErrBiOperationNotAllowed{LHS: LHS.Type(), RHS: RHS.Type(), Operation: bo.Op.Type}.Error())
			return tast.ErrorExpr{
				NodeInfo: bo.NodeInfo,
			}
		}
	}

	return &tast.BinaryOperation{
		NodeInfo:     bo.NodeInfo,
		LHS:          LHS,
		RHS:          RHS,
		Op:           a.analyze_operator(bo.Op),
		ResolvedType: resolved_type,
	}
}

func (a *Analyzer) analyze_unary_operation(uo *ast.UnaryOperation) tast.TypedExpr {
	/* Analyze the inner expression */
	typed_expr := a.analyze_expression(uo.Expr)
	/* Invalid by default */
	var resolved_type semantic.Type = semantic.InvalidType

	ok := a.assert_nonvoid(typed_expr)
	if !ok {
		return tast.ErrorExpr{
			NodeInfo: uo.NodeInfo,
		}
	}

	if !semantic.IsInvalid(typed_expr.Type()) {
		/* Check if we can perform the operation */
		rule, ok := semantic.ResolveUnaryOperation(uo.Op.Type, typed_expr.Type())
		if ok {
			/* Check if we got to do a cast */
			if rule.Conversion != semantic.ConversionNone {
				typed_expr = applyConversion(typed_expr, rule.Conversion)
			}
			/* First attach the resolved type */
			resolved_type = rule.Result
		} else {
			a.report(uo.NodeInfo.Span, "%s", ErrUnOperationNotAllowed{Type: typed_expr.Type(), Operation: uo.Op.Type}.Error())
			return tast.ErrorExpr{
				NodeInfo: uo.NodeInfo,
			}
		}
	}

	return &tast.UnaryOperation{
		NodeInfo:     uo.NodeInfo,
		Expr:         typed_expr,
		Op:           a.analyze_operator(uo.Op),
		ResolvedType: resolved_type,
	}
}

func (a *Analyzer) analyze_variable(variable *ast.Variable) tast.TypedExpr {
	diagnostics.Dbg("Analyzing variable: ", variable.Name)
	/* Get variable symbol */
	symbol, ok := a.scope.Lookup(variable.Name)

	if !ok || symbol == nil {
		diagnostics.Dbg("Not found!")
		a.report(variable.NodeSpan(), "%s", ErrVariableDoesntExist{Name: variable.Name}.Error())
		return &tast.ErrorExpr{NodeInfo: variable.NodeInfo}
	}

	diagnostics.Dbg("Found! Symbol: ", symbol)
	return &tast.VariableExpr{
		NodeInfo: variable.NodeInfo,
		Symbol:   symbol,
	}
}

func (a *Analyzer) analyze_operator(op ast.Operator) *tast.Operator {
	return &tast.Operator{
		NodeInfo: op.NodeInfo,
		Op:       op.Type,
	}
}

func (a *Analyzer) assert_type(target semantic.Type, exprs ...tast.TypedExpr) bool {
	ok := true
	for _, expr := range exprs {
		if expr.Type() != target {
			a.report(expr.NodeSpan(), "%s", ErrExpectedType{Type: target}.Error())
			ok = false
		}
	}
	return ok
}

func (a *Analyzer) assert_nontype(target semantic.Type, msg string, exprs ...tast.TypedExpr) bool {
	ok := true
	for _, expr := range exprs {
		if expr.Type() == target {
			a.report(expr.NodeSpan(), "%s", msg)
			ok = false
		}
	}
	return ok
}

func (a *Analyzer) assert_nonvoid(exprs ...tast.TypedExpr) bool {
	return a.assert_nontype(semantic.VoidType, ErrUnexpectedVoid, exprs...)
}

func tryConvert(src tast.TypedExpr, to semantic.Type) (*tast.Cast, bool) {
	conversion := semantic.CanCast(src.Type(), to)
	if conversion == nil {
		return nil, false
	}
	cast := applyConversion(src, *conversion)
	return cast, cast != nil
}

func applyConversion(operand tast.TypedExpr, conversion semantic.ConversionKind) *tast.Cast {
	lookup := semantic.SearchConversion(operand.Type(), conversion)
	if lookup == nil {
		return nil
	}
	return &tast.Cast{
		TargetType:     lookup.To,
		ConversionKind: &lookup.ConversionKind,
		Expr:           operand,
	}
}
