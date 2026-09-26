package semantic

import (
	"fmt"
	"pseint-compiled/internal/utils"
	"strings"
)

type TypeSignature struct {
	BaseType      string
	GenericParams []TypeSignature
	Metadata      string
}

func (t TypeSignature) String() string {
	var s strings.Builder

	s.WriteString(t.BaseType)
	if len(t.GenericParams) > 0 || t.Metadata != "" {
		params := strings.Join(utils.Map(t.GenericParams, func(gp TypeSignature) string {
			return gp.String()
		}), ", ")

		final := strings.Join([]string{params, t.Metadata}, ";")

		fmt.Fprintf(&s, "<%s>", final)
	}

	return s.String()
}

func (t TypeSignature) Equals(other TypeSignature) bool {
	if t.BaseType != other.BaseType {
		return false
	}
	if len(t.GenericParams) != len(other.GenericParams) {
		return false
	}
	for i := range t.GenericParams {
		if !t.GenericParams[i].Equals(other.GenericParams[i]) {
			return false
		}
	}
	return true
}

type Type interface{ Signature() TypeSignature }

type ArrayType struct {
	Elem Type
	Rank int
}

func (t ArrayType) String() string {
	return t.Signature().String()
}

func (t ArrayType) Signature() TypeSignature {
	return TypeSignature{
		BaseType: "array",
		GenericParams: []TypeSignature{
			t.Elem.Signature(),
		},
		Metadata: fmt.Sprintf("%d", t.Rank),
	}
}

/* Based on type signatures */
func EqualTypes(a, b Type) bool {
	return a.Signature().Equals(b.Signature())
}

/* Checks if two types are equal at the top level (array vs array) */
func EqualTopLevelType(a, b Type) bool {
	return a.Signature().BaseType == b.Signature().BaseType
}

func IsInvalid(t Type) bool {
	prim, ok := t.(PrimitiveType)
	return ok && prim.Kind() == INVALID
}

func IsVoid(t Type) bool {
	return EqualTypes(t, VoidType{})
}

func IsNumeric(t Type) bool {
	prim, ok := t.(PrimitiveType)
	if !ok {
		return false
	}

	return prim.Kind() == INTEGER || prim.Kind() == REAL
}

func IsInteger(t Type) bool {
	return EqualTypes(t, IntegerType{})
}

func IsReal(t Type) bool {
	return EqualTypes(t, RealType{})
}

func IsString(t Type) bool {
	return EqualTypes(t, StringType{})
}

func IsBoolean(t Type) bool {
	return EqualTypes(t, BooleanType{})
}

func IsGeneric(t Type) bool {
	return len(t.Signature().GenericParams) > 0
}

func IsArray(t Type) bool {
	return EqualTopLevelType(t, ArrayType{Elem: VoidType{}})
}
