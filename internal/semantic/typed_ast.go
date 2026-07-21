package semantic

type TypedExpr interface {
	exprNode()
	Type() Type
}

type StringLiteral struct {
	Value string
}

func (*StringLiteral) exprNode() {}

func (*StringLiteral) Type() Type {
	return PrimitiveType{Kind: STRING}
}

type NumberLiteral struct {
	Int  uint64
	Frac uint64
}

func (*NumberLiteral) exprNode() {}

func (nl *NumberLiteral) Type() Type {
	if nl.Frac == 0 {
		return PrimitiveType{Kind: INTEGER}
	} else {
		return PrimitiveType{Kind: REAL}
	}
}

type BooleanLiteral struct {
	Value bool
}

func (*BooleanLiteral) exprNode() {}

func (*BooleanLiteral) Type() Type {
	return PrimitiveType{Kind: BOOLEAN}
}

type VariableExpr struct {
	Symbol *Symbol
}

func (*VariableExpr) exprNode() {}

func (v *VariableExpr) Type() Type {
	return v.Symbol.Type
}

type Assignment struct {
	Target *Symbol
	Value  TypedExpr
}

type Cast struct {
	TargetType Type
	Expr       TypedExpr
}

func (c *Cast) Type() Type {
	return c.TargetType
}
