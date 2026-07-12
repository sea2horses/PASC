package lexer

import "strings"

//go:generate stringer -type=Keyword
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

func MapToKeyword(str []rune) (Keyword, bool) {
	val, ok := keywordMap[strings.ToLower(string(str))]
	return val, ok
}

var keywordMap map[string]Keyword = map[string]Keyword{
	"algoritmo":    ALGORITMO,
	"finalgoritmo": FINALGORITMO,
	"escribir":     ESCRIBIR,
	"leer":         LEER,
}
