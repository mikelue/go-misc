package types

import (
	"fmt"
)

func ExampleBuiltinTypeSpace_basic() {
	fmt.Printf("%v", BasicTypes.OfInt16().Kind())
	// Output:
	// int16
}

func ExampleBuiltinTypeSpace_pointer() {
	fmt.Printf("%v", PointerTypes.OfString().Kind())
	// Output:
	// ptr
}

func ExampleBuiltinTypeSpace_slice() {
	fmt.Printf("%v", SliceTypes.OfFloat64().Kind())
	// Output:
	// slice
}

func ExampleBuiltinTypeSpace_array() {
	uint8ArrayType := ArrayTypes(5).OfUint8()

	fmt.Printf("[%d]%v", uint8ArrayType.Len(), uint8ArrayType.Kind())
	// Output:
	// [5]array
}
