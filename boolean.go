package resp

import (
	"bytes"
	"errors"
)

type Boolean struct {
	bool
}

func (b Boolean) Value() any {
	return b.bool
}

func (b Boolean) Type() Type {
	return TypeBoolean
}

func (b Boolean) Contents() []byte {
	if b.bool {
		return []byte{'t'}
	}
	return []byte{'f'}
}

func (b *Boolean) UnmarshalRESP3(src []byte) error {
	if err := CanUnmarshalObject(src, b); err != nil {
		return err
	}

	char := Contents(src)
	if char[0] == 't' {
		b.bool = true
	} else if char[0] == 'f' {
		b.bool = false
	} else {
		return errors.New("invalid value for RESP boolean")
	}
	return nil
}

func (b *Boolean) UnmarshalRESP(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("boolean is not available in RESP 2")
	}
	return b.UnmarshalRESP3(src)
}

func (b Boolean) MarshalRESP3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(b, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b Boolean) MarshalRESP(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("boolean is not available in RESP 2")
	}
	return b.MarshalRESP3()
}

func NewBoolean(ok bool) Boolean {
	return Boolean{ok}
}

func ExtractBoolean(src []byte) (Boolean, []byte, error) {
	var v Boolean

	term := bytes.Index(src, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of Boolean")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.UnmarshalRESP3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
