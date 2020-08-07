/*
FuncInfo

As information keeper for meta-data of functions.

  funcInfo := TypeExtBuilder(someFuncRef).FuncInfo()
  funcInfo.InAsTypes()
  funcInfo.OutAsTypes()

*/
package reflect

import (
	"reflect"
)

// Use the methods provided by this struct to access the meta-data.
type FuncInfo struct {
	inTypes []*TypeExt
	outTypes []*TypeExt
}

// Gets the types of in params(as "[]*TypeExt")
//
// See: "reflect.Type.In()"
func (self *FuncInfo) InTypes() []*TypeExt {
	newInTypes := make([]*TypeExt, len(self.inTypes))
	copy(newInTypes, self.inTypes)
	return newInTypes
}
// Gets the types of out params(as "[]*TypeExt")
//
// See: "reflect.Type.Out()"
func (self *FuncInfo) OutTypes() []*TypeExt {
	newOutTypes := make([]*TypeExt, len(self.outTypes))
	copy(newOutTypes, self.outTypes)
	return newOutTypes
}
// Gets the types of in params(as "[]reflect.Type")
//
// See: "reflect.Type.In()"
func (self *FuncInfo) InAsTypes() []reflect.Type {
	newInTypes := make([]reflect.Type, 0, len(self.inTypes))

	for _, t := range self.inTypes {
		newInTypes = append(newInTypes, t.AsType())
	}

	return newInTypes
}
// Gets the types of out params(as "[]*TypeExt")
//
// See: "reflect.Type.Out()"
func (self *FuncInfo) OutAsTypes() []reflect.Type {
	newOutTypes := make([]reflect.Type, 0, len(self.outTypes))

	for _, t := range self.outTypes {
		newOutTypes = append(newOutTypes, t.AsType())
	}

	return newOutTypes
}
