package resp

import "errors"

// Interface Unmarshaler2 declares that a type can be unmarshaled from RESP v2.
// This is the custom unmarshaler you can implement for a type.
//
// The source bytes are exactly as they are from RESP including the type
// identifying byte.
type Unmarshaler2 interface {
	UnmarshalRESP2(src []byte) error
}

// Interface Unmarshaler3 declares that a type can be unmarshaled from RESP v3.
// This is the custom unmarshaler you can implement for a type.
//
// The source bytes are exactly as they are from RESP including the type
// identifying byte.
type Unmarshaler3 interface {
	UnmarshalRESP3(src []byte) error
}

// Interface Unmarshaler declares that a type can be unmarshaled from RESP.
// This is the custom unmarshaler you can implement for a type.
//
// Unlike [Unmarshaler3] or [Unmarshaler2] this version accepts a version
// identifier to choose the appropriate conversion method. You can use the helper
// function [ProxyUnmarshaler] to delegate to the [Unmarshaler3] or [Unmarshaler2]
// interfaces automatically.
//
// The source bytes are exactly as they are from RESP including the type
// identifying byte.
type Unmarshaler interface {
	UnmarshalRESP(src []byte, version Version) error
}

// ProxyUnmarshaler takes a source of RESP data, a target destination object, and
// a version identifier. With this, it will attempt to type-assert the destination
// and choose one of the [Unmarshaler3] or [Unmarshaler2] interfaces to unmarshal
// the given source RESP.
//
// If the desired RESP format is 3, it will look for [Unmarshaler3] to decode it.
// If that unmarshaler is not implement, it will default to [Unmarshaler2] for
// backwards compatibility.
//
// Note this does not use any built-in unmarshalers or stringer types. It will error
// if the unmarshaler interfaces are not implemented. Additionally, it errors out
// if the version identifier is invalid.
func ProxyUnmarshaler(src []byte, dst any, version Version) error {
	if !version.Valid() {
		return errors.New("unsupported version")
	}

	if version == Version3 {
		if v3, ok := dst.(Unmarshaler3); ok {
			return v3.UnmarshalRESP3(src)
		} else if v2, ok := dst.(Unmarshaler2); ok {
			return v2.UnmarshalRESP2(src)
		}

		// For now we return an error
		return errors.New("object does not implement either Unmarshaler3 or Unmarshaler2 interfaces")
	}

	if v2, ok := dst.(Unmarshaler2); ok {
		return v2.UnmarshalRESP2(src)
	}

	// For now we return an error
	return errors.New("object does not implement Marshaler2 interface")
}
