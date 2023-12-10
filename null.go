package resp

import (
	"bytes"
	"errors"
)

// Null implements the RESP3 null type allowing for null (nothing) values to be
// encoded/decoded. The contents of this type are none. It is only the type
// identifier byte and the terminator.
type Null struct{}

func (n Null) Value() any {
	return nil
}

func (n Null) Type() Type {
	return TypeNull
}

func (n Null) Contents() []byte {
	return []byte{}
}

func (n *Null) Unmarshal3(src []byte) error {
	if err := CanUnmarshalObject(src, n); err != nil {
		return err
	}

	return nil
}

func (n *Null) Unmarshal(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("null is not available in RESP 2")
	}
	return n.Unmarshal3(src)
}

func (n Null) Marshal3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(n, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n Null) Marshal(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("null is not available in RESP 2")
	}
	return n.Marshal3()
}

func NewNull() Null {
	return Null{}
}

func ExtractNull(src []byte) (Null, []byte, error) {
	var v Null

	term := bytes.Index(src, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of Null")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.Unmarshal3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
