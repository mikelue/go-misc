package frangipani

import (
	"fmt"
)

func ExamplePropertyResolver_containsProperty() {
	testedResolver := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"k1": 20,
	})

	fmt.Printf("k1: %v. k2: %v.",
		testedResolver.ContainsProperty("k1"),
		testedResolver.ContainsProperty("k2"),
	)
	// Output:
	// k1: true. k2: false.
}

func ExamplePropertyResolver_getProperty() {
	testedResolver := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"z1": 20, "z2": "hello world",
	})

	fmt.Printf("z1: %v. z2: %v.",
		testedResolver.GetProperty("z1"),
		testedResolver.GetProperty("z2"),
	)
	// Output:
	// z1: 20. z2: hello world.
}

func ExamplePropertyResolver_getRequiredProperty() {
	testedResolver := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"a1": 98,
	})

	_, err := testedResolver.GetRequiredProperty("a2")
	fmt.Printf("err: %v", err)
	// Output:
	// err: Property[a2] is not existing
}

func ExampleTypedR_getIntFamily() {
	typedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"v1": int8(20), "v2": 40,
	}).Typed()

	fmt.Printf("int: %v\n", typedR.GetInt("v1"))
	fmt.Printf("int32: %v\n", typedR.GetInt32("v1"))
	fmt.Printf("int64: %v\n", typedR.GetInt64("v1"))
	// Output:
	// int: 20
	// int32: 20
	// int64: 20
}

func ExampleTypedR_getUintFamily() {
	typedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"v1": uint8(27), "v2": 40,
	}).Typed()

	fmt.Printf("uint: %v\n", typedR.GetUint("v1"))
	fmt.Printf("uint32: %v\n", typedR.GetUint32("v1"))
	fmt.Printf("uint64: %v\n", typedR.GetUint64("v1"))
	// Output:
	// uint: 27
	// uint32: 27
	// uint64: 27
}

func ExampleTypedR_getByteSize() {
	typedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"file.limit.size": "10 mb",
	}).Typed()

	fmt.Printf(
		typedR.GetByteSize("file.limit.size").
			Format("Limit: %.0f", "kb", false),
	)
	// Output:
	// Limit: 10240KB
}

func ExampleRequiredTypedR_getByteSize() {
	requiredTypedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"file.limit.size": "10 mb",
		"size.wrong-format": "10 cm",
	}).RequiredTyped()

	sizeObj, _ := requiredTypedR.GetByteSize("file.limit.size")
	fmt.Println(sizeObj.Format("Limit: %.0f", "kb", false))

	_, err := requiredTypedR.GetByteSize("size.wrong-format")
	fmt.Printf("Format error: %v\n", err)

	// Output:
	// Limit: 10240KB
	// Format error: Unrecognized size suffix cm
}

func ExampleRequiredTypedR_notExistingProperty() {
	requiredTypedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"v3": 30, "v4": 80,
	}).RequiredTyped()

	_, err := requiredTypedR.GetInt("v1")

	fmt.Printf("err: %v", err)
	// Output:
	// err: Property[v1] is not existing
}

func ExampleRequiredTypedR_getSliceFamily() {
	requiredTypedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"s1": []int{ 90, 91, 92 },
		"s2": []string{ "GP-1", "GP-2", "GP-3" },
		"s-wrong": map[int]int {},
	}).RequiredTyped()

	s1, _ := requiredTypedR.GetIntSlice("s1")
	fmt.Printf("int slice: %v\n", s1)

	s2, _ := requiredTypedR.GetStringSlice("s2")
	fmt.Printf("string slice: %v\n", s2)

	_, e := requiredTypedR.GetStringSlice("s-wrong")
	fmt.Printf("wrong type: %v\n", e)
	// Output:
	// int slice: [90 91 92]
	// string slice: [GP-1 GP-2 GP-3]
	// wrong type: unable to cast map[int]int{} of type map[int]int to []string
}
