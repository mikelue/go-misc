package utils

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils", func() {
	Context("EnvExecutor", func() {
		const sampleEnvVar = "GMT_VAR_1"
		const oldValue = "old_value"

		BeforeEach(func() {
			os.Setenv(sampleEnvVar, "old_value")
		})
		AfterEach(func() {
			os.Unsetenv(sampleEnvVar)
		})

		It("Variable is unchanged for Run()", func() {
			const (
				newValue = "new_value"
				newVar = "GMT_VAR_78"
				newVarValue = "new_value_2"
			)

			sampleFunc := func() {
				Expect(os.Getenv(sampleEnvVar)).
					To(BeEquivalentTo(newValue))
				Expect(os.Getenv(newVar)).
					To(BeEquivalentTo(newVarValue))
			}

			testedExecutor := NewEnvExecutor(map[string]string {
				sampleEnvVar: newValue,
				newVar: newVarValue,
			})
			testedExecutor.Run(sampleFunc)

			/**
			 * Asserts the un-changed env-variable and non-existing env-variable
			 */
			Expect(os.Getenv(sampleEnvVar)).
				To(BeEquivalentTo(oldValue))

			_, newOne := os.LookupEnv(newVar)
			Expect(newOne).To(BeFalse())
			// :~)
		})
	})
})
