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
