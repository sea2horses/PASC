package lexer

import "pseint-compiled/internal/models"

//go:generate stringer -type=TokenType
type TokenType int8

/* TODO: Add character tokens like , = etc. */
const (
	// Literals
	KEYWORD TokenType = iota
	IDENTIFIER
	STRING_LITERAL
	NUMBER_LITERAL
	BOOLEAN_LITERAL

	// Characters
	DOT           // .
	COMMA         // ,
	EQUALS        // =
	L_PARENTHESES // (
	R_PARENTHESES // )

	// Binary Operators
	// Arithmetic
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	// Comparison
	L_ANGLE   // <
	R_ANGLE   // >
	AMPERSAND // &
	PIPE      // |
	EX_MARK   // !

	// Special
	EOF
)

// Token struct
type Token struct {
	Type  TokenType   /* Type of the token */
	Value string      /* Value of the token */
	Span  models.Span /* Span of the token */
}
