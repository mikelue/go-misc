package reflect

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FuncInfo", func() {
	testedInfo := &FuncInfo{
		inTypes: []*TypeExt{ TypeExtBuilder.NewByAny("A"), TypeExtBuilder.NewByAny((*int16)(nil)) },
		outTypes: []*TypeExt{ TypeExtBuilder.NewByAny(int32(0)), TypeExtBuilder.NewByAny(fmt.Errorf("E1")) },
	}

	Context("InTypes", func() {
		It("As []*TypeExt", func() {
			testedTypes := testedInfo.InTypes()

			Expect(testedTypes).To(HaveLen(2))
			Expect(testedTypes[0].Kind()).To(BeEquivalentTo(reflect.String))
		})

		It("As []reflect.Type", func() {
			testedTypes := testedInfo.InAsTypes()

			Expect(testedTypes).To(HaveLen(2))
			Expect(testedTypes[0].Kind()).To(BeEquivalentTo(reflect.String))
		})
	})

	Context("OutTypes", func() {
		It("As []*TypeExt", func() {
			testedTypes := testedInfo.OutTypes()

			Expect(testedTypes).To(HaveLen(2))
			Expect(testedTypes[1].Kind()).To(BeEquivalentTo(reflect.Ptr))
		})

		It("As []reflect.Type", func() {
			testedTypes := testedInfo.OutAsTypes()

			Expect(testedTypes).To(HaveLen(2))
			Expect(testedTypes[1].Kind()).To(BeEquivalentTo(reflect.Ptr))
		})
	})
})
