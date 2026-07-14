package generator

import (
	"bytes"
	"fmt"
	"pseint-compiled/internal/parser"
)

type CodeGenerator struct {
	buffer bytes.Buffer
	padding uint8
}

func (cg *CodeGenerator) Generate(ast parser.Node) []byte {
	// Initialize a new empty buffer
	cg.buffer = *bytes.NewBuffer([]byte{})
	cg.write_node(ast, &cg.buffer)

	return cg.buffer.Bytes()
}

func (cg *CodeGenerator) newline() {
	cg.buffer.WriteRune('\n')
	// Write a tab
	cg.buffer.Write(bytes.Repeat([]byte{'\t'}, int(cg.padding)))
}

func (cg *CodeGenerator) write_node(node parser.Node, buffer *bytes.Buffer) {
	// 
	switch n := node.(type) {
	case *parser.StringLiteral:
		buffer.WriteRune('"')
		buffer.WriteString(n.Content)
		buffer.WriteRune('"')
	case *parser.Write:
		buffer.WriteString("Write(")
		cg.write_node(n.Print, buffer)
		buffer.WriteString(")")
		cg.newline()
	case *parser.MainFunction:
		fmt.Fprintf(buffer, "/* Nombre original de la función: '%s' */", n.Name)
		cg.newline()
		buffer.WriteString("func main() {")
		cg.newline()
		cg.padding++
		for _, stmt := range n.Stmts {
			cg.write_node(stmt, buffer)
		}
		cg.padding--
		buffer.WriteString("}")
	}
}
