package ast

import (
	"fmt"
	"pseint-compiled/internal/models"
	"pseint-compiled/internal/operators"
	"strings"
)

type NodeInfo struct {
	Span models.Span
}

func (ni NodeInfo) JoinSpan(span models.Span) {
	ni.Span = models.JoinSpans(ni.Span, span)
}

func (ni NodeInfo) Info() NodeInfo {
	return ni
}

func (ni NodeInfo) NodeSpan() models.Span {
	return ni.Span
}

/* Node implemented by every AST Node */
type Node interface {
	Info() NodeInfo
	NodeSpan() models.Span
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

type TypeRef struct {
	NodeInfo
	Name string
}

func (t TypeRef) String() string {
	return fmt.Sprintf("type<%s>", t.Name)
}

type Operator struct {
	NodeInfo
	Type operators.OperatorType
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

type Declaration struct {
	NodeInfo
	Name string
	Type *TypeRef
}

func (d Declaration) String() string {
	return fmt.Sprintf("decl %s as %s", d.Name, d.Type)
}

type Assignment struct {
	NodeInfo
	Target  Expr
	Content Expr
}

func (a Assignment) String() string {
	return fmt.Sprintf("%s = %s", a.Target, a.Content)
}

type Write struct {
	NodeInfo
	Print []Expr
}

func (w Write) String() string {
	return fmt.Sprintf("Write %s", w.Print)
}

type Read struct {
	NodeInfo
	Into []Expr
}

func (r Read) String() string {
	return fmt.Sprintf("Read %s", r.Into)
}

type If struct {
	NodeInfo
	Condition Expr
	Stmts     []Stmt
	Else      *Else
}

func (i If) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("If '%s' {\n", i.Condition))
	for _, stmt := range i.Stmts {
		builder.WriteString(stmt.String() + "\n")
	}
	builder.WriteString("} ")
	if i.Else != nil {
		builder.WriteString(i.Else.String())
	}

	return builder.String()
}

type Else struct {
	NodeInfo
	Stmts []Stmt
}

func (e Else) String() string {
	var builder strings.Builder

	builder.WriteString("Else {\n")
	for _, stmt := range e.Stmts {
		builder.WriteString(stmt.String() + "\n")
	}
	builder.WriteString("}")

	return builder.String()
}

type While struct {
	NodeInfo
	Condition Expr
	Stmts     []Stmt
}

func (w While) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("While '%s' {\n", w.Condition))
	for _, stmt := range w.Stmts {
		builder.WriteString(stmt.String() + "\n")
	}
	builder.WriteString("}")

	return builder.String()
}

type MainFunction struct {
	NodeInfo
	Name  string
	Stmts []Stmt
}

func (mf MainFunction) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("MainFunction '%s' {\n", mf.Name))
	for _, stmt := range mf.Stmts {
		builder.WriteString(stmt.String() + "\n")
	}
	builder.WriteString("}")

	return builder.String()
}
