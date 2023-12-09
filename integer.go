package resp

import (
	"bytes"
	"errors"
	"strconv"
)

type Integer struct {
	int64
}

// Value returns the inner value of this type without the additional protocol
// implementation. This is for using reflection and other marshalers.
func (i Integer) Value() any {
	return i.int64
}

// Type returns the underlying RESP type identifier for this object
func (i Integer) Type() Type {
	return TypeInteger
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
func (i Integer) Contents() []byte {
	return []byte(strconv.FormatInt(i.int64, 10))
}

func (i *Integer) Unmarshal2(src []byte) error {
	if err := CanUnmarshalObject(src, i); err != nil {
		return err
	}

	var err error
	i.int64, err = strconv.ParseInt(string(Contents(src)), 10, 64)
	return err
}

func (i *Integer) Unmarshal3(src []byte) error {
	return i.Unmarshal2(src)
}

func (i *Integer) Unmarshal(src []byte, _ Version) error {
	return i.Unmarshal2(src)
}

func (i Integer) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(i, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i Integer) Marshal3() ([]byte, error) {
	return i.Marshal2()
}

func (i Integer) Marshal(_ Version) ([]byte, error) {
	return i.Marshal2()
}

func NewInteger(num int64) Integer {
	return Integer{num}
}

// ExtractInteger takes a byte slice that may be larger than an individual
// object and extracts the needed RESP data to fill a [Integer] type. It will
// check the initial type identifier for you.
//
// It returns the object, the remaining bytes after the error AND terminator,
// and an error if one occurred.
//
// If an error did happen, the object is returned as is, and the source is
// returned un-altered.
func ExtractInteger(src []byte) (Integer, []byte, error) {
	var v Integer

	term := bytes.Index(src, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of Integer")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.Unmarshal2(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
