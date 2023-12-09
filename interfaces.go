package resp

// Valuer is any type that returns its inner type information as an interface{}
type Valuer interface {
	Value() any
}

// Object is anything in RESP that can be encoded/decoded by providing a type
// identifier and the binary contents of the statement.
type Object interface {
	Type() Type
	Contents() []byte
}
