package semantic

import "pseint-compiled/internal/operators"

type ConversionKind uint8

const (
	ConversionNone ConversionKind = iota
	ConversionIntegerToReal
	ConversionRealToIntegerExact
)

type ImplicitCast struct {
	From           Type
	To             Type
	ConversionKind ConversionKind
}

type BinaryOperationRule struct {
	Operator        operators.OperatorType
	Left            Type
	Right           Type
	Result          Type
	LeftConversion  ConversionKind
	RightConversion ConversionKind
}

type UnaryOperationRule struct {
	Operator   operators.OperatorType
	Operand    Type
	Result     Type
	Conversion ConversionKind
}

var AllowedCasts = []ImplicitCast{
	{From: IntegerType, To: RealType, ConversionKind: ConversionIntegerToReal},
	{From: RealType, To: IntegerType, ConversionKind: ConversionRealToIntegerExact},
}

var BinaryOperationTable = []BinaryOperationRule{
	// Addition and string concatenation.
	{operators.ADD, IntegerType, IntegerType, IntegerType, ConversionNone, ConversionNone},
	{operators.ADD, IntegerType, RealType, RealType, ConversionIntegerToReal, ConversionNone},
	{operators.ADD, RealType, IntegerType, RealType, ConversionNone, ConversionIntegerToReal},
	{operators.ADD, RealType, RealType, RealType, ConversionNone, ConversionNone},
	{operators.ADD, StringType, StringType, StringType, ConversionNone, ConversionNone},

	// Arithmetic.
	{operators.SUBTRACT, IntegerType, IntegerType, IntegerType, ConversionNone, ConversionNone},
	{operators.SUBTRACT, IntegerType, RealType, RealType, ConversionIntegerToReal, ConversionNone},
	{operators.SUBTRACT, RealType, IntegerType, RealType, ConversionNone, ConversionIntegerToReal},
	{operators.SUBTRACT, RealType, RealType, RealType, ConversionNone, ConversionNone},

	{operators.MULTIPLY, IntegerType, IntegerType, IntegerType, ConversionNone, ConversionNone},
	{operators.MULTIPLY, IntegerType, RealType, RealType, ConversionIntegerToReal, ConversionNone},
	{operators.MULTIPLY, RealType, IntegerType, RealType, ConversionNone, ConversionIntegerToReal},
	{operators.MULTIPLY, RealType, RealType, RealType, ConversionNone, ConversionNone},

	// Division and exponentiation always produce a real value.
	{operators.DIVIDE, IntegerType, IntegerType, RealType, ConversionIntegerToReal, ConversionIntegerToReal},
	{operators.DIVIDE, IntegerType, RealType, RealType, ConversionIntegerToReal, ConversionNone},
	{operators.DIVIDE, RealType, IntegerType, RealType, ConversionNone, ConversionIntegerToReal},
	{operators.DIVIDE, RealType, RealType, RealType, ConversionNone, ConversionNone},

	// Power will be shelved since it needs a function
	// {operators.POWER, IntegerType, IntegerType, RealType, ConversionIntegerToReal, ConversionIntegerToReal},
	// {operators.POWER, IntegerType, RealType, RealType, ConversionIntegerToReal, ConversionNone},
	// {operators.POWER, RealType, IntegerType, RealType, ConversionNone, ConversionIntegerToReal},
	// {operators.POWER, RealType, RealType, RealType, ConversionNone, ConversionNone},

	// MOD is strict and integer-only.
	{operators.MODULO, IntegerType, IntegerType, IntegerType, ConversionNone, ConversionNone},
	{operators.MODULO, IntegerType, RealType, IntegerType, ConversionNone, ConversionRealToIntegerExact},
	{operators.MODULO, RealType, IntegerType, RealType, ConversionRealToIntegerExact, ConversionNone},
	{operators.MODULO, RealType, RealType, IntegerType, ConversionRealToIntegerExact, ConversionRealToIntegerExact},

	// Numeric ordering.
	{operators.GREATER, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.GREATER, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.GREATER, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.GREATER, RealType, RealType, BooleanType, ConversionNone, ConversionNone},
	{operators.GREATER_EQ, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.GREATER_EQ, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.GREATER_EQ, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.GREATER_EQ, RealType, RealType, BooleanType, ConversionNone, ConversionNone},
	{operators.LESSER, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.LESSER, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.LESSER, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.LESSER, RealType, RealType, BooleanType, ConversionNone, ConversionNone},
	{operators.LESSER_EQ, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.LESSER_EQ, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.LESSER_EQ, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.LESSER_EQ, RealType, RealType, BooleanType, ConversionNone, ConversionNone},

	// Equality is intentionally unavailable for strings.
	{operators.EQUALS, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.EQUALS, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.EQUALS, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.EQUALS, RealType, RealType, BooleanType, ConversionNone, ConversionNone},
	{operators.EQUALS, BooleanType, BooleanType, BooleanType, ConversionNone, ConversionNone},
	{operators.NOT_EQUALS, IntegerType, IntegerType, BooleanType, ConversionNone, ConversionNone},
	{operators.NOT_EQUALS, IntegerType, RealType, BooleanType, ConversionIntegerToReal, ConversionNone},
	{operators.NOT_EQUALS, RealType, IntegerType, BooleanType, ConversionNone, ConversionIntegerToReal},
	{operators.NOT_EQUALS, RealType, RealType, BooleanType, ConversionNone, ConversionNone},
	{operators.NOT_EQUALS, BooleanType, BooleanType, BooleanType, ConversionNone, ConversionNone},

	// Logical operators.
	{operators.AND, BooleanType, BooleanType, BooleanType, ConversionNone, ConversionNone},
	{operators.OR, BooleanType, BooleanType, BooleanType, ConversionNone, ConversionNone},
}

var UnaryOperationTable = []UnaryOperationRule{
	{operators.SUBTRACT, IntegerType, IntegerType, ConversionNone},
	{operators.SUBTRACT, RealType, RealType, ConversionNone},
	{operators.NOT, BooleanType, BooleanType, ConversionNone},
}

func CanCast(from Type, to Type) *ConversionKind {
	for _, cast := range AllowedCasts {
		if EqualTypes(cast.From, from) && EqualTypes(cast.To, to) {
			return &cast.ConversionKind
		}
	}
	return nil
}

func SearchConversion(with Type, c ConversionKind) *ImplicitCast {
	for _, cast := range AllowedCasts {
		if EqualTypes(cast.From, with) && cast.ConversionKind == c {
			return &cast
		}
	}
	return nil
}

func ResolveBinaryOperation(op operators.OperatorType, left, right Type) (BinaryOperationRule, bool) {
	for _, rule := range BinaryOperationTable {
		if rule.Operator == op && EqualTypes(rule.Left, left) && EqualTypes(rule.Right, right) {
			return rule, true
		}
	}

	return BinaryOperationRule{}, false
}

func ResolveUnaryOperation(op operators.OperatorType, operand Type) (UnaryOperationRule, bool) {
	for _, rule := range UnaryOperationTable {
		if rule.Operator == op && EqualTypes(rule.Operand, operand) {
			return rule, true
		}
	}

	return UnaryOperationRule{}, false
}
