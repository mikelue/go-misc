package types

import (
	"reflect"
)

// Defines the interface to get "reflect.Type" instance for builtin types
//
// See: "BasicTypes", "PointerTypes", "SliceTypes"
type BuiltinTypeSpace interface {
	// Gets instance of "reflect.Type" relating to "int"
	OfInt() reflect.Type
	// Gets instance of "reflect.Type" relating to "int64"
	OfInt64() reflect.Type
	// Gets instance of "reflect.Type" relating to "int32"
	OfInt32() reflect.Type
	// Gets instance of "reflect.Type" relating to "int16"
	OfInt16() reflect.Type
	// Gets instance of "reflect.Type" relating to "int8"
	OfInt8() reflect.Type
	// Gets instance of "reflect.Type" relating to "uint"
	OfUint() reflect.Type
	// Gets instance of "reflect.Type" relating to "uint64"
	OfUint64() reflect.Type
	// Gets instance of "reflect.Type" relating to "uint32"
	OfUint32() reflect.Type
	// Gets instance of "reflect.Type" relating to "uint16"
	OfUint16() reflect.Type
	// Gets instance of "reflect.Type" relating to "uint8"
	OfUint8() reflect.Type
	// Gets instance of "reflect.Type" relating to "float32"
	OfFloat32() reflect.Type
	// Gets instance of "reflect.Type" relating to "float64"
	OfFloat64() reflect.Type
	// Gets instance of "reflect.Type" relating to "complex64"
	OfComplex64() reflect.Type
	// Gets instance of "reflect.Type" relating to "complex128"
	OfComplex128() reflect.Type
	// Gets instance of "reflect.Type" relating to "byte"
	OfByte() reflect.Type
	// Gets instance of "reflect.Type" relating to "bool"
	OfBool() reflect.Type
	// Gets instance of "reflect.Type" relating to "string"
	OfString() reflect.Type
}
