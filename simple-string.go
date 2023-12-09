package resp

import (
	"bytes"
	"errors"
	"io"
)

// SimpleString implements the RESP2 Simple String type allowing for individual
// strings to be encoded/decoded. The specification dictates that the string not
// contain the terminator (\r\n) in the string, but the encoders make no validation
// check of this.
type SimpleString struct {
	Type
	byts []byte
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
func (s SimpleString) Contents() []byte {
	return s.byts
}

func (s *SimpleString) Unmarshal2(src []byte) error {
	if len(src) <= 2 {
		return errors.New("source content is not long enough to be valid")
	} else if src[0] != byte(s.Type) {
		return errors.New("invalid type identifier for unmarshaling SimpleString")
	} else if !EndsWithTerminator(src) {
		return errors.New("source does not end with terminator")
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

func (s SimpleString) WriteTo(w io.Writer) (n int64, err error) {
	var l int
	if l, err = w.Write([]byte{byte(s.Type)}); err != nil {
		return
	} else {
		n += int64(l)
	}
	if l, err = w.Write(s.Contents()); err != nil {
		return
	} else {
		n += int64(l)
	}
	l, err = w.Write(eol)
	n += int64(l)
	return
}

func (s SimpleString) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := s.WriteTo(&buf); err != nil {
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
	return SimpleString{
		TypeSimpleString,
		[]byte(str),
	}
}
