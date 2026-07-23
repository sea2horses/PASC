package generator

import (
	"bytes"
	"pseint-compiled/internal/ast"
)

const (
	RUNTIME_WRITE_FUNCTION = "Write"
)

type CodeGenerator struct {
	buffer  bytes.Buffer
	padding uint8
}

func (cg *CodeGenerator) Generate(ast ast.Node) []byte {
	// Initialize a new empty buffer
	cg.buffer = *bytes.NewBuffer([]byte{})

	return cg.buffer.Bytes()
}

// func (cg *CodeGenerator) newline() {
// 	cg.buffer.WriteRune('\n')
// 	// Write a tab
// 	cg.buffer.Write(bytes.Repeat([]byte{'\t'}, int(cg.padding)))
// }

// func (cg *CodeGenerator) write_node(node ast.Node, buffer *bytes.Buffer) {
// 	//
// 	switch n := node.(type) {
// 	case *ast.StringLiteral:
// 		buffer.WriteRune('"')
// 		buffer.WriteString(n.Content)
// 		buffer.WriteRune('"')
// 	case *ast.Write:
// 		buffer.WriteString(RUNTIME_WRITE_FUNCTION)
// 		buffer.WriteString("(")
// 		cg.write_node(n.Print, buffer)
// 		buffer.WriteString(")")
// 		cg.newline()
// 	case *ast.MainFunction:
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
