package generator

import (
	"bytes"
	"fmt"
	tast "pseint-compiled/internal/ast_typed"
	"pseint-compiled/internal/semantic"
	"strings"
)

const (
	RUNTIME_WRITE_FUNCTION = "Write"
)

// ALL conversion kinds must be mapped to a runtime function
var conversionMap map[semantic.ConversionKind]string = map[semantic.ConversionKind]string{
	semantic.ConversionIntegerToReal:      "IntegerToReal",
	semantic.ConversionRealToIntegerExact: "ToIntegerExact",
}

const (
	WriteFn = "Write"
	ReadFn  = "Read"
)

type CodeGenerator struct {
	buffer  bytes.Buffer
	padding uint8
}

func (cg *CodeGenerator) up_scope() {
	cg.padding++
}

func (cg *CodeGenerator) down_scope() {
	cg.padding--
}

func (cg *CodeGenerator) write(msg string, args ...any) {
	// Generate final string
	s := fmt.Sprintf(msg, args...)
	// Pad it
	if cg.padding > 0 {
		target := "\n" + strings.Repeat("\t", int(cg.padding))
		s = strings.ReplaceAll(s, "\n", target)
	}
	// Append to buffer
	cg.buffer.WriteString(s)
}

func (cg *CodeGenerator) newline() {
	cg.buffer.WriteRune('\n')
	// Write a tab
	cg.buffer.Write(bytes.Repeat([]byte{'\t'}, int(cg.padding)))
}

func (cg *CodeGenerator) Generate(ast tast.Node) []byte {
	// Initialize a new empty buffer
	cg.buffer = *bytes.NewBuffer([]byte{})
	cg.write_node(ast)

	return cg.buffer.Bytes()
}

func (cg *CodeGenerator) write_node(node tast.Node) {
	switch n := node.(type) {
	case *tast.StringLiteral:
		cg.write("\"%s\"", n.Value)
	case *tast.NumberLiteral:
		cg.write("%d.%d", n.Int, n.Frac)
	case *tast.BooleanLiteral:
		if n.Value {
			cg.write("true")
		} else {
			cg.write("false")
		}
	case *tast.VariableExpr:
		cg.write("%s", n.Symbol.Name)
	case *tast.Cast:
		cg.write("%s(", conversionMap[*n.ConversionKind])
		cg.write_node(n.Expr)
		cg.write(")")
	case *tast.Operator:
		cg.write("%s", n.Op.String())
	case *tast.UnaryOperation:
		cg.write("(")
		cg.write_node(n.Op)
		cg.write(" ")
		cg.write_node(n.Expr)
		cg.write(")")
	case *tast.BinaryOperation:
		cg.write("(")
		cg.write_node(n.LHS)
		cg.write(" ")
		cg.write_node(n.Op)
		cg.write(" ")
		cg.write_node(n.RHS)
		cg.write(")")
	case *tast.Assignment:
		cg.write_node(n.Target)
		if n.Declarative {
			cg.write(" := ")
		} else {
			cg.write(" = ")
		}
		cg.write_node(n.Value)
	case *tast.Write:
		cg.write("%s(", WriteFn)
		for i, exp := range n.Content {
			if i != 0 {
				cg.write(", ")
			}
			cg.write_node(exp)
		}
		cg.write(")")
	case *tast.Read:
		cg.write("%s(&", ReadFn)
		cg.write_node(n.Into)
		cg.write(")")
	case *tast.If:
		cg.write("if ")
		cg.write_node(n.Condition)
		cg.write(" {")
		cg.write_statement_block(n.Stmts)
		cg.write("}")
		if n.Else != nil {
			cg.write(" else {")
			cg.write_statement_block(n.Else.Stmts)
			cg.write("}")
		}
	case *tast.While:
		cg.write("for ")
		cg.write("(")
		cg.write_node(n.Condition)
		cg.write(")")
		cg.write("{")
		cg.write_statement_block(n.Stmts)
		cg.write("}")
	case *tast.MainFunction:
		cg.write("/* Original Pseint Name: %s */\n", n.Name)
		cg.write("func main() {")
		cg.write_statement_block(n.Stmts)
		cg.write("}")
	}
}

func (cg *CodeGenerator) write_statement_block(stmts []tast.Stmt) {
	cg.up_scope()
	cg.write("\n")
	for i, stmt := range stmts {
		if i != 0 {
			cg.write("\n")
		}
		cg.write_node(stmt)
	}
	cg.down_scope()
	cg.write("\n")
}

// func (cg *CodeGenerator) newline() {
// 	cg.buffer.WriteRune('\n')
// 	// Write a tab
// 	cg.buffer.Write(bytes.Repeat([]byte{'\t'}, int(cg.padding)))
// }

// func (cg *CodeGenerator) write_node(node ast.Node, buffer *bytes.Buffer) {
// 	//
// 	switch n := node.(type) {
// 	case **ast.StringLiteral:
// 		buffer.WriteRune('"')
// 		buffer.WriteString(n.Content)
// 		buffer.WriteRune('"')
// 	case **ast.Write:
// 		buffer.WriteString(RUNTIME_WRITE_FUNCTION)
// 		buffer.WriteString("(")
// 		cg.write_node(n.Print, buffer)
// 		buffer.WriteString(")")
// 		cg.newline()
// 	case **ast.MainFunction:
// 		fmt.Fprintf(buffer, "/* Nombre original de la función: '%s' */", n.Name)
// 		cg.newline()
// 		buffer.WriteString("func main() {")
// 		cg.newline()
// 		cg.padding++
// 		for _, stmt := range n.Stmts {
// 			cg.write_node(stmt, buffer)
// 		}
// 		cg.padding--
// 		buffer.WriteString("}")
// 	}
// }
