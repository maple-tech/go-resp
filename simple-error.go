package resp

import (
	"bytes"
	"errors"
	"io"
)

// SimpleError implements the RESP2 Simple Error type allowing for individual
// strings to be encoded/decoded. The specification dictates that the string not
// contain the terminator (\r\n) in the string, but the encoders make no validation
// check of this.
type SimpleError struct {
	Type
	byts []byte
}

// Error implements the [Error] interface since this is actually an error type.
func (e SimpleError) Error() string {
	return string(e.byts)
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
func (e SimpleError) Contents() []byte {
	return e.byts
}

func (e *SimpleError) Unmarshal2(src []byte) error {
	if len(src) <= 2 {
		return errors.New("source content is not long enough to be valid")
	} else if src[0] != byte(e.Type) {
		return errors.New("invalid type identifier for unmarshaling SimpleError")
	} else if !EndsWithTerminator(src) {
		return errors.New("source does not end with terminator")
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

func (e SimpleError) WriteTo(w io.Writer) (n int64, err error) {
	var l int
	if l, err = w.Write([]byte{byte(e.Type)}); err != nil {
		return
	} else {
		n += int64(l)
	}
	if l, err = w.Write(e.Contents()); err != nil {
		return
	} else {
		n += int64(l)
	}
	l, err = w.Write(eol)
	n += int64(l)
	return
}

func (e SimpleError) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := e.WriteTo(&buf); err != nil {
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
	return SimpleError{
		TypeSimpleError,
		[]byte(str),
	}
}
