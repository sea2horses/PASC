package semantic

import (
	"pseint-compiled/internal/models"
	"strings"
)

type SymbolKind uint8

const (
	SymbolVariable SymbolKind = iota
	SymbolParameter
	SymbolFunction
)

type Symbol struct {
	Name     string
	Kind     SymbolKind
	Type     Type
	Declared models.Span
	Mutable  bool
}

type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: map[string]*Symbol{},
	}
}

/* TODO: Add option to remove case insensitive assignment */
func NormalizeName(name string) string {
	return strings.ToUpper(name)
}

func (s *Scope) Parent() *Scope {
	return s.parent
}

func (s *Scope) Declare(symbol *Symbol) bool {
	key := NormalizeName(symbol.Name)
	if _, exists := s.symbols[key]; exists {
		return false
	}

	s.symbols[key] = symbol
	return true
}

func (s *Scope) LookupLocal(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[NormalizeName(name)]
	return symbol, ok
}

func (s *Scope) Lookup(name string) (*Symbol, bool) {
	for current := s; current != nil; current = current.parent {
		if symbol, ok := current.LookupLocal(name); ok {
			return symbol, true
		}
	}

	return nil, false
}
