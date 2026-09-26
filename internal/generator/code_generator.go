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

var typeMap map[semantic.Type]string = map[semantic.Type]string{
	semantic.IntegerType{}: "int64",
	semantic.StringType{}:  "string",
	semantic.BooleanType{}: "bool",
	semantic.RealType{}:    "float64",
}

const (
	WriteFn = "Write"
	ReadFn  = "Read"
	ClearFn = "ClearScreen"
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

/* TODO: Valdiation pass, don't generate code for statements with uninferred types */
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
	case *tast.Declaration:
		typerep, _ := typeMap[n.Symbol.Type]
		cg.write("var %s %s", n.Symbol.Name, typerep)
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
	case *tast.Case:
		cg.write("case ")
		cg.write_node(n.Value)
		cg.write(": {\n")
		cg.write_statement_block(n.Stmts)
		cg.write("}\n")
	case *tast.Switch:
		cg.write("switch ")
		cg.write_node(n.Base)
		cg.write(" {\n")
		for _, c := range n.Cases {
			cg.write_node(c)
		}
		if n.Default != nil {
			cg.write("default: {\n")
			cg.write_statement_block(n.Default)
			cg.write("}\n")
		}
		cg.write("}\n")
	case *tast.Dimension:
		arrayType := n.Symbol.Type.(semantic.ArrayType)
		elemType := semantic.Resolve(arrayType.Elem)
		typerep, _ := typeMap[elemType]
		cg.write("%s := NewTensor[%s]", n.Symbol.Name, typerep)
		cg.write("(")
		for i, dim := range n.Dimensions {
			if i > 0 {
				cg.write(", ")
			}
			cg.write("int(")
			cg.write_node(dim)
			cg.write(")")
		}
		cg.write(")")
	case *tast.Index:
		cg.write("*IndexTensor(")
		cg.write_node(n.Target)

		for _, index := range n.Indexes {
			cg.write(", int(")
			cg.write_node(index)
			cg.write(")")
		}

		cg.write(")")
	case *tast.Call:
		cg.generateCall(n)
	case *tast.ClearScreen:
		cg.write("%s()", ClearFn)
	case *tast.MainFunction:
		cg.write("/* Original Pseint Name: %s */\n", n.Name)
		cg.write("func main() {\n\tdefer writer.Flush()")
		cg.write_statement_block(n.Stmts)
		cg.write("}")
	}
}

func (cg *CodeGenerator) generateCall(call *tast.Call) {
	if variable, ok := call.Callable.(*tast.VariableExpr); ok {
		if variable.Symbol.Builtin != nil {
			cg.generateBuiltinCall(variable.Symbol.Builtin, call.Arguments)
			return
		}
	}

	cg.write_node(call.Callable)
	cg.write("(")

	for i, arg := range call.Arguments {
		if i > 0 {
			cg.write(", ")
		}

		cg.write_node(arg)
	}

	cg.write(")")
}

func (cg *CodeGenerator) generateBuiltinCall(
	builtin *semantic.BuiltinInfo,
	args []tast.TypedExpr,
) {
	switch builtin.Kind {
	case semantic.BuiltinSin:
		cg.write("math.Sin(")
		cg.write_node(args[0])
		cg.write(")")
	case semantic.BuiltinCos:
		cg.write("math.Cos(")
		cg.write_node(args[0])
		cg.write(")")
	case semantic.BuiltinTrunc:
		cg.write("math.Trunc(")
		cg.write_node(args[0])
		cg.write(")")
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
