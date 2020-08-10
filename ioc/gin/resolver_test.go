package gin

import (
	"reflect"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("resolver", func() {
	Context("resolverController", func() {
		It("Has found viable resolver", func() {
			testedController := resolverController{ &sampleParamResolver{ resolvable: true, value: 40 } }
			testedBuilder := testedController.resolveBuilder(nil)

			Expect(testedBuilder).ToNot(BeNil())

			value, err := testedBuilder(nil)
			Expect(value.Interface()).To(BeEquivalentTo(40))
			Expect(err).To(Succeed())
		})
		It("Could not find resolver", func() {
			testedController := resolverController{ &sampleParamResolver{ resolvable: false } }
			Expect(testedController.resolveBuilder(nil)).To(BeNil())
		})
	})
})

type sampleParamResolver struct {
	resolvable bool
	value interface{}
}
func (self *sampleParamResolver) CanResolve(targetType reflect.Type) bool {
	return self.resolvable
}
func (self *sampleParamResolver) Resolve(context *gin.Context, targetType reflect.Type) (interface{}, error) {
	return self.value, nil
}
