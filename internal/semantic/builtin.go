package semantic

type BuiltinKind uint8

const (
	BuiltinSqrt BuiltinKind = iota
	BuiltinSin
	BuiltinCos
)

type BuiltinInfo struct {
	Kind BuiltinKind
}
