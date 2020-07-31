package reflect

import (
	"reflect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type simStruct struct{}
type simFunc func()
type simError struct{}
func (e *simError) Error() string {
	return "OK"
}

var _ = Describe("Is viable", func() {
	ch1, ch2 := make(chan bool, 1), make(chan bool, 1)
	ch1 <- true

	var nilErr1 error = (*simError)(nil)
	var nilErr2 *simError = (*simError)(nil)

	DescribeTable("result as expected one",
		func(sampleValue interface{}, expected bool) {
			testedValue := reflect.ValueOf(sampleValue)
			testedResult := ValueExt(testedValue).IsViable()

			Expect(testedResult).To(Equal(expected))
		},
		Entry("30 is viable", 30, true),
		Entry("0 is viable", 0, true),
		Entry("Initialized pointer to *struct is viable", &simStruct{}, true),
		Entry("Nil pointer to *struct is not viable", (*simStruct)(nil), false),
		Entry("Non-empty array is viable", []int{20}, true),
		Entry("Empty array is not viable", []int{}, false),
		Entry("Nil array is not viable", []string(nil), false),
		Entry("Non-empty array(element's type is pointer) is viable", []*simStruct{{}}, true),
		Entry("Nil array(element's type is pointer) is not viable", []*simStruct{}, false),
		Entry("Non-empty map is viable", map[int]bool{20: true}, true),
		Entry("Empty map is not viable", map[int]bool{}, false),
		Entry("Function is viable", simFunc(func() {}), true),
		Entry("Nil function is not viable", simFunc(nil), false),
		Entry("Non-empty channel is viable", ch1, true),
		Entry("Empty channel is not viable", ch2, false),
		Entry("Nil error(pure) is not viable", (error)(nil), false),
		Entry("Nil error is not viable", nilErr1, false),
		Entry("Nil error(alias) is not viable", nilErr2, false),
	)
})
