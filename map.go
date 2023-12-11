package resp

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type Map struct {
	entries map[Object]Object
}

// MapPair is an array of 2 [Object] mapping as `[Key, Value]` pair for a [Map]
// container.
type MapPair [2]Object

func (m Map) Value() any {
	return m.entries
}

func (m Map) Type() Type {
	return TypeMap
}

func (m Map) Contents() []byte {
	content := LenBytes(len(m.entries))
	for key, val := range m.entries {
		content = append(content, eol...)
		content = append(content, key.Contents()...)
		content = append(content, eol...)
		content = append(content, val.Contents()...)
	}
	return content
}

func (m *Map) UnmarshalRESP3(src []byte) error {
	if err := CanUnmarshalObject(src, m); err != nil {
		return err
	}

	// Split from the first terminator
	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return errors.New("invalid array value, missing intermediate terminator")
	}

	// Parse the length of the array
	ln, err := strconv.ParseUint(string(src[1:interTerm]), 10, 64)
	if err != nil {
		return err
	}

	// Establish the map
	m.entries = make(map[Object]Object, ln)

	// We now only care about the rest, but need to figure out what it is.
	rest := src[interTerm+len(eol) : len(src)-len(eol)]

	// Walk through the objects, expecting a key and value
	var key Object
	var val Object

	for i := 0; i < int(ln); i++ {
		key, rest, err = Extract(rest)
		if err != nil {
			return fmt.Errorf("failed to extract Map key: %s", err.Error())
		}

		val, rest, err = Extract(rest)
		if err != nil {
			return fmt.Errorf("failed to extract Map value: %s", err.Error())
		}

		m.entries[key] = val
	}

	return nil
}

func (m *Map) UnmarshalRESP(src []byte, ver Version) error {
	if ver == Version2 {
		return errors.New("map is not available in RESP 2")
	}
	return m.UnmarshalRESP3(src)
}

func (m Map) MarshalRESP3() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(m, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m Map) MarshalRESP(ver Version) ([]byte, error) {
	if ver == Version2 {
		return nil, errors.New("map is not available in RESP 2")
	}
	return m.MarshalRESP3()
}

func NewMap(entries ...MapPair) Map {
	v := Map{
		entries: make(map[Object]Object),
	}
	for _, p := range entries {
		v.entries[p[0]] = p[1]
	}
	return v
}

func ExtractMap(src []byte) (Map, []byte, error) {
	var v Map

	// Here is the issue. We don't actually know by any easy standards where the
	// end of a map is. We need to ingress the whole thing. Basically this
	// resembles the Array setup, but doubled for key/value pairs

	// Check the type
	ident := Type(src[0])
	if ident != TypeMap {
		return v, src, errors.New("expected Map type indicator for extracting")
	}

	// Split from the first terminator
	interTerm := bytes.Index(src, eol)
	if interTerm == -1 || interTerm == len(src)-len(eol) {
		return v, src, errors.New("invalid map value, missing intermediate terminator")
	}

	// Parse the length of the map, how many KV pairs
	ln, err := strconv.ParseUint(string(src[1:interTerm]), 10, 64)
	if err != nil {
		return v, src, err
	}

	// Establish the map
	v.entries = make(map[Object]Object, ln)

	// Start extracting, be aware this is recursive
	rest := src[interTerm+len(eol):]

	var key Object
	var val Object
	for i := 0; i < int(ln); i++ {
		key, rest, err = Extract(rest)
		if err != nil {
			return v, src, fmt.Errorf("error extracting map key at indice %d: %s", i, err.Error())
		}

		val, rest, err = Extract(rest)
		if err != nil {
			return v, src, fmt.Errorf("error extracting map value at indice %d: %s", i, err.Error())
		}

		v.entries[key] = val
	}

	return v, rest, nil
}

type rawMap struct {
	entries map[string][]byte
}

func (m rawMap) Value() any {
	return m.entries
}

func (m rawMap) Contents() []byte {
	content := LenBytes(len(m.entries))
	for key, val := range m.entries {
		content = append(content, eol...)
		content = append(content, key...)
		content = append(content, eol...)
		content = append(content, val...)
	}
	return content
}

func (m rawMap) Type() Type {
	return TypeArray
}

func (m rawMap) Marshal2() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := WriteTo(m, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m rawMap) Marshal3() ([]byte, error) {
	return m.Marshal2()
}

func (m rawMap) Marshal(_ Version) ([]byte, error) {
	return m.Marshal2()
}
