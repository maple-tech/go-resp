package resp

// WithoutTerminator returns the given byte slice with the last two characters
// removed since they are from the terminator (EOL)
func WithoutTerminator(byts []byte) []byte {
	return byts[:len(byts)-len(eol)]
}

// WithoutTypeIdentifier returns the given byte slice with the first byte skipped
// since that is the type identifier.
func WithoutTypeIdentifier(byts []byte) []byte {
	return byts[1:]
}

// Contents returns the given byte slice with the first byte, and last two bytes
// removed as those are the type identifier and the terminator (eol).
func Contents(byts []byte) []byte {
	return byts[1 : len(byts)-len(eol)]
}

// EndsWithTerminator returns true if the given bytes end with the desired
// terminator (eol)
func EndsWithTerminator(byts []byte) bool {
	ln := len(byts)
	tl := len(eol)
	return ln > 2 && byts[ln-tl] == eol[0] && byts[(ln-tl)+1] == eol[1]
}
