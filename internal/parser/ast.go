package parser

import (
	"fmt"
	"pseint-compiled/internal/models"
	"strings"
)

type NodeInfo struct {
	span models.Span
}

func (ni NodeInfo) Span() models.Span {
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

type Operator struct {
	NodeInfo
	Type OperatorType
}

func (o Operator) String() string {
	return o.Type.String()
}

type UnaryOperation struct {
	NodeInfo
	Op   Operator
	Expr Expr
}

func (uo UnaryOperation) String() string {
	return fmt.Sprintf("(%s %s)", uo.Op, uo.Expr)
}

type BinaryOperation struct {
	NodeInfo
	LHS Expr
	Op  Operator
	RHS Expr
}

func (bo BinaryOperation) String() string {
	return fmt.Sprintf("(%s %s %s)", bo.LHS, bo.Op, bo.RHS)
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
