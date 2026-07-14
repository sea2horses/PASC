package parser

import (
	"fmt"
	"pseint-compiled/internal/models"
	"pseint-compiled/internal/semantic"
	"strings"
)

type NodeInfo struct {
	span models.Span
}

func (ni *NodeInfo) Span() models.Span {
	return ni.span
}

/* Node implemented by every AST Node */
type Node interface {
	Span() models.Span
	String() string
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

type TypedExpr interface {
	Node
	Type() semantic.Type
}

type Stmt interface {
	Node
}

type StringLiteral struct {
	NodeInfo
	Content string
}

func (sl StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Content)
}

func (sl *StringLiteral) Accept(visitor Visitor) error {
	return visitor.VisitStringLiteral(sl)
}

type Write struct {
	NodeInfo
	Print Expr
}

func (w Write) String() string {
	return fmt.Sprintf("Write %s", w.Print)
}

func (w *Write) Accept(visitor Visitor) error {
	return visitor.VisitWrite(w)
}

type MainFunction struct {
	NodeInfo
	Name  string
	Stmts []Stmt
}

func (mf MainFunction) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("MainFunction '%s' {", mf.Name))
	for _, stmt := range mf.Stmts {
		builder.WriteString(stmt.String())
	}
	builder.WriteString("}")

	return builder.String()
}

func (mf *MainFunction) Accept(visitor Visitor) error {
	return visitor.VisitMainFunction(mf)
}
