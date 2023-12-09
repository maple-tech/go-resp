package resp

import (
	"bytes"
	"errors"
)

// SimpleError implements the RESP2 Simple Error type allowing for individual
// strings to be encoded/decoded. The specification dictates that the string not
// contain the terminator (\r\n) in the string, but the encoders make no validation
// check of this.
type SimpleError struct {
	byts []byte
}

// Error implements the [Error] interface since this is actually an error type.
func (e SimpleError) Error() string {
	return string(e.byts)
}

// Value returns the inner value of this type without the additional protocol
// implementation. This is for using reflection and other marshalers.
func (e SimpleError) Value() any {
	return errors.New(string(e.byts))
}

// Type returns the underlying RESP type identifier for this object
func (e SimpleError) Type() Type {
	return TypeSimpleError
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
func (e SimpleError) Contents() []byte {
	return e.byts
}

func (e *SimpleError) Unmarshal2(src []byte) error {
	if err := CanUnmarshalObject(src, e); err != nil {
		return err
	}

	e.byts = Contents(src)
	return nil
}

func (e *SimpleError) Unmarshal3(src []byte) error {
	return e.Unmarshal2(src)
}

func (e *SimpleError) Unmarshal(src []byte, _ Version) error {
	return e.Unmarshal2(src)
}

func (e SimpleError) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(e, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e SimpleError) Marshal3() ([]byte, error) {
	return e.Marshal2()
}

func (e SimpleError) Marshal(_ Version) ([]byte, error) {
	return e.Marshal2()
}

func NewSimpleError(str string) SimpleError {
	return SimpleError{[]byte(str)}
}
