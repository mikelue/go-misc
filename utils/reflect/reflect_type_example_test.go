package reflect

import (
	"fmt"
	"reflect"
)

func ExampleITypeExtBuilder_newByAny() {
	typeExt := TypeExtBuilder.NewByAny(20)
	fmt.Printf("Type extension: %v", typeExt.Kind())
	// Output:
	// Type extension: int
}

func ExampleITypeExtBuilder_newByType() {
	sampleString := "Hello"
	sampleType := reflect.TypeOf(&sampleString)

	typeExt := TypeExtBuilder.NewByType(sampleType)
	fmt.Printf("Type extension: %v", typeExt.Kind())
	// Output:
	// Type extension: ptr
}

func ExampleTypeExt_recursiveIndirect() {
	var u1 uint32 = 981
	u1p := &u1
	u1pp := &u1p

	u1ppTypeExt := TypeExtBuilder.NewByAny(u1pp)

	fmt.Printf("Type: %v. Leaf type: %v.",
		u1ppTypeExt.Kind(),
		u1ppTypeExt.RecursiveIndirect().Kind(),
	)
	// Output:
	// Type: ptr. Leaf type: uint32.
}
