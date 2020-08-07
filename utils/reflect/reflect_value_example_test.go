package reflect

import (
	"fmt"
	"reflect"
)

func ExampleIValueExtBuilder_newByAny() {
	valueExt := ValueExtBuilder.NewByAny(20)
	fmt.Printf("Value: %v", valueExt.AsValue().Interface())
	// Output:
	// Value: 20
}

func ExampleIValueExtBuilder_newByValue() {
	sampleValue := reflect.ValueOf("Hello")

	valueExt := ValueExtBuilder.NewByValue(sampleValue)
	fmt.Printf("Value: %v", valueExt.AsValue().Interface())
	// Output:
	// Value: Hello
}
