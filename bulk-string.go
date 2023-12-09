package resp

import (
	"bytes"
	"errors"
	"strconv"
)

// BulkString implements the RESP2 Bulk String type allowing for larger blocks
// strings or other binary data to be encoded/decoded. The specification dictates
// that the statement be in two parts. The first being the type identifier and
// a base-10 integer declaring the length, followed by a terminator. The next
// portion is the actual contents, followed by the final terminator.
type BulkString struct {
	byts []byte
}

func (s BulkString) Value() any {
	return string(s.byts)
}

func (s BulkString) Type() Type {
	return TypeBulkString
}

// Contents returns the inner contents of the item without its type identifier,
// or terminators applied. It is used as a quick way to get the actual body out-
// side of the RESP protocol specifics.
//
// Unlike simple types, these are split into two parts. The first is the length,
// and the second being the content. They are separated by a terminator.
func (s BulkString) Contents() []byte {
	slen := strconv.FormatUint(uint64(len(s.byts)), 10)

	content := make([]byte, len(slen))
	copy(content, slen)
	content = append(content, eol...)
	content = append(content, s.byts...)
	return content
}

// Unmarshal2 implements the [resp.Unmarshaler2] interface.
//
// NOTE: Bulk String messages are two part, as such this unmarshaler requires
// the intermediate terminators be left in place. Because we have the whole
// message, the first length portion is actually ignored and not validated.
func (s *BulkString) Unmarshal2(src []byte) error {
	if err := CanUnmarshalObject(src, s); err != nil {
		return err
	}

	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return errors.New("invalid bulk string value, missing intermediate terminator")
	}

	// We don't care about the length here since we have the whole message

	// Extract the actual content
	s.byts = src[interTerm+len(eol) : len(src)-len(eol)]

	return nil
}

func (s *BulkString) Unmarshal3(src []byte) error {
	return s.Unmarshal2(src)
}

func (s *BulkString) Unmarshal(src []byte, _ Version) error {
	return s.Unmarshal2(src)
}

func (s BulkString) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(s, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s BulkString) Marshal3() ([]byte, error) {
	return s.Marshal2()
}

func (s BulkString) Marshal(_ Version) ([]byte, error) {
	return s.Marshal2()
}

func NewBulkString(str string) BulkString {
	return BulkString{[]byte(str)}
}
