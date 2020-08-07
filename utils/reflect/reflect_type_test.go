package reflect

import (
	"reflect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeExt", func() {
	Context("TypeExtBuilder", func() {
		It("NewByType", func() {
			testedType := TypeExtBuilder.NewByType(reflect.TypeOf(int(3)))
			Expect(testedType.Kind()).To(BeEquivalentTo(reflect.Int))
		})
		It("NewByAny", func() {
			testedType := TypeExtBuilder.NewByAny(int16(3))
			Expect(testedType.Kind()).To(BeEquivalentTo(reflect.Int16))
		})
	})

	Context("TypeExt", contextOfTypeExt)
})

var contextOfTypeExt = func() {
	It("AsType", func() {
		testedType := TypeExtBuilder.NewByAny("Cool")
		Expect(testedType.Kind()).To(BeEquivalentTo(reflect.String))
	})

	It("Kind", func() {
		testedType := TypeExtBuilder.NewByAny(int64(981))
		Expect(testedType.Kind()).To(BeEquivalentTo(reflect.Int64))
	})

	DescribeTable("IsReflectType",
		func(v interface{}, expected bool) {
			testedType := TypeExtBuilder.NewByAny(v)
			Expect(testedType.IsReflectType()).To(BeEquivalentTo(expected))
		},
		Entry(`Implement "reflect.Type"`, reflect.TypeOf("Hello"), true),
		Entry(`Not implement "reflect.Type"`, "Hello", false),
	)

	Context("InterfaceType", func() {
		It("Success extracted", func() {
			testedType := TypeExtBuilder.NewByAny((*sampleI1)(nil))
			Expect(testedType.InterfaceType().Kind()).To(BeEquivalentTo(reflect.Interface))
		})

		DescribeTable("Failed extracting",
			func(value interface{}, matchedPattern string) {
				Expect(
					func() { TypeExtBuilder.NewByAny(value).InterfaceType() },
				).To(PanicWith(MatchError(MatchRegexp(matchedPattern))))
			},
			Entry("Non-pointer", 98, "should be pointer"),
			Entry("Not pointer-to-interface", &sampleBox{}, "kind of expected type"),
		)
	})

	It("NewAsPointer", func() {
		testedValue := TypeExtBuilder.NewByAny(int64(876)).
			NewAsPointer().AsValue()

		Expect(testedValue.Kind()).To(BeEquivalentTo(reflect.Ptr))
	})

	It("NewAsValue", func() {
		testedValue := TypeExtBuilder.NewByAny((*int16)(nil)).
			NewAsValue().AsAny()

		Expect(testedValue).To(BeEquivalentTo(0))
	})

	Context("FuncInfo", func() {
		It("Valid information", func() {
			testedFuncInfo := TypeExtBuilder.NewByAny(sampleF1).
				FuncInfo()

			testedInTypes := testedFuncInfo.InTypes()
			Expect(testedInTypes).To(HaveLen(2))
			Expect(testedInTypes[0].Kind()).To(BeEquivalentTo(reflect.String))

			testedOutType := testedFuncInfo.OutTypes()
			Expect(testedOutType).To(HaveLen(2))
			Expect(testedOutType[0].Kind()).To(BeEquivalentTo(reflect.Int64))
		})

		It("Non-function", func() {
			testedType := TypeExtBuilder.NewByAny(23)

			Expect(
				func() { testedType.FuncInfo() },
			).To(PanicWith(MatchError(MatchRegexp(`is not "reflect.Func"`))))
		})
	})

	Context("RecursiveIndirect", func() {
		a := 20
		ap := &a
		app := &ap

		DescribeTable("Different level of pointer",
			func(v interface{}) {
				testedType := TypeExtBuilder.NewByAny(v).
					RecursiveIndirect().AsType()

				Expect(testedType.Kind()).To(BeEquivalentTo(reflect.Int))
			},
			Entry("Non-pointer", a),
			Entry("Pointer to value", ap),
			Entry("Pointer of pointer", app),
		)
	})
}

func sampleF1(a string, b *uint32) (int64, error) {
	return 0, nil
}

type sampleI1 interface {
	Do1()
}
