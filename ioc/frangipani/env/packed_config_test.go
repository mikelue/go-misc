package env

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packed configurations", func() {
	Context("packedConfig", contextOfPackedConfig)

	Context("readInByString", func() {
		It("Successful read in", func() {
			viper, isSuccess := readInByString("sample.prop.json", "json", `{ "some.key1": 20 }`)

			Expect(isSuccess).To(BeTrue())
			Expect(viper.GetInt("some.key1")).To(BeEquivalentTo(20))
		})
		It("Failed read in", func() {
			viper, isSuccess := readInByString("sample.prop.json", "json", `{ some.key1: 20 }`)

			Expect(isSuccess).To(BeFalse())
			Expect(viper).To(BeNil())
		})
	})
})

var contextOfPackedConfig = func() {
	It("loadFormattedProps", func() {
		sampleKycArgs := &packedConfig {
			yamlProps: `{ "v1": 72 }`,
			jsonProps: `{ "v1": 40 }`,
		}

		testedVipers := sampleKycArgs.loadFormattedProps()

		Expect(testedVipers).To(HaveLen(2))
		Expect(testedVipers[0].GetInt("v1")).To(BeEquivalentTo(72))
		Expect(testedVipers[1].GetInt("v1")).To(BeEquivalentTo(40))
	})

	It("loadFormattedProps with error(cannot parse)", func() {
		sampleKycArgs := &packedConfig {
			yamlProps: `{ "v1": [[ }`,
			jsonProps: `{ "v1": ]] }`,
		}

		testedVipers := sampleKycArgs.loadFormattedProps()

		Expect(testedVipers).To(HaveLen(0))
	})

	It("configFileProp", func() {
		sampleKycArgs := &packedConfig {
			externalFiles: "config-test.yaml, config-assure.yaml ,,",
		}

		testedVipers := sampleKycArgs.configFileProp()
		Expect(testedVipers.GetStringSlice(PROP_CONFIG_FILES)).
			To(ConsistOf("config-test.yaml", "config-assure.yaml"))
	})

	It("loadActiveProfiles", func() {
		sampleKycArgs := &packedConfig {
			activeProfiles: "a1,a2",
		}

		testedVipers := sampleKycArgs.loadActiveProfiles()
		Expect(testedVipers).ToNot(BeNil())
		Expect(testedVipers.GetString(PROP_PROFILES_ACTIVE)).
			To(BeEquivalentTo("a1,a2"))
	})

	DescribeTable("loadAll",
		func(sampleArgs *packedConfig, matcher OmegaMatcher) {
			testedResult := sampleArgs.load()

			Expect(testedResult).To(matcher)
		},
		Entry("Nothing", &packedConfig{}, HaveLen(0)),
		Entry("YAML and JSON", &packedConfig{ yamlProps: `{}`, jsonProps: `{}` },
			HaveLen(2),
		),
		Entry("External file", &packedConfig{ externalFiles: fmt.Sprintf("%s/%s", currentSrcDir, "sample-1.yaml") },
			HaveLen(1),
		),
		Entry("Active profiles", &packedConfig{ activeProfiles: `pf1,pf2` },
			HaveLen(1),
		),
	)
}
