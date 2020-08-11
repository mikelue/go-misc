package frangipani

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Profiles", func() {
	DescribeTable("OfProfiles",
		func(sampleProfiles []string, expected []interface{}) {
			testedProfiles := OfProfiles(sampleProfiles...)

			Expect(testedProfiles).
				To(ConsistOf(expected...))
		},
		Entry("normal case", []string{ "e1", "e2" }, []interface{}{ "e1", "e2" }),
		Entry("trimming space", []string{ " g1 ", " g2 ", "  " }, []interface{}{ "g1", "g2" }),
	)

	DescribeTable("Matches",
		func(sampleProfiles []string, expected bool) {
			Expect(
				OfProfiles(sampleProfiles...).Matches(matchM1AndM2),
			).To(BeEquivalentTo(expected))
		},
		Entry("matched", []string{ "m1", "m2" }, true),
		Entry("not matched", []string{ "m1", "n1" }, false),
		Entry("nothing, matched", []string{}, true),
	)
})

func matchM1AndM2(profile string) bool {
	return profile == "m1" || profile == "m2"
}
