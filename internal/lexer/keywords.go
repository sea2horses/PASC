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

	// Operators
	MOD

	// Flow Control
	SI
	SINO
	FINSI

	MIENTRAS
	HACER
	FINMIENTRAS
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
	"mientras":     MIENTRAS,
	"hacer":        HACER,
	"finmientras":  FINMIENTRAS,
	"mod":          MOD,
	"si":           SI,
	"sino":         SINO,
	"finsi":        FINSI,
}
