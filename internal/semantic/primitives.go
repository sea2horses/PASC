package semantic

type PrimKind uint8

const (
	INVALID PrimKind = iota
	VOID
	INTEGER // translates to int64
	REAL    // translates to float64
	STRING  // translates to string
	BOOLEAN // translates to bool
)

type PrimitiveType interface {
	String() string
	Kind() PrimKind
	Signature() TypeSignature
}

type IntegerType struct{}

func (i IntegerType) Kind() PrimKind {
	return INTEGER
}

func (i IntegerType) String() string {
	return "integer"
}

func (i IntegerType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      i.String(),
		GenericParams: nil,
	}
}

type RealType struct{}

func (RealType) Kind() PrimKind {
	return REAL
}

func (r RealType) String() string {
	return "real"
}

func (r RealType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      r.String(),
		GenericParams: nil,
	}
}

type StringType struct{}

func (s StringType) Kind() PrimKind {
	return STRING
}

func (s StringType) String() string {
	return "string"
}

func (s StringType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      s.String(),
		GenericParams: nil,
	}
}

type BooleanType struct{}

func (b BooleanType) Kind() PrimKind {
	return BOOLEAN
}

func (b BooleanType) String() string {
	return "boolean"
}

func (b BooleanType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      b.String(),
		GenericParams: nil,
	}
}

type VoidType struct{}

func (v VoidType) Kind() PrimKind {
	return VOID
}

func (v VoidType) String() string {
	return "void"
}

func (v VoidType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      v.String(),
		GenericParams: nil,
	}
}

type InvalidType struct{}

func (i InvalidType) Kind() PrimKind {
	return INVALID
}

func (i InvalidType) String() string {
	return "<invalid>"
}

func (i InvalidType) Signature() TypeSignature {
	return TypeSignature{
		BaseType:      i.String(),
		GenericParams: nil,
	}
}
