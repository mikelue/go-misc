package ginkgo

import (
	"fmt"

	"github.com/mikelue/go-misc/utils"
)

// Offset can be used to report error in Before/After blocks
//
//  lifeCycle := LifeCycleBuilder.ByRollbackContainer(your_container).
//    OutputIfE(GinkgoT(OFFSET).Error)
//
//  BeforeEach(func() {
//    lifeCycle.DoBefore()
//  })
//  AfterEach(func() {
//    lifeCycle.DoAfter()
//  })
const OFFSET = 4

// Defines the life cycle of a test
type LifeCycle interface {
	// Runs the before action and output something.
	DoBefore()
	// Runs the after action and output something.
	DoAfter()
}

// Direct usage to "ILifeCycleBuilder".
var LifeCycleBuilder ILifeCycleBuilder = 0

// Method space used to build "CycleWithOutputter"s
type ILifeCycleBuilder int
// Builds a CycleWithOutputter by "RollbackContainer"
func (ILifeCycleBuilder) ByRollbackContainer(container utils.RollbackContainer) CycleWithOutputter {
	return &lifeCycleForContainer{ container }
}
// Builds a CycleWithOutputter by "RollbackContainerP"
func (ILifeCycleBuilder) ByRollbackContainerP(containerP utils.RollbackContainerP) CycleWithOutputter {
	return &lifeCycleForContainer{ utils.RollbackContainerBuilder.ToContainer(containerP) }
}

type CycleWithOutputter interface {
	// Output error with method of "Error", "Log", or "Panic" of "GinkgoTInterface" if
	// the "Setup()" or "TearDown()" method returns any viable error.
	OutputIfE(Output) LifeCycle
	// The string argument is the "format"(should have a "%s") of "GinkgoT().Errorf(format)" if
	// the "Setup()" or "TearDown()" method returns any viable error.
	//
	// The "[Setup]" or "[TearDown]" would be prefixed to the output mesasge.
	OutputfIfE(Outputf, string) LifeCycle
}

// The function used to output normal message.
type Output func(...interface{})
// The function used to output message with designated string of formatting.
type Outputf func(string, ...interface{})

type lifeCycleForContainer struct {
	container utils.RollbackContainer
}
func (self *lifeCycleForContainer) OutputIfE(output Output) LifeCycle {
	return &lifeCycleAndOutputForContainer {
		container: self.container,
		outputImpl: newOutput(output),
	}
}
func (self *lifeCycleForContainer) OutputfIfE(outputf Outputf, format string) LifeCycle {
	return &lifeCycleAndOutputForContainer {
		container: self.container,
		outputImpl: newOutputf(format, outputf),
	}
}

type lifeCycleAndOutputForContainer struct {
	params utils.Params
	container utils.RollbackContainer
	outputImpl *outputImpl
}
func (self *lifeCycleAndOutputForContainer) DoBefore() {
	err := self.container.Setup()
	self.outputError("Setup", err)
}
func (self *lifeCycleAndOutputForContainer) DoAfter() {
	err := self.container.TearDown()
	self.outputError("TearDown", err)
}
func (self *lifeCycleAndOutputForContainer) outputError(lifeCycle string, err error) {
	if err == nil {
		return
	}

	outputImpl := self.outputImpl
	switch outputImpl.outputType {
	case 1:
		self.outputImpl.output(fmt.Sprintf("[%s] has error: ", lifeCycle), err)
	case 2:
		self.outputImpl.outputf(fmt.Sprintf("[%s]", lifeCycle), err)
	}
}

func newOutputf(newFormat string, newOutputfImpl func(string, ...interface{})) *outputImpl {
	return &outputImpl {
		outputType: 2,
		format: newFormat,
		outputfImpl: newOutputfImpl,
	}
}
func newOutput(newOutputImpl func(s...interface{})) *outputImpl {
	return &outputImpl {
		outputType: 1,
		outputImpl: newOutputImpl,
	}
}

type outputImpl struct {
	outputType int
	format string
	outputImpl func(... interface{})
	outputfImpl func(string, ...interface{})
}
func (self *outputImpl) output(message string, e error) {
	self.outputImpl(message, e)
}
func (self *outputImpl) outputf(prefix string, e error) {
	self.outputfImpl(fmt.Sprintf("%s %s", prefix, self.format), e)
}
