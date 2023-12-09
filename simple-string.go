package resp

import (
	"bytes"
	"errors"
)

// SimpleString implements the RESP2 Simple String type allowing for individual
// strings to be encoded/decoded. The specification dictates that the string not
// contain the terminator (\r\n) in the string, but the encoders make no validation
// check of this.
type SimpleString struct {
	byts []byte
}

// Value returns the inner value of this type without the additional protocol
// implementation. This is for using reflection and other marshalers.
func (s SimpleString) Value() any {
	return string(s.byts)
}

// Type returns the underlying RESP type identifier for this object
func (s SimpleString) Type() Type {
	return TypeSimpleString
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
func (s SimpleString) Contents() []byte {
	return s.byts
}

func (s *SimpleString) Unmarshal2(src []byte) error {
	if err := CanUnmarshalObject(src, s); err != nil {
		return err
	}

	s.byts = Contents(src)
	return nil
}

func (s *SimpleString) Unmarshal3(src []byte) error {
	return s.Unmarshal2(src)
}

func (s *SimpleString) Unmarshal(src []byte, _ Version) error {
	return s.Unmarshal2(src)
}

func (s SimpleString) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(s, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s SimpleString) Marshal3() ([]byte, error) {
	return s.Marshal2()
}

func (s SimpleString) Marshal(_ Version) ([]byte, error) {
	return s.Marshal2()
}

func NewSimpleString(str string) SimpleString {
	return SimpleString{[]byte(str)}
}

// ExtractSimpleString takes a byte slice that may be larger than an individual
// object and extracts the needed RESP data to fill a Simple String type. It will
// check the initial type identifier for you.
//
// It returns the object, the remaining bytes after the string and terminator,
// and an error if one occurred.
//
// If an error did happen, the object is returned as is, and the source is
// returned un-altered.
func ExtractSimpleString(src []byte) (SimpleString, []byte, error) {
	var s SimpleString
	if len(src) < 3 {
		return s, src, errors.New("cannot extract from empty source")
	}

	typ := Type(src[0])
	if typ != TypeSimpleString {
		return s, src, errors.New("attempted to extract simple string from incorrect type identifier")
	}

	term := bytes.Index(src, eol)
	if term == -1 {
		return s, src, errors.New("no terminator found for end of string")
	}

	s.byts = src[1:term]

	return s, src[term+len(eol):], nil
}
