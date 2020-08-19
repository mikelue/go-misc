package utils

import (
	"fmt"
)

type weightContainerP int
func (weightContainerP) Setup() (Params, error) {
	fmt.Println("[1] Setup container")
	return Params{ "v1": 871 }, nil
}
func (weightContainerP) TearDown(params Params) error {
	fmt.Printf("[3] Teardown container: %v\n", params["v2"])
	return nil
}

func ExampleRollbackContainerP() {
	RollbackExecutor.RunP(
		func(params Params) {
			fmt.Printf("[2] Perform job: %v\n", params["v1"])
			params["v2"] = 145
		},
		weightContainerP(0),
	)

	// Output:
	// [1] Setup container
	// [2] Perform job: 871
	// [3] Teardown container: 145
}
