package semantic

type BuiltinKind uint8

const (
	BuiltinSqrt BuiltinKind = iota
	BuiltinSin
	BuiltinCos
	BuiltinTrunc
)

type BuiltinInfo struct {
	Kind BuiltinKind
}
