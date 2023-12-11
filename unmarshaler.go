package resp

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var reflectJSONUnmarshaler = reflect.TypeOf(new(json.Unmarshaler)).Elem()
var reflectTextUnmarshaler = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
var reflectNil = reflect.New(reflect.TypeOf(nil))

// Interface Unmarshaler2 declares that a type can be unmarshaled from RESP v2.
// This is the custom unmarshaler you can implement for a type.
//
// The source bytes are exactly as they are from RESP including the type
// identifying byte.
type Unmarshaler2 interface {
	UnmarshalRESP2(src []byte) error
}

var reflectUnmarshaler2 = reflect.TypeOf(new(Unmarshaler2)).Elem()

// Interface Unmarshaler3 declares that a type can be unmarshaled from RESP v3.
// This is the custom unmarshaler you can implement for a type.
//
// The source bytes are exactly as they are from RESP including the type
// identifying byte.
type Unmarshaler3 interface {
	UnmarshalRESP3(src []byte) error
}

var reflectUnmarshaler3 = reflect.TypeOf(new(Unmarshaler3)).Elem()

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

var reflectUnmarshaler = reflect.TypeOf(new(Unmarshaler)).Elem()

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

func Unmarshal2(src []byte, dst any) error {
	typ := reflect.TypeOf(dst)
	if typ.Kind() != reflect.Ptr {
		return errors.New("resp.Unmarshal2 requires a destination pointer")
	}
	val := reflect.ValueOf(dst)
	typ = val.Elem().Type()

	// Check the interface values first.
	if typ.Implements(reflectUnmarshaler2) {
		v, ok := val.Interface().(Unmarshaler2)
		if !ok {
			return fmt.Errorf("failed to type assert value as resp.Unmarshaler2")
		}
		return v.UnmarshalRESP2(src)
	} else if typ.Implements(reflectUnmarshaler) {
		v, ok := val.Interface().(Unmarshaler)
		if !ok {
			return fmt.Errorf("failed to type assert value as resp.Unmarshaler")
		}
		return v.UnmarshalRESP(src, 2)
	}

	obj, rest, err := Extract(src)
	if err != nil {
		return fmt.Errorf("could not extract RESP2 data: %s", err.Error())
	} else if len(rest) > 0 {
		return fmt.Errorf("found %d extra bytes after unmarshaling", len(rest))
	}

	switch typ.Kind() {
	case reflect.Bool:
		if obj.Type() != TypeInteger && obj.Type() != TypeBoolean {
			return fmt.Errorf("unexpected RESP object %s for bool destination", obj.Type().String())
		}

		ok := false
		if i, _ := obj.(Integer); i.int64 > 0 {
			ok = true
		}
		val.SetBool(ok)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if obj.Type() != TypeInteger {
			return fmt.Errorf("unexpected RESP object %s for integer destination", obj.Type().String())
		}
		v, _ := obj.(Integer)
		val.SetInt(v.int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if obj.Type() != TypeInteger {
			return fmt.Errorf("unexpected RESP object %s for unsigned integer destination", obj.Type().String())
		}
		v, _ := obj.(Integer)
		val.SetUint(uint64(v.int64))
	case reflect.Float32, reflect.Float64:
		if obj.Type() == TypeSimpleString {
			str, _ := obj.(SimpleString)
			flt, err := strconv.ParseFloat(string(str.byts), 64)
			if err != nil {
				return fmt.Errorf("failed to parse expected float string for RESP2 float compatibility: %s", err.Error())
			}
			val.SetFloat(flt)
		} else if obj.Type() == TypeDouble {
			dbl, _ := obj.(Double)
			val.SetFloat(dbl.float64)
		}
		return fmt.Errorf("unexpected RESP object %s for float destination", obj.Type().String())
	case reflect.Slice:
		return unmarshalSlice(obj, dst)
	case reflect.Map, reflect.Struct:
		if obj.Type() != TypeBulkString {
			return fmt.Errorf("unexpected RESP object %s for RESP2 compatibility json unmarshaling", obj.Type().String())
		}
		str, _ := obj.(BulkString)

		return json.Unmarshal(str.byts, dst)
	default:
		if typ.Implements(reflectJSONUnmarshaler) {
			v, _ := dst.(json.Unmarshaler)
			if obj.Type() == TypeSimpleString {
				str, _ := obj.(SimpleString)
				return v.UnmarshalJSON(str.byts)
			} else if obj.Type() == TypeBulkString {
				str, _ := obj.(BulkString)
				return v.UnmarshalJSON(str.byts)
			} else if obj.Type() == TypeVerbatimString {
				str, _ := obj.(VerbatimString)
				return v.UnmarshalJSON(str.byts)
			} else {
				return fmt.Errorf("json unmarshalling for RESP2 compatibility requires the incoming object be a string type, instead found %s", obj.Type())
			}
		} else if typ.Implements(reflectTextUnmarshaler) {
			v, _ := dst.(encoding.TextUnmarshaler)
			if obj.Type() == TypeSimpleString {
				str, _ := obj.(SimpleString)
				return v.UnmarshalText(str.byts)
			} else if obj.Type() == TypeBulkString {
				str, _ := obj.(BulkString)
				return v.UnmarshalText(str.byts)
			} else if obj.Type() == TypeVerbatimString {
				str, _ := obj.(VerbatimString)
				return v.UnmarshalText(str.byts)
			} else {
				return fmt.Errorf("text unmarshalling for RESP2 compatibility requires the incoming object be a string type, instead found %s", obj.Type())
			}
		}

		return fmt.Errorf("unsupported type %s for RESP2 unmarshaling", typ.String())
	}

	return nil
}

func Unmarshal3(src []byte, dst any) error {
	typ := reflect.TypeOf(dst)
	if typ.Kind() != reflect.Ptr {
		return errors.New("resp.Unmarshal2 requires a destination pointer")
	}
	val := reflect.ValueOf(dst)
	typ = val.Elem().Type()

	// Check the interface values first.
	if typ.Implements(reflectUnmarshaler3) {
		v, ok := val.Interface().(Unmarshaler3)
		if !ok {
			return fmt.Errorf("failed to type assert value as resp.Unmarshaler3")
		}
		return v.UnmarshalRESP3(src)
	} else if typ.Implements(reflectUnmarshaler2) {
		v, ok := val.Interface().(Unmarshaler2)
		if !ok {
			return fmt.Errorf("failed to type assert value as resp.Unmarshaler2")
		}
		return v.UnmarshalRESP2(src)
	} else if typ.Implements(reflectUnmarshaler) {
		v, ok := val.Interface().(Unmarshaler)
		if !ok {
			return fmt.Errorf("failed to type assert value as resp.Unmarshaler")
		}
		return v.UnmarshalRESP(src, 2)
	}

	obj, rest, err := Extract(src)
	if err != nil {
		return fmt.Errorf("could not extract RESP2 data: %s", err.Error())
	} else if len(rest) > 0 {
		return fmt.Errorf("found %d extra bytes after unmarshaling", len(rest))
	}

	switch typ.Kind() {
	case reflect.Bool:
		if obj.Type() != TypeBoolean {
			return fmt.Errorf("unexpected RESP object %s for bool destination", obj.Type().String())
		}
		v, _ := obj.(Boolean)
		val.SetBool(v.bool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if obj.Type() != TypeInteger {
			return fmt.Errorf("unexpected RESP object %s for integer destination", obj.Type().String())
		}
		v, _ := obj.(Integer)
		val.SetInt(v.int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if obj.Type() != TypeInteger {
			return fmt.Errorf("unexpected RESP object %s for unsigned integer destination", obj.Type().String())
		}
		v, _ := obj.(Integer)
		val.SetUint(uint64(v.int64))
	case reflect.Float32, reflect.Float64:
		if obj.Type() != TypeDouble {
			return fmt.Errorf("unexpected RESP object %s for float destination", obj.Type().String())
		}
		dbl, _ := obj.(Double)
		val.SetFloat(dbl.float64)
	case reflect.Slice:
		return unmarshalSlice(obj, dst)
	case reflect.Map:
		return unmarshalMap(obj, dst)
	case reflect.Struct:
		return unmarshalStruct(obj, dst)
	default:
		if typ.Implements(reflectJSONUnmarshaler) {
			v, _ := dst.(json.Unmarshaler)
			if obj.Type() == TypeSimpleString {
				str, _ := obj.(SimpleString)
				return v.UnmarshalJSON(str.byts)
			} else if obj.Type() == TypeBulkString {
				str, _ := obj.(BulkString)
				return v.UnmarshalJSON(str.byts)
			} else if obj.Type() == TypeVerbatimString {
				str, _ := obj.(VerbatimString)
				return v.UnmarshalJSON(str.byts)
			} else {
				return fmt.Errorf("json unmarshalling for RESP3 compatibility requires the incoming object be a string type, instead found %s", obj.Type())
			}
		} else if typ.Implements(reflectTextUnmarshaler) {
			v, _ := dst.(encoding.TextUnmarshaler)
			if obj.Type() == TypeSimpleString {
				str, _ := obj.(SimpleString)
				return v.UnmarshalText(str.byts)
			} else if obj.Type() == TypeBulkString {
				str, _ := obj.(BulkString)
				return v.UnmarshalText(str.byts)
			} else if obj.Type() == TypeVerbatimString {
				str, _ := obj.(VerbatimString)
				return v.UnmarshalText(str.byts)
			} else {
				return fmt.Errorf("text unmarshalling for RESP3 compatibility requires the incoming object be a string type, instead found %s", obj.Type())
			}
		}

		return fmt.Errorf("unsupported type %s for RESP3 unmarshaling", typ.String())
	}

	return nil
}

func Unmarshal(src []byte, dst any, ver Version) error {
	if ver == Version2 {
		return Unmarshal2(src, dst)
	} else if ver == Version3 {
		return Unmarshal3(src, dst)
	}

	return errors.New("unsupported version used for resp.Unmarshal")
}

func unmarshalSlice(obj Object, dst any) error {
	if obj.Type() != TypeArray {
		return fmt.Errorf("unexpected RESP object %s for slice destination", obj.Type())
	}
	arr, _ := obj.(Array)

	val := reflect.ValueOf(dst)

	if len(arr.entries) == 0 {
		val.SetLen(0)
	} else {
		arrType := arr.entries[0].Type()

		val.Set(reflect.MakeSlice(val.Type(), len(arr.entries), len(arr.entries)))
		for i, ent := range arr.entries {
			if ent.Type() != arrType {
				return fmt.Errorf("slices can only be unmarshaled using the same data types, detected %s instead of %s", ent.Type(), arrType)
			}

			if tmp, err := valueForObject(ent); err != nil {
				return fmt.Errorf("could not get value for object during slice unmarshaling: %s", err.Error())
			} else {
				if tmp.Type() != val.Type() {
					return fmt.Errorf("cannot apply type %s onto slice of %s", tmp.Type(), val.Type())
				}
				val.Index(i).Set(tmp)
			}
		}
	}

	return nil
}

func unmarshalMap(obj Object, dst any) error {
	if obj.Type() != TypeMap {
		return fmt.Errorf("unexpected RESP object %s for map destination", obj.Type())
	}
	mp, _ := obj.(Map)

	val := reflect.ValueOf(dst)
	val.Set(reflect.MakeMap(val.Type()))
	if len(mp.entries) == 0 {
		return nil
	}

	var err error
	var keyType, valType reflect.Type
	var keyVal, valVal reflect.Value
	for k, v := range mp.entries {
		keyVal, err = valueForObject(k)
		if err != nil {
			return fmt.Errorf("could not get value for key type %s: %s", k.Type(), err.Error())
		}

		if keyType != nil && reflect.TypeOf(keyVal) != keyType {
			return fmt.Errorf("maps can only be unmarshaled using the same key type, detected %s instead of %s", keyType, reflect.TypeOf(keyVal))
		}
		keyType = reflect.TypeOf(keyVal)

		valVal, err = valueForObject(v)
		if err != nil {
			return fmt.Errorf("could not get value for value type %s: %s", v.Type(), err.Error())
		}

		if valType != nil && reflect.TypeOf(valVal) != valType {
			return fmt.Errorf("maps can only be unmarshaled using the same value type, detected %s instead of %s", valType, reflect.TypeOf(valVal))
		}
		valType = reflect.TypeOf(valType)

		val.SetMapIndex(keyVal, valVal)
	}

	return nil
}

func unmarshalStruct(obj Object, dst any) error {
	if obj.Type() != TypeMap {
		return fmt.Errorf("unexpected RESP object %s for struct destination", obj.Type())
	}
	mp, _ := obj.(Map)

	val := reflect.ValueOf(dst)
	val.Set(reflect.New(val.Type()))
	if len(mp.entries) == 0 {
		return nil
	}

	meta := make(map[string]struct {
		name string
		typ  reflect.Type
	})
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		outName := field.Name

		if tag, ok := field.Tag.Lookup("resp"); ok {
			parts := strings.Split(tag, ",")
			if len(parts) == 1 {
				outName = parts[0]
			} else if len(parts) > 1 {
				return fmt.Errorf("found struct tag '%s' that is not yet supported", tag)
			}
		}

		meta[outName] = struct {
			name string
			typ  reflect.Type
		}{
			name: field.Name,
			typ:  field.Type,
		}
	}

	for k, v := range mp.entries {
		if k.Type() != TypeSimpleString {
			return fmt.Errorf("struct unmarshalling requires a map with simple string keys, instead found %s", k.Type())
		}
		keyStr, _ := k.(SimpleString)

		entry, ok := meta[string(keyStr.byts)]
		if !ok {
			return fmt.Errorf("could not find destination field %s for struct unmarshaling", keyStr.byts)
		}

		valVal, err := valueForObject(v)
		if err != nil {
			return fmt.Errorf("could not get value for value type %s: %s", v.Type(), err.Error())
		} else if valVal.Type() != entry.typ {
			return fmt.Errorf("mismatched types for struct field %s, received %s but trying to apply to %s", entry.name, valVal.Type(), entry.typ)
		}

		val.FieldByName(entry.name).Set(valVal)
	}

	return nil
}

func valueForObject(obj Object) (reflect.Value, error) {
	val, ok := obj.(Valuer)
	if !ok {
		return reflectNil, errors.New("expected Object to have Value() method")
	}

	return reflect.ValueOf(reflect.TypeOf(val.Value())), nil
}

// CanUnmarshalObject returns nil if the source data can be unmarshaled onto the
// destination [Object]. It does this by checking the first type identifier byte,
// and that it ends with the RESP terminator.
//
// It does not verify that the contents are correct, just that the type and
// general format is correct.
func CanUnmarshalObject(src []byte, dst Object) error {
	if len(src) <= 2 {
		return errors.New("source content is not long enough to be valid")
	} else if src[0] != byte(dst.Type()) {
		return errors.New("invalid type identifier for unmarshaling object")
	} else if !EndsWithTerminator(src) {
		return errors.New("source does not end with terminator")
	}
	return nil
}
