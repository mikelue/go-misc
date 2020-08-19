/*
Value Extension

"ValueExt" is the extension of "reflect.Value".

  // As value of int16(20)
  valueExt := ValueExtBuilder.NewByAny(int16(20))

ValueExtBuilder -
See provided methods in "IValueExtBuilder".
*/
package reflect

import (
	"fmt"
	"reflect"
)

// Collects functions used to construct "ValueExt", see "IValueExtBuilder".
const ValueExtBuilder IValueExtBuilder = 0

type IValueExtBuilder int

// Constructs "ValueExt" by "reflect.Value"
func (IValueExtBuilder) NewByValue(value reflect.Value) ValueExt {
	return ValueExt(value)
}

// Constructs "ValueExt" by any interface("interface{}")
func (self IValueExtBuilder) NewByAny(v interface{}) ValueExt {
	return self.NewByValue(reflect.ValueOf(v))
}

// Alias of reflect.Value, provides some convenient methods of features.
type ValueExt reflect.Value

// Converts this object to "reflect.Value"
func (self ValueExt) AsValue() reflect.Value {
	return reflect.Value(self)
}

// Converts this object to "interface{}"
func (self ValueExt) AsAny() interface{} {
	return reflect.Value(self).Interface()
}

// Gets the type of this value as "*TypeExt"
//
// See: "TypeExtBuilder.NewByType"
func (self ValueExt) TypeExt() *TypeExt {
	return TypeExtBuilder.NewByType(self.AsValue().Type())
}

// Gets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
//
// See: "ValueExt.SetFieldValue"
func (self ValueExt) GetFieldValue(tree ...string) ValueExt {
	currentValue := self.AsValue()

	for _, fieldName := range tree {
		if currentValue.Kind() == reflect.Ptr {
			currentValue = ValueExt(currentValue).RecursiveIndirect().AsValue()
		}

		if currentValue.Kind() != reflect.Struct {
			panic(fmt.Errorf("Field[%s] --> Current type is not struct: %v",
				fieldName, currentValue.Type(),
			))
		}

		currentValue = currentValue.FieldByName(fieldName)
		if !currentValue.IsValid() {
			panic(fmt.Errorf("Field[%s] is INVALID(!IsValid())", fieldName))
		}
	}

	return ValueExt(currentValue)
}

// Sets value of field, supporting tree visiting whether or not the value is
// struct or pointer to struct.
//
// Note: Only instance of addressable("reflect.Value.CanAddr") struct can be set value of field.
//
// returns: The original value
//
// See: "ValueExt.GetFieldValue"
func (self ValueExt) SetFieldValue(newValue ValueExt, tree ...string) ValueExt {
	targetValue := self.GetFieldValue(tree...).AsValue()

	if !targetValue.CanSet() {
		panic(fmt.Errorf("Value for type[%v] cannot be set", targetValue.Type()))
	}

	targetValue.Set(newValue.AsValue())
	return self
}

// Gets the value of struct(following points) represented by this value
func (self ValueExt) RecursiveIndirect() ValueExt {
	value := self.AsValue()

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return ValueExt(value)
}

// Returns true value if the value is array or slice
func (self ValueExt) IsArrayOrSlice() bool {
	switch reflect.Value(self).Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}

	return false
}

// Returns true value if the value is reflect.Ptr, reflect.Uintptr, or reflect.UnsafePointer
func (self ValueExt) IsPointer() bool {
	switch reflect.Value(self).Kind() {
	case reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
		return true
	}

	return false
}

// Checks if a value is viable
//
// 	For array, slice, map, chan: the value of "reflect.Value.Len()" must be > 0
//	For pointer, interface, or function: the value of "reflect.Value.IsNil()" must be false
//
//	Othewise: use reflect.Value.IsValid()
func (self ValueExt) IsViable() bool {
	reflectValue := reflect.Value(self)

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
