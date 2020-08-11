package frangipani

import (
	"github.com/spf13/viper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment by Viper", func() {
	Context("newEnvViperImpl", func() {
		v1 := viper.New()
		v1.MergeConfigMap(map[string]interface{} {
			"v1": 20, "v2": 40,
		})
		v2 := viper.New()
		v2.MergeConfigMap(map[string]interface{} {
			"v1": 22, "v2": 42, "v3": 52,
		})

		DescribeTable("Constructs and test single property",
			func(sampleVipers []*viper.Viper, name string, expectedValue int) {
				testedEnv := newEnvViperImpl(sampleVipers...)

				Expect(testedEnv.Typed().GetInt(name)).
					To(BeEquivalentTo(expectedValue))
			},
			Entry("Single *Viper", []*viper.Viper{ v1 }, "v1", 20),
			Entry("Overriding *Viper", []*viper.Viper{ v2, v1 }, "v1", 22),
			Entry("Non-overriding *Viper", []*viper.Viper{ v1, v2 }, "v3", 52),
		)
	})
})
