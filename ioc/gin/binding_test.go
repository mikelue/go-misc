package gin

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Binding", func() {
	Context("MvcConfig", testContextOfMvcConfig)
	Context("MvcBuilder", testContextOfMvcBuilder)
})

var testContextOfMvcBuilder = func() {
	var testedBuilder MvcBuilder
	var handlerContent *mvcSampleHandler
	var errHandler *errHandler1

	BeforeEach(func() {
		handlerContent = &mvcSampleHandler{}

		errHandler = &errHandler1{}
		testedBuilder = NewMvcConfig().
			RegisterErrorHandlers(errHandler).
			ToBuilder()
	})

	It("Handler with nothing", func() {
		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerP0R0)
		testedGinHandler(nil)

		Expect(handlerContent.v0).To(BeEquivalentTo(5909))
		Expect(errHandler.handled).To(BeFalse())
	})
	It("Handler(2 input parameters/no returned value)", func() {
		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerP2R0)
		testedGinHandler(nil)

		Expect(handlerContent.v1).To(BeEquivalentTo(3462))
		Expect(handlerContent.v2).To(BeEquivalentTo("tub"))
		Expect(errHandler.handled).To(BeFalse())
	})
	It("Handler(0 input parameter/returned value)", func() {
		context, _ := newContext()

		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerP0R1)
		testedGinHandler(context)

		Expect(context.Writer.Status()).To(BeEquivalentTo(http.StatusNotModified))
		Expect(errHandler.handled).To(BeFalse())
	})
	It("Handler(resolve parameter has error)", func() {
		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerPE1R0)
		testedGinHandler(nil)

		Expect(errHandler.handled).To(BeTrue())
	})
	It("Handler(return viable error)", func() {
		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerP0RE)
		testedGinHandler(nil)

		Expect(errHandler.handled).To(BeTrue())
	})
	It("Handler(return nil error)", func() {
		testedGinHandler := testedBuilder.WrapToGinHandler(handlerContent.handlerP0RE0)
		testedGinHandler(nil)

		Expect(errHandler.handled).To(BeFalse())
	})
}

var testContextOfMvcConfig = func() {
	var testedConfig *MvcConfig

	BeforeEach(func() {
		testedConfig = NewMvcConfig()
	})

	It("RegisterParamResolvers", func() {
		testedConfig.RegisterParamResolvers(alaniParamResolver(0), alaniParamResolver(0))
		Expect(testedConfig.paramResolvers).To(HaveLen(2))
	})

	It("RegisterParamAsFieldResolvers", func() {
		testedConfig.RegisterParamAsFieldResolvers(delanoParamAsFieldResolver(0), delanoParamAsFieldResolver(0))
		Expect(testedConfig.paramAsFieldResolvers).To(HaveLen(2))
	})

	It("RegisterErrorHandlers", func() {
		testedConfig.RegisterErrorHandlers(mungbeansErrorHandler(0), mungbeansErrorHandler(0))
		Expect(testedConfig.errorController).To(HaveLen(2))
	})
}

type param1 int
func (self *param1) Resolve(c *gin.Context) error {
	*self = 3462
	return nil
}
type param2 string
func (self *param2) Resolve(c *gin.Context) error {
	*self = "tub"
	return nil
}
type errParam string
func (self *errParam) Resolve(c *gin.Context) error {
	return fmt.Errorf("handle-1")
}

type mvcSampleHandler struct {
	v0 int
	v1 int
	v2 string
}

func (self *mvcSampleHandler) handlerP0R0() {
	self.v0 = 5909
}
func (self *mvcSampleHandler) handlerP2R0(v1 param1, v2 param2) {
	self.v1 = int(v1)
	self.v2 = string(v2)
}
func (self *mvcSampleHandler) handlerP0R1() int {
	return http.StatusNotModified
}
func (self *mvcSampleHandler) handlerPE1R0(errParam errParam) {}
func (self *mvcSampleHandler) handlerP0RE() error {
	return fmt.Errorf("handle-1")
}
func (self *mvcSampleHandler) handlerP0RE0() error {
	return nil
}

type mungbeansErrorHandler int
func (mungbeansErrorHandler) CanHandle(context *gin.Context, err error) bool {
	return true
}
func (mungbeansErrorHandler) HandleError(context *gin.Context, err error) error {
	return nil
}

type alaniParamResolver int
func (alaniParamResolver) CanResolve(srcType reflect.Type) bool {
	return true
}
func (alaniParamResolver) Resolve(context *gin.Context, srcType reflect.Type) (interface{}, error) {
	return 0, nil
}

type delanoParamAsFieldResolver int
func (delanoParamAsFieldResolver) CanResolve(*reflect.StructField) bool {
	return true
}
func (delanoParamAsFieldResolver) Resolve(*gin.Context, *reflect.StructField) (interface{}, error) {
	return 0, nil
}
