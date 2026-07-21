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

type NumberLiteral struct {
	NodeInfo
	Int  uint64
	Frac uint64
}

func (nl NumberLiteral) String() string {
	return fmt.Sprintf("%d.%d", nl.Int, nl.Frac)
}

type BoolLiteral struct {
	NodeInfo
	Value bool
}

func (bl BoolLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

type Variable struct {
	NodeInfo
	Name string
}

func (v Variable) String() string {
	return v.Name
}

type Write struct {
	NodeInfo
	Print Expr
}

func (w Write) String() string {
	return fmt.Sprintf("Write %s", w.Print)
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
