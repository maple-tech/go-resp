package resp

// Type declares the type identifier and therefore, the encoding type, of a RESP
// message.
type Type byte

const (
	TypeSimpleString Type = '+'
	TypeSimpleError  Type = '-'
	TypeInteger      Type = ':'
	TypeBulkString   Type = '$'
	TypeArray        Type = '*'

	TypeNull           Type = '_'
	TypeBoolean        Type = '#'
	TypeDouble         Type = ','
	TypeBigNumber      Type = '('
	TypeBulkError      Type = '!'
	TypeVerbatimString Type = '='
	TypeMap            Type = '%'
	TypeSet            Type = '~'
	TypePush           Type = '>'
)

var (
	// ValidTypesV2 declares all the type identifier bytes valid in the RESP v2
	// specification.
	ValidTypesV2 = []byte{
		byte(TypeSimpleString),
		byte(TypeSimpleError),
		byte(TypeInteger),
		byte(TypeBulkString),
		byte(TypeArray),
	}

	// ValidTypesV3 declares all the type identifier bytes that where added in
	// the v3 specification. See [ValidTypes] for the combined values.
	ValidTypesV3 = []byte{
		byte(TypeNull),
		byte(TypeBoolean),
		byte(TypeDouble),
		byte(TypeBigNumber),
		byte(TypeBulkError),
		byte(TypeVerbatimString),
		byte(TypeMap),
		byte(TypeSet),
		byte(TypePush),
	}

	// ValidTypes contains all the valid type identifier bytes for both version
	// 2 and 3 of the RESP specification.
	ValidTypes = append(ValidTypesV2, ValidTypesV3...)
)

func (t Type) String() string {
	switch t {
	case TypeSimpleString:
		return "Simple String"
	case TypeSimpleError:
		return "Simple Error"
	case TypeInteger:
		return "Integer"
	case TypeBulkString:
		return "Bulk String"
	case TypeArray:
		return "Array"

	case TypeNull:
		return "Null"
	case TypeBoolean:
		return "Boolean"
	case TypeDouble:
		return "Double"
	case TypeBigNumber:
		return "Big Number"
	case TypeBulkError:
		return "Bulk Error"
	case TypeVerbatimString:
		return "Verbatim String"
	case TypeMap:
		return "Map"
	case TypeSet:
		return "Set"
	case TypePush:
		return "Push"
	default:
		return "Unknown"
	}
}

// Valid returns true if this Type is in the combined v2 and v3 spec.
func (t Type) Valid() bool {
	byt := byte(t)
	for _, b := range ValidTypes {
		if byt == b {
			return true
		}
	}
	return false
}

// IsVersion2 returns true if the type identifier is under v2 spec.
func (t Type) IsVersion2() bool {
	byt := byte(t)
	for _, b := range ValidTypesV2 {
		if byt == b {
			return true
		}
	}
	return false
}

// IsVersion3 returns true if the type identifier is under v3 spec.
func (t Type) IsVersion3() bool {
	byt := byte(t)
	for _, b := range ValidTypesV3 {
		if byt == b {
			return true
		}
	}
	return false
}

// Version returns the Version identifier for which set this type identifier is
// in. It returns [VersionUnknown] instead of an error if it is in neither.
func (t Type) Version() Version {
	if t.IsVersion3() {
		return Version3
	} else if t.IsVersion2() {
		return Version2
	}
	return VersionUnknown
}

// Typer is anything that returns a Type value and is used for asserting RESP
// type payloads.
type Typer interface {
	Type() Type
}
