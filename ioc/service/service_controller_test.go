package service

import (
	"context"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceController", func() {
	allServices := sampleServices {
		&timedService{ false, 100 * time.Millisecond, sync.Mutex{} },
		&timedService{ false, 100 * time.Millisecond, sync.Mutex{} },
	}

	startServices := func(ctrl ServiceController) {
		allServices.startAll(ctrl)

		Eventually(
			allServices.allStatus,
			2 * time.Second, time.Second,
		).Should(BeTrue())
	}
	assertStop := func(ctrl ServiceController) {
		stopBox := make(chan bool, 0)

		go func() {
			ctrl.WaitForStop()
			stopBox <- true
		} ()

		Eventually(
			allServices.allStatus,
			2 * time.Second, time.Second,
		).Should(BeFalse())

		Eventually(
			func() bool { return <-stopBox },
			2 * time.Second, time.Second,
		).Should(BeTrue())
	}
	simulateRunning := func(stopper func()) {
		<-time.After(250 * time.Millisecond)
		stopper()
	}

	Context("Start/Stop services", func() {
		It("By chan os.Signal", func() {
			sampleSignalChan := make(chan os.Signal, 1)
			defer close(sampleSignalChan)
			testedController := ServiceControllerBuilder.BySignalChan(sampleSignalChan)

			startServices(testedController)

			go simulateRunning(
				func() {
					sampleSignalChan<-os.Kill
				},
			)

			assertStop(testedController)
		})
		It("By context.Context", func() {
			sampleContext, cancel := context.WithCancel(context.TODO())

			testedController := ServiceControllerBuilder.ByContext(sampleContext)
			startServices(testedController)

			go simulateRunning(cancel)

			assertStop(testedController)
		})

		Context("By StopChannel", func() {
			It("Buffered channel", func() {
				sampleStopChan := make(StopChannel, 1)
				defer close(sampleStopChan)

				testedController := ServiceControllerBuilder.ByStopChan(sampleStopChan)
				startServices(testedController)

				go simulateRunning(
					func() {
						sampleStopChan <- 0
					},
				)

				assertStop(testedController)
			})
			It("Un-buffered channel", func() {
				sampleStopChan := make(StopChannel, 0)

				testedController := ServiceControllerBuilder.ByStopChan(sampleStopChan)
				startServices(testedController)

				go simulateRunning(
					func() {
						close(sampleStopChan)
					},
				)

				assertStop(testedController)
			})
		})
	})
})

type sampleServices []*timedService
func (self sampleServices) allStatus() bool {
	for _, service := range self {
		if !service.getStatus() {
			return false
		}
	}

	return true
}
func (self sampleServices) startAll(ctrl ServiceController) {
	for _, service := range self {
		ctrl.StartService(
			ServiceBuilder.New(service),
		)
	}
}

type timedService struct {
	running bool
	shutdownWait time.Duration
	lock sync.Mutex
}
func (self *timedService) Start(context context.Context) error {
	self.setStatus(true)
	for self.getStatus() {
		select {
		case <-context.Done():
			self.setStatus(false)
		case <-time.After(3 * time.Second):
		}
	}

	return nil
}
func (self *timedService) Stop(context context.Context) error {
	self.setStatus(false)
	<-time.After(self.shutdownWait)
	return nil
}
func (self *timedService) setStatus(value bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.running = value
}
func (self *timedService) getStatus() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.running
}
