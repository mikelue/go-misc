package gin

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("error handler", func() {
	var handler1 *errHandler1
	var handler2 *errHandler2
	var testedController errorController

	BeforeEach(func() {
		handler1 = &errHandler1{ handled: false }
		handler2 = &errHandler2{ handled: false }
		testedController = errorController{ handler1, handler2 }
	})

	Context("errorController", func() {
		It("handle by first", func() {
			testedController.handle(nil, fmt.Errorf("handle-1"))

			Expect(handler1.handled).To(BeTrue())
			Expect(handler2.handled).To(BeFalse())
		})
		It("handle by second", func() {
			testedController.handle(nil, fmt.Errorf("handle-2"))

			Expect(handler1.handled).To(BeFalse())
			Expect(handler2.handled).To(BeTrue())
		})
	})
})

type errHandler1 struct {
	handled bool
}
func (*errHandler1) CanHandle(c *gin.Context, err error) bool {
	return strings.Contains(err.Error(), "handle-1")
}
func (self *errHandler1) HandleError(c *gin.Context, err error) error {
	self.handled = true
	return nil
}

type errHandler2 struct {
	handled bool
}
func (*errHandler2) CanHandle(c *gin.Context, err error) bool {
	return strings.Contains(err.Error(), "handle-2")
}
func (self *errHandler2) HandleError(c *gin.Context, err error) error {
	self.handled = true
	return nil
}
