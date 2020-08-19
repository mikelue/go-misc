package utils

import (
	"fmt"
)

type weightContainer int
func (weightContainer) Setup() error {
	fmt.Println("[1] Setup container")
	return nil
}
func (weightContainer) TearDown() error {
	fmt.Println("[3] Teardown container")
	return nil
}

func ExampleRollbackContainer() {
	RollbackExecutor.Run(
		func() {
			fmt.Println("[2] Perform job")
		},
		weightContainer(0),
	)

	// Output:
	// [1] Setup container
	// [2] Perform job
	// [3] Teardown container
}
