package frangipani

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Context("AcceptsProfiles", func() {
		sampleEnv := EnvBuilder.NewByMap(map[string]interface{} {
			PROP_ACITVE_PROFILES: "a1,a2",
		})

		DescribeTable("Matching check",
			func(profiles []string, expected bool) {
				sampleProfiles := OfProfiles(profiles...)

				Expect(sampleEnv.AcceptsProfiles(sampleProfiles)).
					To(BeEquivalentTo(expected))
			},
			Entry("Matched", []string{ "a1", "a2" }, true),
			Entry("Matched", []string{ "a1", "a2", "c3" }, false),
			Entry("Matched", []string{}, true),
		)
	})

	DescribeTable("GetActiveProfiles",
		func(profiles string, expected []interface{}) {
			sampleMap := make(map[string]interface{})

			if profiles != "" {
				sampleMap[PROP_ACITVE_PROFILES] = profiles
			}

			testedEnv := EnvBuilder.NewByMap(sampleMap)

			Expect(testedEnv.GetActiveProfiles()).
				To(ConsistOf(expected...))
		},
		Entry("nothing", "", []interface{}{ DEFAULT_PROFILE }),
		Entry("2 profiles", "a1,b2", []interface{}{ "b2", "a1", DEFAULT_PROFILE }),
		Entry("2 profiles(trimming space)", "  a3 , c3  ,,", []interface{}{ "a3", "c3", DEFAULT_PROFILE }),
		Entry("duplicated profiles", "a1,b2,a1", []interface{}{ "b2", "a1", DEFAULT_PROFILE }),
	)
})
