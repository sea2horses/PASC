package parser;

/* Node implemented by every AST Node */
type Node any;

type Expr interface {
	Node
}

type Stmt interface {
	Node
}

type StringLiteral struct {
	Content string
}

type Write struct {
	Print Expr
}

type MainFunction struct {
	Name string
	Stmts []Stmt
}
