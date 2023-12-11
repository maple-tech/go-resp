package resp

import (
	"errors"
	"fmt"
)

type Push struct {
	*Array
}

func (p Push) Type() Type {
	return TypeSet
}

func NewPush(entries ...Object) Push {
	return Push{&Array{entries}}
}

func ExtractPush(src []byte) (Push, []byte, error) {
	var v Push

	// Check the type
	ident := Type(src[0])
	if ident != TypeSet {
		return v, src, errors.New("expected Push type indicator for extracting")
	}

	// Use the same Array one
	arr, rest, err := ExtractArray(src)
	if err != nil {
		return v, src, fmt.Errorf("failed to extract push due to underlying array error: %s", err.Error())
	}
	v.Array = &arr

	return v, rest, err
}
