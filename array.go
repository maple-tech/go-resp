package resp

import (
	"bytes"
	"errors"
	"strconv"
)

// Array implements the RESP2 Array type allowing for a list of objects to be
// encoded/decoded. Note that the specification dictates that multiple types be
// allowed during his process. To facilitate this, this version of array extracts
// using the [Extract] utilities to consume the objects needed to fill the array.
//
// Inside the array is a slice of [Object] that you can use to further assert
// the types.
type Array struct {
	entries []Object
}

func (a Array) Value() any {
	return a.entries
}

func (a Array) Type() Type {
	return TypeArray
}

func (a Array) Contents() []byte {
	content := LenBytes(len(a.entries))
	for _, ent := range a.entries {
		content = append(content, eol...)
		content = append(content, ent.Contents()...)
	}
	return content
}

func (a *Array) Unmarshal2(src []byte) error {
	if err := CanUnmarshalObject(src, a); err != nil {
		return err
	}

	// Split from the first terminator
	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return errors.New("invalid array value, missing intermediate terminator")
	}

	// Parse the length of the array
	ln, err := strconv.ParseUint(string(src[1:interTerm]), 10, 64)
	if err != nil {
		return err
	}

	// Establish the array
	a.entries = make([]Object, ln)

	// We now only care about the rest, but need to figure out what it is.
	rest := src[interTerm+len(eol) : len(src)-len(eol)]

	// Walk through extracting all that we can
	var obj Object
	for i := range a.entries {
		obj, rest, err = Extract(rest)
		if err != nil {
			return err
		}
		a.entries[i] = obj
	}

	return nil
}

func (a *Array) Unmarshal3(src []byte) error {
	return a.Unmarshal2(src)
}

func (a *Array) Unmarshal(src []byte, _ Version) error {
	return a.Unmarshal2(src)
}

func (a Array) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(a, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (a Array) Marshal3() ([]byte, error) {
	return a.Marshal2()
}

func (a Array) Marshal(_ Version) ([]byte, error) {
	return a.Marshal2()
}

func NewArray(entries ...Object) Array {
	return Array{entries}
}
