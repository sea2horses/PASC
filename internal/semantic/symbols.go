package semantic

import "pseint-compiled/internal/models"

type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota
	SymbolParameter
	SymbolFunction
)

type Symbol struct {
	Name string
	Kind SymbolKind
	Type Type
	Declared models.Span
	Mutable bool
}
