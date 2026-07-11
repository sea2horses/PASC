package lexer;

type Keyword int8

const (
	// Main
	ALGORITMO Keyword = iota
	FINALGORITMO

	// i/o
	ESCRIBIR
	LEER

	// Flow Control
)

type TokenType int8

const (
	// Literals
	KEYWORD TokenType = iota
	IDENTIFIER
	STRING_LITERAL
	NUMBER_LITERAL
	BOOLEAN_LITERAL
)

// Token struct
type Token struct {
	Type TokenType /* Type of the token */
	Value string /* Value of the token */
}
