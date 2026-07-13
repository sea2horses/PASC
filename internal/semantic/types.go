package semantic;

type Type uint8

const (
	TypeInvalid Type = iota
	TypeInteger
	TypeReal
	TypeString
	TypeBoolean
	Unresolved
)

func IsNumeric(t Type) bool {
	return t == TypeInteger || t == TypeReal
}

