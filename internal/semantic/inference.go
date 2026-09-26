package semantic

import (
	"fmt"
)

type InferType struct {
	ID       int
	Resolved Type
}

func (i *InferType) Signature() TypeSignature {
	if i.Resolved == nil {
		return TypeSignature{
			BaseType: fmt.Sprintf("_infer T%d", i.ID),
		}
	} else {
		return i.Resolved.Signature()
	}
}

func (i *InferType) String() string {
	return i.Signature().String()
}

func InferTo(target Type)
