package resp

import (
	"bytes"
	"errors"
)

type BulkError struct {
	byts []byte
}

func (e BulkError) Value() any {
	return errors.New(string(e.byts))
}

func (e BulkError) Type() Type {
	return TypeBulkError
}

func (e BulkError) Contents() []byte {
	content := LenBytes(len(e.byts))
	content = append(content, eol...)
	content = append(content, e.byts...)
	return content
}

func (e *BulkError) Unmarshal3(src []byte) error {
	if err := CanUnmarshalObject(src, e); err != nil {
		return err
	}

	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return errors.New("invalid bulk error value, missing intermediate terminator")
	}

	// We don't care about the length here since we have the whole message

	// Extract the actual content
	e.byts = src[interTerm+len(eol) : len(src)-len(eol)]

	return nil
}

func (e *BulkError) Unmarshal(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("bulk error is not available in RESP 2")
	}
	return e.Unmarshal3(src)
}

func (e BulkError) Marshal3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(e, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e BulkError) Marshal(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("bulk error is not available in RESP 2")
	}
	return e.Marshal3()
}

func NewBulkError(str string) BulkError {
	return BulkError{[]byte(str)}
}

func ExtractBulkError(src []byte) (BulkError, []byte, error) {
	var v BulkError

	term := IndexN(src, 2, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of BulkError")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.Unmarshal3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
