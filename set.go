package resp

import (
	"errors"
	"fmt"
)

type Set struct {
	*Array
}

func (s Set) Type() Type {
	return TypeSet
}

func NewSet(entries ...Object) Set {
	return Set{&Array{entries}}
}

func ExtractSet(src []byte) (Set, []byte, error) {
	var v Set

	// Check the type
	ident := Type(src[0])
	if ident != TypeSet {
		return v, src, errors.New("expected Set type indicator for extracting")
	}

	// Use the same Array one
	arr, rest, err := ExtractArray(src)
	if err != nil {
		return v, src, fmt.Errorf("failed to extract set due to underlying array error: %s", err.Error())
	}
	v.Array = &arr

	return v, rest, err
}
