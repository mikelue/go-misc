package ginkgo

import (
	"strings"
	"fmt"
	"github.com/mikelue/go-misc/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LifeCycle", func() {
	Context("CycleWithOutputter", contextOfCycleWithOutputter)
})

func contextOfCycleWithOutputter() {
	var mockT *mockGinkgoT

	BeforeEach(func() {
		mockT = &mockGinkgoT {
			outputCapture: new(strings.Builder),
		}
	})

	Context("OutputIfE", func() {
		var testedOutput LifeCycle

		BeforeEach(func() {
			testedOutput = LifeCycleBuilder.ByRollbackContainerP(
				&sampleRollbackContainer { true, true },
			).OutputIfE(mockT.Log)
		})

		It("DoBefore", func() {
			testedOutput.DoBefore()
			Expect(mockT.outputCapture.String()).
				To(BeEquivalentTo("[Setup] has error: Sample error: Setup"))
		})
		It("DoAfter", func() {
			testedOutput.DoAfter()
			Expect(mockT.outputCapture.String()).
				To(BeEquivalentTo("[TearDown] has error: Sample error: TearDown"))
		})
	})

	Context("OutputfIfE", func() {
		var testedOutput LifeCycle

		BeforeEach(func() {
			testedOutput = LifeCycleBuilder.ByRollbackContainerP(
				&sampleRollbackContainer { true, true },
			).OutputfIfE(mockT.Logf, "MyE: %v")
		})

		It("DoBefore", func() {
			testedOutput.DoBefore()
			Expect(mockT.outputCapture.String()).
				To(BeEquivalentTo("[Setup] MyE: Sample error: Setup"))
		})
		It("DoAfter", func() {
			testedOutput.DoAfter()
			Expect(mockT.outputCapture.String()).
				To(BeEquivalentTo("[TearDown] MyE: Sample error: TearDown"))
		})
	})

	PContext("Sample output(by error)", func() {
		lifeCycle := LifeCycleBuilder.ByRollbackContainerP(
			&sampleRollbackContainer { true, true },
		).
			OutputIfE(GinkgoT(OFFSET).Error)

		BeforeEach(func() {
			lifeCycle.DoBefore()
		})

		It("See the output", func() {
			LifeCycleBuilder.ByRollbackContainerP(
				&sampleRollbackContainer { true, true },
			).
				OutputIfE(GinkgoT(OFFSET).Log).
				DoBefore()
		})
	})
}

type sampleRollbackContainer struct {
	setupError bool
	tearDownError bool
}
func (self *sampleRollbackContainer) Setup() (utils.Params, error) {
	if self.setupError {
		return nil, fmt.Errorf("Sample error: Setup")
	}

	return utils.Params{}, nil
}
func (self *sampleRollbackContainer) TearDown(params utils.Params) error {
	if self.setupError {
		return fmt.Errorf("Sample error: TearDown")
	}

	return nil
}

type mockGinkgoT struct {
	lastOutput int
	outputCapture *strings.Builder
}

func (self *mockGinkgoT) Fail() {
}
func (self *mockGinkgoT) Error(args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprint(args...))
	self.lastOutput = 1
}
func (self *mockGinkgoT) Errorf(format string, args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprintf(format, args...))
	self.lastOutput = 1
}
func (self *mockGinkgoT) FailNow() {
}
func (self *mockGinkgoT) Fatal(args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprint(args...))
	self.lastOutput = 2
}
func (self *mockGinkgoT) Fatalf(format string, args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprintf(format, args...))
	self.lastOutput = 2
}
func (self *mockGinkgoT) Log(args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprint(args...))
	self.lastOutput = 3
}
func (self *mockGinkgoT) Logf(format string, args ...interface{}) {
	self.outputCapture.WriteString(fmt.Sprintf(format, args...))
	self.lastOutput = 3
}
func (self *mockGinkgoT) Failed() bool {
	return true
}
func (self *mockGinkgoT) Parallel() {
}
func (self *mockGinkgoT) Skip(args ...interface{}) {
}
func (self *mockGinkgoT) Skipf(format string, args ...interface{}) {
}
func (self *mockGinkgoT) SkipNow() {
}
func (self *mockGinkgoT) Skipped() bool {
	return true
}
