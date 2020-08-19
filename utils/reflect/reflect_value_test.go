package reflect

import (
	"reflect"
	"unsafe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValueExt", func() {
	Context("IValueExtBuilder", func() {
		It("NewByAny", func() {
			Expect(
				ValueExtBuilder.NewByAny(int16(20)).AsAny(),
			).
				To(BeEquivalentTo(20))
		})
		It("NewByValue", func() {
			Expect(
				ValueExtBuilder.NewByValue(
					reflect.ValueOf(int32(29)),
				).AsAny(),
			).
				To(BeEquivalentTo(29))
		})
	})

	Context("Methods", contextOfValueExt)
})

var contextOfValueExt = func() {
	It("AsValue", func() {
		testedValueExt := ValueExtBuilder.NewByAny(109)
		Expect(testedValueExt.AsValue().Interface()).To(BeEquivalentTo(109))
	})

	It("AsAny", func() {
		testedValueExt := ValueExtBuilder.NewByAny(133)
		Expect(testedValueExt.AsAny()).To(BeEquivalentTo(133))
	})

	Context("GetFieldValue", func() {
		box1 := sampleBox{ Age: 30, Name: "a1" }
		box2 := &sampleBox{ Age: 33, Name: "b1", ChildBox: &box1 }
		box3 := &sampleBox{ Age: 45, Name: "c1", ChildBoxL2: &box2 }

		DescribeTable("Get value from field of struct",
			func(object interface{}, tree []string, expected interface{}) {
				testedValueExt := ValueExtBuilder.NewByAny(object)
				Expect(testedValueExt.GetFieldValue(tree...).AsAny()).
					To(BeEquivalentTo(expected))
			},
			Entry("struct", box1, []string{ "Age" }, 30),
			Entry("struct field as pointer", box2, []string{ "ChildBox", "Name" }, "a1"),
			Entry("pointer of struct", box3, []string{ "Name" }, "c1"),
			Entry("field of pointer of struct as pointer of pointer", box3, []string{ "ChildBoxL2", "ChildBox", "Age" }, 30),
		)

		DescribeTable("Panic situation",
			func(invalidObject interface{}, matchedError string) {
				testedValueExt := ValueExtBuilder.NewByAny(invalidObject)

				errorMatcher := MatchError(MatchRegexp(matchedError))
				Expect(func() {
					testedValueExt.GetFieldValue("a1")
				}).To(PanicWith(errorMatcher))
			},
			Entry("Non-struct", 20, "Current type is not struct"),
			Entry("Invalid field", sampleBox{}, "is INVALID"),
		)
	})

	Context("SetFieldValue", func() {
		var sampleData [2]interface{}

		BeforeEach(func() {
			box1 := &sampleBox{ Age: 70, Name: "a3" }
			box2 := &sampleBox{ Age: 33, Name: "b1", ChildBox: box1 }

			sampleData[0] = box1
			sampleData[1] = box2
		})

		DescribeTable("Set value to field of struct",
			func(objectIndex int, tree []string, expected interface{}) {
				sampleObject := sampleData[objectIndex]

				testedValueExt := ValueExtBuilder.NewByAny(sampleObject)
				testedValueExt.SetFieldValue(ValueExtBuilder.NewByAny(expected), tree...)

				testedValueExt = ValueExtBuilder.NewByAny(sampleObject)
				Expect(testedValueExt.GetFieldValue(tree...).AsAny()).
					To(BeEquivalentTo(expected))
			},
			Entry("simple type", 0, []string{ "Age" }, 98),
			Entry("pointer", 1, []string{ "ChildBox" }, &sampleBox{ Age: 29, Name: "z0" }),
		)

		It("Set un-settable value", func() {
			testedValueExt := ValueExtBuilder.NewByAny(sampleBox{})

			Expect(func() {
				testedValueExt.SetFieldValue(
					ValueExtBuilder.NewByAny(20),
					"Age",
				)
			}).To(PanicWith(MatchError(MatchRegexp("cannot be set"))))
		})
	})

	Context("RecursiveIndirect", func() {
		s1 := 20
		s1p := &s1
		s1pp := &s1p

		DescribeTable("Structs with various indirect",
			func(v interface{}) {
				testedValueExt := ValueExtBuilder.NewByAny(v)
				Expect(
					testedValueExt.RecursiveIndirect().
						AsAny(),
				).
					To(BeEquivalentTo(20))
			},
			Entry("pure value", s1),
			Entry("pointer", s1p),
			Entry("pointer of pointer", s1pp),
		)
	})

	DescribeTable("IsArrayOrSlice",
		func(v interface{}, expected bool) {
			testedValueExt := ValueExtBuilder.NewByAny(v)
			Expect(testedValueExt.IsArrayOrSlice()).To(BeEquivalentTo(expected))
		},
		Entry("Slice", []string {}, true),
		Entry("Array", [3]int{}, true),
		Entry("string", "", false),
	)

	It("TypeExt", func() {
		testedValue := ValueExtBuilder.NewByAny(int32(345))
		Expect(testedValue.TypeExt().Kind()).To(BeEquivalentTo(reflect.Int32))
	})

	Context("IsPointer", func() {
		a := 30
		ap := &a

		DescribeTable("3 type of pointers",
			func(v interface{}, expected bool) {
				testedValueExt := ValueExtBuilder.NewByAny(v)
				Expect(testedValueExt.IsPointer()).To(BeEquivalentTo(expected))
			},
			Entry("Pointer", ap, true),
			Entry("Unsafe pointer", unsafe.Pointer(ap), true),
			Entry("Uintptr", reflect.ValueOf(ap).Pointer(), true),
			Entry("Non-pointer", a, false),
		)
	})

	Context("IsViable", func() {
		ch1, ch2 := make(chan bool, 1), make(chan bool, 1)
		ch1 <- true

		var nilErr1 error = (*simError)(nil)
		var nilErr2 *simError = (*simError)(nil)

		DescribeTable("For all of supported types",
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
}

type simStruct struct{}
type simFunc func()
type simError struct{}
func (e *simError) Error() string {
	return "OK"
}

type sampleBox struct {
	Age        int
	Name       string
	ChildBox   *sampleBox
	ChildBoxL2 **sampleBox
}
