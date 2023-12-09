package resp

import "io"

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

// WriteTo is a generic utility for writing a type identifier, the contents, and
// the terminator to the given [io.Writer]. It returns the number of bytes
// written, and an error if one occurred while writing.
func WriteTo(obj Object, w io.Writer) (n int64, err error) {
	var l int
	if l, err = w.Write([]byte{byte(obj.Type())}); err != nil {
		return
	} else {
		n += int64(l)
	}
	if l, err = w.Write(obj.Contents()); err != nil {
		return
	} else {
		n += int64(l)
	}
	l, err = w.Write(eol)
	n += int64(l)
	return
}
