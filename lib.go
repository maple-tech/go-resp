// Package resp provides encoders, decoders, and utilities for formatting and
// reading RESP version 2 & 3 messages. This is the primary protocol used by
// Redis for communication, but it's also a neat encoding format.
package resp

// EOLString is a constant declaring the "end of line" for RESP.
const EOLString = "\r\n"

// eol is a constant declaring the "end of line" for RESP.
//
// Do not change this.
var eol = []byte{'\r', '\n'}

// EOL returns the constant for end-of-line in RESP. If you need a string there
// is [EOLString] instead.
//
// This is a function so we don't accidentally change the var that holds this.
func EOL() []byte {
	return eol
}
