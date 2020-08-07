/*
Type Extension

"TypeExt" is the extension of "reflect.Type"

  // As type of "int32"
  valueExt := TypeExtBuilder.NewByAny(int32(0))

TypeExtBuilder -
See provided methods in "ITypeExtBuilder".
*/
package reflect

import (
	"fmt"
	"reflect"
)

// Collects functions used to construct "TypeExt", see "ITypeExtBuilder".
const TypeExtBuilder ITypeExtBuilder = 0

type ITypeExtBuilder int
// Constructs "TypeExt" by "reflect.Type"
func (ITypeExtBuilder) NewByType(t reflect.Type) *TypeExt {
	return &TypeExt{ t }
}
// Constructs "TypeExt" by any interface("interface{}")
func (self ITypeExtBuilder) NewByAny(v interface{}) *TypeExt {
	return self.NewByType(reflect.TypeOf(v))
}

// Holding the actual instance of "reflect.Type",
// see methods provided by this struct for features.
type TypeExt struct {
	typeInstance reflect.Type
}

// New returns a Value representing a pointer to a new zero value for the specified type.
//
// See: "reflect.New"
func (self *TypeExt) NewAsPointer() ValueExt {
	return ValueExtBuilder.NewByValue(reflect.New(self.AsType()))
}

// New returns a Value representing a instance for target type.
//
// The type is indirected to non-pointer type recursively.
//
// See: "reflect.New"
func (self *TypeExt) NewAsValue() ValueExt {
	return self.RecursiveIndirect().NewAsPointer().
		RecursiveIndirect()
}

// Gets the instance as "reflect.Type"
func (self *TypeExt) AsType() reflect.Type {
	return self.typeInstance
}

// As same as "reflect.Type.Kind()"
func (self *TypeExt) Kind() reflect.Kind {
	return self.AsType().Kind()
}

// Returns "true" if the type is "Reflect.Type"
func (self *TypeExt) IsReflectType() bool {
	return self.AsType().Implements(typeOfReflectType)
}

// Gets the type to which of a pointer refer
func (self *TypeExt) RecursiveIndirect() *TypeExt {
	finalType := self.AsType()

	for finalType.Kind() == reflect.Ptr {
		finalType = finalType.Elem()
	}

	return TypeExtBuilder.NewByType(finalType)
}

// In GoLang, you should use (*<Interface>)(nil) to
// reflect the type of an interface.
//
// This method uses the "reflect.Type.Elem()" to get the "real type" of an interface.
//
// returns: the type of target interface
func (self *TypeExt) InterfaceType() *TypeExt {
	originalType := self.AsType()

	if originalType.Kind() != reflect.Ptr {
		panic(fmt.Errorf("The original type should be pointer to \"interface{}\". Got: %v", originalType))
	}

	interfaceType := originalType.Elem()
	if interfaceType.Kind() != reflect.Interface {
		panic(fmt.Errorf("The kind of expected type is not \"reflect.Interface\". Got: %v", interfaceType))
	}

	return TypeExtBuilder.NewByType(interfaceType)
}

// Gets function information of this type
func (self *TypeExt) FuncInfo() *FuncInfo {
	if self.Kind() != reflect.Func {
		panic(fmt.Errorf("The kind of expected type is not \"reflect.Func\". Got: %v", self.AsType()))
	}

	targetType := self.AsType()

	inTypes := make([]*TypeExt, 0, targetType.NumIn())
	outTypes := make([]*TypeExt, 0, targetType.NumOut())

	for i := 0; i < targetType.NumIn(); i++ {
		inTypes = append(inTypes, TypeExtBuilder.NewByType(targetType.In(i)))
	}
	for i := 0; i < targetType.NumOut(); i++ {
		outTypes = append(outTypes, TypeExtBuilder.NewByType(targetType.Out(i)))
	}

	return &FuncInfo {
		inTypes: inTypes,
		outTypes: outTypes,
	}
}

var typeOfReflectType = TypeExtBuilder.NewByAny((*reflect.Type)(nil)).
	InterfaceType().AsType()
