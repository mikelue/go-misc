package frangipani

import (
	"fmt"
)

func ExampleTypedR() {
	typedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"v1": 20, "v2": 40,
	}).Typed()

	fmt.Printf("V1: %d. V2: %d.", typedR.GetInt("v1"), typedR.GetInt("v2"))
	// Output:
	// V1: 20. V2: 40.
}

func ExampleTypedR_getGetByteSize() {
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

func ExampleRequiredTypedR() {
	requiredTypedR := PropertyResolverBuilder.NewByMap(map[string]interface{} {
		"v3": 30, "v4": 80,
	}).RequiredTyped()

	_, err := requiredTypedR.GetInt("v1")

	fmt.Printf("err: %v", err)
	// Output:
	// err: Property[v1] is not existing
}
