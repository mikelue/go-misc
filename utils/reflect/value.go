package reflect

import (
	"reflect"
)

func IsViable(v interface{}) bool {
	return ValueExt(reflect.ValueOf(v)).IsViable()
}

// Alias of reflect.Value, provides some convenient functions for programming on reflection.
type ValueExt reflect.Value

// Returns true value if the value is array or slice
func (v ValueExt) IsArray() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}

	return false
}

// Returns true value if the value is reflect.Ptr, reflect.Uintptr, or reflect.UnsafePointer
func (v ValueExt) IsPointer() bool {
	switch reflect.Value(v).Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
		return true
	}

	return false
}

// Checks if a value is viable
//
// 	For array, slice, map, chan: the value.Len() must be > 0
//
//	For pointer, interface, or function: the value.IsNil() must not be true
//
//	Othewise: use reflect.Value.IsValid()
func (v ValueExt) IsViable() bool {
	reflectValue := reflect.Value(v)

	switch reflectValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return reflectValue.Len() > 0
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Func:
		return !reflectValue.IsNil()
	default:
		return reflectValue.IsValid()
	}
}
