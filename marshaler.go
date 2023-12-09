package resp

import "errors"

// Interface Marshaler2 declares that a type can be marshaled into RESP v2. This
// is the custom marshaler you can declare for a type.
//
// The returned value is expected to be a fully valid RESP statement including
// the [EOL] postfixes.
type Marshaler2 interface {
	MarshalRESP2() ([]byte, error)
}

// Interface Marshaler3 declares that a type can be marshaled into RESP v3. This
// is the custom marshaler you can declare for a type.
//
// The returned value is expected to be a fully valid RESP statement including
// the [EOL] postfixes.
type Marshaler3 interface {
	MarshalRESP3() ([]byte, error)
}

// Interface Marshaler declares that a type can be marshaled into RESP, for either
// version 2 or 3. It requires the version be submitted for which one the encoder
// desires.
//
// To make things a bit easier there is a convenience function [ProxyMarshaler]
// that will handle delegating to either the [Marshaler2] or [Marshaler3] interfaces
// for you.
//
// The returned value is expected to be a fully valid RESP statement including
// the [EOL] postfixes.
type Marshaler interface {
	MarshalRESP(version Version) ([]byte, error)
}

// ProxyMarshaler takes an object of any type (interface{}) and attempts to
// marshal it based on the provided RESP version identifier.
//
// If the desired version is 3, it will look for the [Marshaler3] interface and
// use it's marshaler. If unavailable, it will default to the [Marshaler2].
// If the [Marshaler2] interface does not exist either, then an error is returned.
//
// Additionally, an error is returned if the version identifier is invalid.
func ProxyMarshaler(obj any, version Version) ([]byte, error) {
	if !version.Valid() {
		return nil, errors.New("unsupported version")
	}

	if version == Version3 {
		if v3, ok := obj.(Marshaler3); ok {
			return v3.MarshalRESP3()
		} else if v2, ok := obj.(Marshaler2); ok {
			return v2.MarshalRESP2()
		}
		// For now we return an error
		return nil, errors.New("object does not implement either Marshaler3 or Marshaler2 interfaces")
	}

	if v2, ok := obj.(Marshaler2); ok {
		return v2.MarshalRESP2()
	}

	// For now we return an error
	return nil, errors.New("object does not implement Marshaler2 interface")
}
