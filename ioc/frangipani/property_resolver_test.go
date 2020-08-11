package frangipani

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Property resolver", func() {
	sampleProps := map[string]interface{} {
		"v1": 20,
		"v2": "Hello",
		"v11": "a1,a2,b1",
		"bs1": "32 mb",
	}

	Context("propertyResolverMapImpl", func() {
		testedImpl := mapBasedPropertyResolver(sampleProps)

		DescribeTable("ContainsProperty",
			func(name string, expected bool) {
				Expect(testedImpl.ContainsProperty(name)).
					To(BeEquivalentTo(expected))
			},
			Entry("Contains", "v1", true),
			Entry("Not contains", "v3", false),
		)

		DescribeTable("GetProperty",
			func(name string, expected string) {
				Expect(testedImpl.GetProperty(name)).
					To(BeEquivalentTo(expected))
			},
			Entry("Contains", "v1", "20"),
			Entry("Not contains", "v3", ""),
		)

		DescribeTable("GetRequiredProperty",
			func(name string, expected string, matchedErr string) {
				v, err := testedImpl.GetRequiredProperty(name)

				errMatcher := Succeed()
				if matchedErr != "" {
					errMatcher = MatchError(MatchRegexp(matchedErr))
				}

				Expect(err).To(errMatcher)
				if matchedErr == "" {
					Expect(v).To(BeEquivalentTo(expected))
				}
			},
			Entry("Contains", "v1", "20", ""),
			Entry("Not contains", "v3", "", `Property\[v3\]`),
		)
	})

	Context("requiredTypedRImpl", func() {
		testedImpl := requiredTypedRImpl(sampleProps)

		It("Get value", func() {
			v, err := testedImpl.GetUint64("v1")

			Expect(err).To(Succeed())
			Expect(v).To(BeEquivalentTo(20))
		})

		It("Get byte size(parse error)", func() {
			_, err := testedImpl.GetByteSize("v2")

			Expect(err).To(MatchError(MatchRegexp(`Unrecognized size`)))
		})

		It("Get byte size", func() {
			v, err := testedImpl.GetByteSize("bs1")

			Expect(err).To(Succeed())
			Expect(v).To(BeEquivalentTo(33554432))
		})

		It("Non-existing property", func() {
			_, err := testedImpl.Get("v3")

			Expect(err).To(MatchError(MatchRegexp(`Property\[v3\]`)))
		})

		It("Conversion error", func() {
			_, err := testedImpl.GetUint32("v2")

			Expect(err).To(MatchError(MatchRegexp(`unable to cast`)))
		})
	})

	Context("typedRImpl", func() {
		testedImpl := typedRImpl(sampleProps)

		It("Get value", func() {
			v := testedImpl.GetUint64("v1")
			Expect(v).To(BeEquivalentTo(20))
		})

		It("Get byte size(parse error)", func() {
			v := testedImpl.GetByteSize("v2")
			Expect(v).To(BeEquivalentTo(0))
		})

		It("Get byte size", func() {
			v := testedImpl.GetByteSize("bs1")
			Expect(v).To(BeEquivalentTo(33554432))
		})

		It("Non-existing property", func() {
			v := testedImpl.GetUint64("v3")
			Expect(v).To(BeEquivalentTo(0))
		})

		It("Conversion error", func() {
			v := testedImpl.GetUint32("v2")
			Expect(v).To(BeEquivalentTo(0))
		})
	})
})
