/*
BasicTypes

Contains types over builting types of GoLang

  int64Type := BasicTypes.OfInt64()

PointerTypes

Contains types of pointer to builting types of GoLang

  strPtrType := PointerTypes.OfString()

SliceTypes

Contains types of slice for builting types of GoLang

  sliceUint8Type := SliceTypes.OfUint8()

ArrayTypes

Contains types of array for builting types of GoLang.

  array10IntType := ArrayTypes(10).OfInt()

The instances for this reflection is not cached by this package.

ErrorType

As "reflect.Type" for "error" interface.
*/
package types

import (
	"reflect"
)

// Instance of "reflect.Type" for "error"(interface)
var ErrorType reflect.Type = reflect.TypeOf((*error)(nil)).Elem()

// Contains instance of "reflect.Type" over builtin types of GoLang
//
// See: "reflect.TypeOf"
var BasicTypes BuiltinTypeSpace = &spaceImpl{
	instanceOfInt: reflect.TypeOf(int(0)),
	instanceOfInt64: reflect.TypeOf(int64(0)),
	instanceOfInt32: reflect.TypeOf(int32(0)),
	instanceOfInt16: reflect.TypeOf(int16(0)),
	instanceOfInt8: reflect.TypeOf(int8(0)),
	instanceOfUint: reflect.TypeOf(uint(0)),
	instanceOfUint64: reflect.TypeOf(uint64(0)),
	instanceOfUint32: reflect.TypeOf(uint32(0)),
	instanceOfUint16: reflect.TypeOf(uint16(0)),
	instanceOfUint8: reflect.TypeOf(uint8(0)),

	instanceOfFloat32: reflect.TypeOf(float32(0)),
	instanceOfFloat64: reflect.TypeOf(float64(0)),

	instanceOfComplex64: reflect.TypeOf(complex64(0)),
	instanceOfComplex128: reflect.TypeOf(complex128(0)),

	instanceOfByte: reflect.TypeOf(byte(0)),
	instanceOfBool: reflect.TypeOf(true),
	instanceOfString: reflect.TypeOf(""),
}

// Contains instances of "reflect.Type" over pointer of builtin types of GoLang
//
// See: "reflect.PtrTo"
var PointerTypes BuiltinTypeSpace = &spaceImpl {
	instanceOfInt: reflect.PtrTo(BasicTypes.OfInt()),
	instanceOfInt64: reflect.PtrTo(BasicTypes.OfInt64()),
	instanceOfInt32: reflect.PtrTo(BasicTypes.OfInt32()),
	instanceOfInt16: reflect.PtrTo(BasicTypes.OfInt16()),
	instanceOfInt8: reflect.PtrTo(BasicTypes.OfInt8()),
	instanceOfUint: reflect.PtrTo(BasicTypes.OfUint()),
	instanceOfUint64: reflect.PtrTo(BasicTypes.OfUint64()),
	instanceOfUint32: reflect.PtrTo(BasicTypes.OfUint32()),
	instanceOfUint16: reflect.PtrTo(BasicTypes.OfUint16()),
	instanceOfUint8: reflect.PtrTo(BasicTypes.OfUint8()),

	instanceOfFloat32: reflect.PtrTo(BasicTypes.OfFloat32()),
	instanceOfFloat64: reflect.PtrTo(BasicTypes.OfFloat64()),

	instanceOfComplex64: reflect.PtrTo(BasicTypes.OfComplex64()),
	instanceOfComplex128: reflect.PtrTo(BasicTypes.OfComplex128()),

	instanceOfByte: reflect.PtrTo(BasicTypes.OfByte()),
	instanceOfBool: reflect.PtrTo(BasicTypes.OfBool()),
	instanceOfString: reflect.PtrTo(BasicTypes.OfString()),
}

// Contains instances of "reflect.Type" over slice of builtin types of GoLang
//
// See: "reflect.SliceOf"
var SliceTypes BuiltinTypeSpace = &spaceImpl {
	instanceOfInt: reflect.SliceOf(BasicTypes.OfInt()),
	instanceOfInt64: reflect.SliceOf(BasicTypes.OfInt64()),
	instanceOfInt32: reflect.SliceOf(BasicTypes.OfInt32()),
	instanceOfInt16: reflect.SliceOf(BasicTypes.OfInt16()),
	instanceOfInt8: reflect.SliceOf(BasicTypes.OfInt8()),
	instanceOfUint: reflect.SliceOf(BasicTypes.OfUint()),
	instanceOfUint64: reflect.SliceOf(BasicTypes.OfUint64()),
	instanceOfUint32: reflect.SliceOf(BasicTypes.OfUint32()),
	instanceOfUint16: reflect.SliceOf(BasicTypes.OfUint16()),
	instanceOfUint8: reflect.SliceOf(BasicTypes.OfUint8()),

	instanceOfFloat32: reflect.SliceOf(BasicTypes.OfFloat32()),
	instanceOfFloat64: reflect.SliceOf(BasicTypes.OfFloat64()),

	instanceOfComplex64: reflect.SliceOf(BasicTypes.OfComplex64()),
	instanceOfComplex128: reflect.SliceOf(BasicTypes.OfComplex128()),

	instanceOfByte: reflect.SliceOf(BasicTypes.OfByte()),
	instanceOfBool: reflect.SliceOf(BasicTypes.OfBool()),
	instanceOfString: reflect.SliceOf(BasicTypes.OfString()),
}

// Contains instances of "reflect.Type" over array of builtin types of GoLang
//
// count - The length of array
//
// See: "reflect.ArrayOf"
func ArrayTypes(count int) BuiltinTypeSpace {
	return &arrayImpl{ count }
}

type spaceImpl struct {
	instanceOfInt reflect.Type
	instanceOfInt64 reflect.Type
	instanceOfInt32 reflect.Type
	instanceOfInt16 reflect.Type
	instanceOfInt8 reflect.Type
	instanceOfUint reflect.Type
	instanceOfUint64 reflect.Type
	instanceOfUint32 reflect.Type
	instanceOfUint16 reflect.Type
	instanceOfUint8 reflect.Type

	instanceOfFloat32 reflect.Type
	instanceOfFloat64 reflect.Type

	instanceOfComplex64 reflect.Type
	instanceOfComplex128 reflect.Type

	instanceOfByte reflect.Type
	instanceOfBool reflect.Type
	instanceOfString reflect.Type
}

func (self *spaceImpl) OfInt() reflect.Type {
	return self.instanceOfInt
}
func (self *spaceImpl) OfInt64() reflect.Type {
	return self.instanceOfInt64
}
func (self *spaceImpl) OfInt32() reflect.Type {
	return self.instanceOfInt32
}
func (self *spaceImpl) OfInt16() reflect.Type {
	return self.instanceOfInt16
}
func (self *spaceImpl) OfInt8() reflect.Type {
	return self.instanceOfInt8
}
func (self *spaceImpl) OfUint() reflect.Type {
	return self.instanceOfUint
}
func (self *spaceImpl) OfUint64() reflect.Type {
	return self.instanceOfUint64
}
func (self *spaceImpl) OfUint32() reflect.Type {
	return self.instanceOfUint32
}
func (self *spaceImpl) OfUint16() reflect.Type {
	return self.instanceOfUint16
}
func (self *spaceImpl) OfUint8() reflect.Type {
	return self.instanceOfUint8
}
func (self *spaceImpl) OfFloat32() reflect.Type {
	return self.instanceOfFloat32
}
func (self *spaceImpl) OfFloat64() reflect.Type {
	return self.instanceOfFloat64
}
func (self *spaceImpl) OfComplex64() reflect.Type {
	return self.instanceOfComplex64
}
func (self *spaceImpl) OfComplex128() reflect.Type {
	return self.instanceOfComplex128
}
func (self *spaceImpl) OfByte() reflect.Type {
	return self.instanceOfByte
}
func (self *spaceImpl) OfBool() reflect.Type {
	return self.instanceOfBool
}
func (self *spaceImpl) OfString() reflect.Type {
	return self.instanceOfString
}

type arrayImpl struct {
	count int
}

func (self *arrayImpl) OfInt() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfInt())
}
func (self *arrayImpl) OfInt64() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfInt64())
}
func (self *arrayImpl) OfInt32() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfInt32())
}
func (self *arrayImpl) OfInt16() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfInt16())
}
func (self *arrayImpl) OfInt8() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfInt8())
}
func (self *arrayImpl) OfUint() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfUint())
}
func (self *arrayImpl) OfUint64() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfUint64())
}
func (self *arrayImpl) OfUint32() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfUint32())
}
func (self *arrayImpl) OfUint16() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfUint16())
}
func (self *arrayImpl) OfUint8() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfUint8())
}
func (self *arrayImpl) OfFloat32() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfFloat32())
}
func (self *arrayImpl) OfFloat64() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfFloat64())
}
func (self *arrayImpl) OfComplex64() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfComplex64())
}
func (self *arrayImpl) OfComplex128() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfComplex128())
}
func (self *arrayImpl) OfByte() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfByte())
}
func (self *arrayImpl) OfBool() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfBool())
}
func (self *arrayImpl) OfString() reflect.Type {
	return reflect.ArrayOf(self.count, BasicTypes.OfString())
}
