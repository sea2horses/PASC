package semantic;

type Type interface { isType() }

type PrimKind uint8

const (
	INTEGER PrimKind = iota
	REAL
	STRING
	BOOLEAN
)

type PrimitiveType struct { Kind PrimKind }
type ArrayType struct { Elem Type }

func (PrimitiveType) isType() {}
func (ArrayType) isType() {}

func IsNumeric(t Type) bool {
	prim, ok := t.(PrimitiveType)
	if !ok {
		return false
	}

	return prim.Kind == INTEGER || prim.Kind == REAL
}
