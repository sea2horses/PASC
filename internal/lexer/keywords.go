package lexer

import "strings"

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

func mapToKeyword(str []rune) (Keyword, bool) {
	val, ok := keywordMap[strings.ToLower(string(str))]
	return val, ok
}

var keywordMap map[string]Keyword = map[string]Keyword{
	"algoritmo":    ALGORITMO,
	"finalgoritmo": FINALGORITMO,
	"escribir":     ESCRIBIR,
	"leer":         LEER,
}
