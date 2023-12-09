package resp

import "strconv"

// Version declares the RESP protocol version as a byte for usage with constants
type Version byte

const (
	VersionUnknown Version = 0
	Version2       Version = 2
	Version3       Version = 3
)

// Valid returns true if the version identifier is considered "valid". This
// means it is within the accepted values of 2 and 3.
func (v Version) Valid() bool {
	return v >= Version2 && v <= Version3
}

// String implements the [Stringer] interface for converting the RESP version
// into a string.
//
// This is simply just a [strconv.FormatUint] call.
func (v Version) String() string {
	return strconv.FormatUint(uint64(v), 10)
}
