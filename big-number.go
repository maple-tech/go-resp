package resp

import (
	"bytes"
	"errors"
	"math/big"
)

type BigNumber struct {
	big.Int
}

func (n BigNumber) Value() any {
	return n.Int
}

func (n BigNumber) Type() Type {
	return TypeBigNumber
}

func (n BigNumber) Contents() []byte {
	byts, _ := n.MarshalText()
	return byts
}

func (n *BigNumber) UnmarshalRESP3(src []byte) error {
	if err := CanUnmarshalObject(src, n); err != nil {
		return err
	}

	if err := n.UnmarshalText(Contents(src)); err != nil {
		return errors.New("could not unmarshal BigNumber value")
	}
	return nil
}

func (n *BigNumber) UnmarshalRESP(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("big number is not available in RESP 2")
	}
	return n.UnmarshalRESP3(src)
}

func (n BigNumber) MarshalRESP3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(n, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n BigNumber) MarshalRESP(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("big number is not available in RESP 2")
	}
	return n.MarshalRESP3()
}

func NewBigNumber(val big.Int) BigNumber {
	return BigNumber{val}
}

func ExtractBigNumber(src []byte) (BigNumber, []byte, error) {
	var v BigNumber

	term := bytes.Index(src, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of BigNumber")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.UnmarshalRESP3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
