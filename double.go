package resp

import (
	"bytes"
	"errors"
	"strconv"
)

type Double struct {
	float64
}

func (d Double) Value() any {
	return d.float64
}

func (d Double) Type() Type {
	return TypeDouble
}

func (d Double) Contents() []byte {
	return []byte(strconv.FormatFloat(d.float64, 'G', -1, 64))
}

func (d *Double) Unmarshal3(src []byte) error {
	if err := CanUnmarshalObject(src, d); err != nil {
		return err
	}

	var err error
	d.float64, err = strconv.ParseFloat(string(Contents(src)), 64)
	return err
}

func (d *Double) Unmarshal(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("boolean is not available in RESP 2")
	}
	return d.Unmarshal3(src)
}

func (d Double) Marshal3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(d, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d Double) Marshal(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("double is not available in RESP 2")
	}
	return d.Marshal3()
}

func NewDouble(val float64) Double {
	return Double{val}
}

func ExtractDouble(src []byte) (Double, []byte, error) {
	var v Double

	term := bytes.Index(src, eol)
	if term == -1 {
		return v, src, errors.New("no terminator found for end of Double")
	}

	// Unmarshal checks the type and ending terminator for us
	err := v.Unmarshal3(src[:term+len(eol)])
	if err != nil {
		return v, src, err
	}

	return v, src[term+len(eol):], nil
}
