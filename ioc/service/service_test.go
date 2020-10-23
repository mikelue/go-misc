package service

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("serviceInfoImpl", func() {
		It("Start/Stop service", func() {
			testedService := new(sampleService)
			testedRunner := ServiceBuilder.New(testedService)

			By("Start server")
			testedRunner.Start(context.TODO())
			Expect(testedService.running).To(BeTrue())

			By("Stop server")
			testedRunner.Stop(context.TODO())
			Expect(testedService.running).To(BeFalse())
		})
	})
})

type sampleService struct {
	running bool
}
func (self *sampleService) Start(context context.Context) error {
	self.running = true
	return nil
}
func (self *sampleService) Stop(context context.Context) error {
	self.running = false
	return nil
}
