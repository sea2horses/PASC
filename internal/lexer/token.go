package lexer;

type TokenType int8

/* TODO: Add character tokens like , = etc. */
const (
	// Literals
	KEYWORD TokenType = iota
	IDENTIFIER
	STRING_LITERAL
	NUMBER_LITERAL
	BOOLEAN_LITERAL
)

type Position struct {
	Line uint32
	Column uint32
}

// Token struct
type Token struct {
	Type TokenType /* Type of the token */
	Value string /* Value of the token */
}
