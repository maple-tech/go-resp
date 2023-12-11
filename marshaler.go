package resp

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var reflectJSONMarshaler = reflect.TypeOf(new(json.Marshaler)).Elem()
var reflectStringer = reflect.TypeOf(new(fmt.Stringer)).Elem()

// Interface Marshaler2 declares that a type can be marshaled into RESP v2. This
// is the custom marshaler you can declare for a type.
//
// The returned value is expected to be a fully valid RESP statement including
// the [EOL] postfixes.
type Marshaler2 interface {
	MarshalRESP2() ([]byte, error)
}

var reflectMarshaler2 = reflect.TypeOf(new(Marshaler2)).Elem()

// Interface Marshaler3 declares that a type can be marshaled into RESP v3. This
// is the custom marshaler you can declare for a type.
//
// The returned value is expected to be a fully valid RESP statement including
// the [EOL] postfixes.
type Marshaler3 interface {
	MarshalRESP3() ([]byte, error)
}

var reflectMarshaler3 = reflect.TypeOf(new(Marshaler3)).Elem()

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

var reflectMarshaler = reflect.TypeOf(new(Marshaler)).Elem()

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

func Marshal2(obj any) ([]byte, error) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	if typ.Kind() == reflect.Pointer {
		if val.IsNil() {
			return NewNull().MarshalRESP3()
		}

		val = reflect.Indirect(val)
		typ = val.Type()
	}

	// Check the interface values first.
	if typ.Implements(reflectMarshaler2) {
		v, ok := val.Interface().(Marshaler2)
		if !ok {
			return nil, fmt.Errorf("failed to type assert value as resp.Marshaler2")
		}
		return v.MarshalRESP2()
	} else if typ.Implements(reflectMarshaler) {
		v, ok := val.Interface().(Marshaler)
		if !ok {
			return nil, fmt.Errorf("failed to type assert value as resp.Marshaler")
		}
		return v.MarshalRESP(2)
	}

	switch typ.Kind() {
	case reflect.Bool:
		ok := 0
		if val.Bool() {
			ok = 1
		}
		return NewInteger(int64(ok)).MarshalRESP2()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInteger(val.Int()).MarshalRESP2()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewInteger(int64(val.Uint())).MarshalRESP2()
	case reflect.Float32, reflect.Float64:
		str := strconv.FormatFloat(val.Float(), 'G', -1, 64)
		return NewSimpleString(str).MarshalRESP2()
	case reflect.Array, reflect.Slice:
		return marshalSlice(val)
	case reflect.Map, reflect.Struct:
		str, err := json.Marshal(val.Interface())
		if err != nil {
			return nil, fmt.Errorf("could not json encode slice for RESP2 compatibility: %s", err.Error())
		}
		return NewBulkString(string(str)).MarshalRESP2()
	default:
		if typ.Implements(reflectJSONMarshaler) {
			v, ok := val.Interface().(json.Marshaler)
			if !ok {
				return nil, fmt.Errorf("failed to type assert value as json.Marshaler")
			}
			byt, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to json encode for RESP2 compatibility: %s", err.Error())
			}
			return NewBulkString(string(byt)).MarshalRESP2()
		} else if typ.Implements(reflectStringer) {
			v, ok := val.Interface().(fmt.Stringer)
			if !ok {
				return nil, fmt.Errorf("failed to type assert value as fmt.Stringer")
			}
			return NewBulkString(v.String()).MarshalRESP2()
		}

		return nil, fmt.Errorf("unsupported type %s for RESP2 marshaling", typ.String())
	}
}

func Marshal3(obj any) ([]byte, error) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	if typ.Kind() == reflect.Pointer {
		if val.IsNil() {
			return NewNull().MarshalRESP3()
		}

		val = reflect.Indirect(val)
		typ = val.Type()
	}

	// Check the interface values first.
	if typ.Implements(reflectMarshaler2) {
		v, ok := val.Interface().(Marshaler2)
		if !ok {
			return nil, fmt.Errorf("failed to type assert value as resp.Marshaler2")
		}
		return v.MarshalRESP2()
	} else if typ.Implements(reflectMarshaler3) {
		v, ok := val.Interface().(Marshaler3)
		if !ok {
			return nil, fmt.Errorf("failed to type assert value as resp.Marshaler3")
		}
		return v.MarshalRESP3()
	} else if typ.Implements(reflectMarshaler) {
		v, ok := val.Interface().(Marshaler)
		if !ok {
			return nil, fmt.Errorf("failed to type assert value as resp.Marshaler")
		}
		return v.MarshalRESP(3)
	}

	switch typ.Kind() {
	case reflect.Bool:
		return NewBoolean(val.Bool()).MarshalRESP3()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInteger(val.Int()).MarshalRESP2()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return NewInteger(int64(val.Uint())).MarshalRESP2()
	case reflect.Float32, reflect.Float64:
		return NewDouble(val.Float()).MarshalRESP3()
	case reflect.Array, reflect.Slice:
		return marshalSlice(val)
	case reflect.Map:
		return marshalMap(val)
	case reflect.Struct:
		return marshalStruct(val)
	default:
		if typ.Implements(reflectJSONMarshaler) {
			v, ok := val.Interface().(json.Marshaler)
			if !ok {
				return nil, fmt.Errorf("failed to type assert value as json.Marshaler")
			}
			byt, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to json encode for interface compatibility: %s", err.Error())
			}
			return NewBulkString(string(byt)).MarshalRESP2()
		} else if typ.Implements(reflectStringer) {
			v, ok := val.Interface().(fmt.Stringer)
			if !ok {
				return nil, fmt.Errorf("failed to type assert value as fmt.Stringer")
			}
			return NewBulkString(v.String()).MarshalRESP3()
		}

		return nil, fmt.Errorf("unsupported type %s for RESP3 marshaling", typ.String())
	}
}

func Marshal(obj any, ver Version) ([]byte, error) {
	if ver == Version2 {
		return Marshal2(obj)
	} else if ver == Version3 {
		return Marshal3(obj)
	}
	return nil, errors.New("invalid version for resp.Marshal")
}

func marshalSlice(val reflect.Value) ([]byte, error) {
	raw := rawArray{entries: make([][]byte, val.Len())}
	var err error
	for i := 0; i < val.Len(); i++ {
		raw.entries[i], err = Marshal3(val.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal indice %d in array: %s", i, err.Error())
		}
	}
	return raw.MarshalRESP3()
}

func marshalMap(val reflect.Value) ([]byte, error) {
	raw := rawMap{entries: make(map[string][]byte, val.Len())}
	var err error

	iter := val.MapRange()
	var keyBytes []byte
	var valBytes []byte
	for iter.Next() {
		keyBytes, err = Marshal3(iter.Key().Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal map key: %s", err.Error())
		}

		valBytes, err = Marshal3(iter.Value().Interface())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal map value: %s", err.Error())
		}

		raw.entries[string(keyBytes)] = valBytes
	}

	return raw.Marshal3()
}

func marshalStruct(val reflect.Value) ([]byte, error) {
	typ := val.Type()

	raw := rawMap{entries: make(map[string][]byte, val.Len())}
	var err error

	var field reflect.StructField
	var tag string
	var ok bool
	var outName string
	for i := 0; i < typ.NumField(); i++ {
		field = typ.Field(i)
		if !field.IsExported() {
			continue
		}

		outName = typ.Name()

		if tag, ok = field.Tag.Lookup("resp"); ok {
			parts := strings.Split(tag, ",")
			if len(parts) == 1 {
				outName = parts[0]
			} else if len(parts) > 1 {
				return nil, fmt.Errorf("found struct tag '%s' that is not yet supported", tag)
			}
		}

		raw.entries[outName], err = Marshal3(val.Field(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("could not marshal struct field %s: %s", field.Name, err.Error())
		}
	}

	return raw.Marshal3()
}
