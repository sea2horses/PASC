package parser

import "pseint-compiled/internal/models"

type NodeInfo struct {
	span models.Span
}

func (ni *NodeInfo) Span() models.Span {
	return ni.span
}

/* Node implemented by every AST Node */
type Node interface {
	Span() models.Span
	Accept(visitor Visitor) error
}

type Visitor interface {
	VisitStringLiteral(node *StringLiteral) error
	VisitWrite(node *Write) error
	VisitMainFunction(node *MainFunction) error
}

type Expr interface {
	Node
}

type Stmt interface {
	Node
}

type StringLiteral struct {
	NodeInfo
	Content string
}

func (sl *StringLiteral) Accept(visitor Visitor) error {
	return visitor.VisitStringLiteral(sl)
}

type Write struct {
	NodeInfo
	Print Expr
}

func (w *Write) Accept(visitor Visitor) error {
	return visitor.VisitWrite(w)
}

type MainFunction struct {
	NodeInfo
	Name  string
	Stmts []Stmt
}

func (mf *MainFunction) Accept(visitor Visitor) error {
	return visitor.VisitMainFunction(mf)
}
