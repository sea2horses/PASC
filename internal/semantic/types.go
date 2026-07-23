package semantic

import "fmt"

type Type interface{ isType() }

type PrimKind uint8

const (
	INVALID PrimKind = iota
	VOID
	INTEGER // translates to int64
	REAL    // translates to float64
	STRING  // translates to string
	BOOLEAN // translates to bool
)

type PrimitiveType struct{ Kind PrimKind }
type ArrayType struct{ Elem Type }

func (PrimitiveType) isType() {}
func (ArrayType) isType()     {}

func (t PrimitiveType) String() string {
	switch t.Kind {
	case INTEGER:
		return "integer"
	case REAL:
		return "real"
	case STRING:
		return "string"
	case BOOLEAN:
		return "boolean"
	default:
		return "<invalid>"
	}
}

func (t ArrayType) String() string {
	return fmt.Sprintf("array[%s]", t.Elem)
}

var (
	InvalidType = PrimitiveType{Kind: INVALID}
	VoidType    = PrimitiveType{Kind: VOID}
	IntegerType = PrimitiveType{Kind: INTEGER}
	RealType    = PrimitiveType{Kind: REAL}
	StringType  = PrimitiveType{Kind: STRING}
	BooleanType = PrimitiveType{Kind: BOOLEAN}
)

func EqualTypes(a, b Type) bool {
	switch left := a.(type) {
	case PrimitiveType:
		right, ok := b.(PrimitiveType)
		return ok && left.Kind == right.Kind
	case ArrayType:
		right, ok := b.(ArrayType)
		return ok && EqualTypes(left.Elem, right.Elem)
	default:
		return false
	}
}

func IsInvalid(t Type) bool {
	prim, ok := t.(PrimitiveType)
	return ok && prim.Kind == INVALID
}

func IsNumeric(t Type) bool {
	prim, ok := t.(PrimitiveType)
	if !ok {
		return false
	}

	return prim.Kind == INTEGER || prim.Kind == REAL
}

func IsInteger(t Type) bool {
	return EqualTypes(t, IntegerType)
}

func IsReal(t Type) bool {
	return EqualTypes(t, RealType)
}

func IsString(t Type) bool {
	return EqualTypes(t, StringType)
}

func IsBoolean(t Type) bool {
	return EqualTypes(t, BooleanType)
}
