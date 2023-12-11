package resp

import (
	"bytes"
	"errors"
)

type VerbatimString struct {
	encoding []byte
	byts     []byte
}

func (v VerbatimString) Value() any {
	return string(v.byts)
}

func (v VerbatimString) Type() Type {
	return TypeVerbatimString
}

func (v VerbatimString) Contents() []byte {
	content := LenBytes(len(v.encoding) + len(v.byts) + 1)
	content = append(content, eol...)
	content = append(content, v.encoding...)
	content = append(content, ':')
	content = append(content, v.byts...)
	return content
}

func (v *VerbatimString) Unmarshal3(src []byte) error {
	if err := CanUnmarshalObject(src, v); err != nil {
		return err
	}

	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return errors.New("invalid verbatim string value, missing intermediate terminator")
	}

	// We don't care about the length here since we have the whole message

	// Extract the actual content
	content := src[interTerm+len(eol) : len(src)-len(eol)]

	colonInd := bytes.IndexByte(content, ':')
	if colonInd != 3 {
		return errors.New("malformed verbatim string value, encoding separator is in the wrong position")
	}

	v.encoding = content[:colonInd]
	v.byts = content[colonInd+1:]

	return nil
}

func (v *VerbatimString) Unmarshal(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("bulk error is not available in RESP 2")
	}
	return v.Unmarshal3(src)
}

func (v VerbatimString) Marshal3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(v, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewVerbatimString constructs a new "Verbatim String". This is much like a
// Bulk String, with the exception that it includes an encoding hint.
//
// Important! This encoding hint must be exactly 3 characters long.
func NewVerbatimString(encoding, content string) VerbatimString {
	return VerbatimString{
		encoding: []byte(encoding),
		byts:     []byte(content),
	}
}

func ExtractVerbatimString(src []byte) (VerbatimString, []byte, error) {
	var v VerbatimString

	term := IndexN(src, 2, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of VerbatimString")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.Unmarshal3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
