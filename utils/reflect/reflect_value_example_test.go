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

func ExampleValueExt_recursiveIndirect() {
	var s1 string = "hello"
	s1p := &s1
	s1pp := &s1p

	s1ppExt := ValueExtBuilder.NewByAny(s1pp)

	fmt.Printf("Indirected value: %s", s1ppExt.RecursiveIndirect().AsAny())
	// Output:
	// Indirected value: hello
}

func ExampleValueExt_isViable() {
	printOut := func(content string, v interface{}) {
		fmt.Printf(content + ": %v\n", ValueExtBuilder.NewByAny(v).IsViable())
	}

	/**
	 * True
	 */
	printOut("0", 0)
	printOut("false", false)
	printOut("<empty string>", "")
	bufferedChan := make(chan int, 2)
	bufferedChan <- 1
	printOut("chan int(one element)", bufferedChan)
	// :~)

	fmt.Println("--------------------")

	/**
	 * False
	 */
	printOut("(*int)(nil)", (*int)(nil))
	printOut("([]string)(nil)", ([]string)(nil))
	printOut("[0]string", [0]string{})
	printOut("[]uint64(len == 0)", make([]uint64, 0, 4))
	printOut("map[int]string(empty)", make(map[int]string))

	<- bufferedChan
	printOut("chan int(empty)", bufferedChan)
	// :~)

	// Output:
	// 0: true
	// false: true
	// <empty string>: true
	// chan int(one element): true
	// --------------------
	// (*int)(nil): false
	// ([]string)(nil): false
	// [0]string: false
	// []uint64(len == 0): false
	// map[int]string(empty): false
	// chan int(empty): false
}
