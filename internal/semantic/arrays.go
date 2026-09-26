package semantic

import "fmt"

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
