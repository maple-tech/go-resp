package resp

import "errors"

// Extract takes a larger byte slice buffer and extracts the first full RESP
// object it can from it. It returns the RESP [Object], and the remaining bytes.
// If an error occurs then it returns nil, the original source, and the error.
func Extract(src []byte) (Object, []byte, error) {
	if len(src) == 0 {
		return nil, src, errors.New("cannot extract from empty source")
	}

	typ := Type(src[0])
	if !typ.Valid() {
		return nil, src, errors.New("invalid type identifier")
	}

	switch typ {
	case TypeSimpleString:
		return ExtractSimpleString(src)
	case TypeSimpleError:
		return ExtractSimpleError(src)
	case TypeInteger:
		return ExtractInteger(src)
	case TypeBulkString:
		return ExtractBulkString(src)
	}
	return nil, src, errors.New("could not extract valid RESP object")
}
