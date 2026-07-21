package semantic;

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

type VariableExpr struct {
	Symbol *Symbol
}

func (*VariableExpr) exprNode() {}

func (v *VariableExpr) Type() Type {
	return v.Symbol.Type
}

type Assignment struct {
	Target *Symbol
	Value TypedExpr
}
