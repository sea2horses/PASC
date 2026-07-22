package tast

import "pseint-compiled/internal/semantic"

type Stmt interface{}

type TypedExpr interface {
	exprNode()
	Type() semantic.Type
}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) exprNode() {}

func (*StringLiteral) Type() semantic.Type {
	return semantic.PrimitiveType{Kind: semantic.STRING}
}

type NumberLiteral struct {
	Int  uint64
	Frac uint64
}

func (*NumberLiteral) exprNode() {}

/* ALL number literals are real by default */
/* TODO: Add warning to round any operation when assigning to an integer variable */
func (nl *NumberLiteral) Type() semantic.Type {
	return semantic.PrimitiveType{Kind: semantic.REAL}
}

type BooleanLiteral struct {
	Value bool
}

func (*BooleanLiteral) exprNode() {}

func (*BooleanLiteral) Type() semantic.Type {
	return semantic.PrimitiveType{Kind: semantic.BOOLEAN}
}

type VariableExpr struct {
	Symbol *semantic.Symbol
}

func (*VariableExpr) exprNode() {}

func (v *VariableExpr) Type() semantic.Type {
	return v.Symbol.Type
}

type Assignment struct {
	Target *semantic.Symbol
	Value  TypedExpr
}

type Cast struct {
	TargetType semantic.Type
	Expr       TypedExpr
}

func (c *Cast) Type() semantic.Type {
	return c.TargetType
}

type While struct {
	Condition TypedExpr
	Stmts     []Stmt
}

type MainFunction struct {
	Name  string
	Stmts []Stmt
}
