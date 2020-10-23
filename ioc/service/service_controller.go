package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Default signals for stopping service.
//
// These signals are used by "ServiceControllerBuilder.ByDefaultStopSignals()".
var DEFAULT_STOP_SIGNALS = []os.Signal {
	syscall.SIGINT, syscall.SIGKILL,
	syscall.SIGTERM, syscall.SIGQUIT,
}

// This channel used to capture asking for stopping services.
type StopChannel chan int

// Builder space used to build instance of "ServiceController".
var ServiceControllerBuilder IServiceControllerBuilder

type IServiceControllerBuilder int
// Constructs a controller with signals, these signal would be used with "signal.Notify" method.
//
// See: https://pkg.go.dev/os/signal#Notify
func (self *IServiceControllerBuilder) ByTrapSignals(signals ...os.Signal) ServiceController {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	stopChannel := make(StopChannel, 0)
	go func() {
		<-signalChan
		defer close(signalChan)
		close(stopChannel)
	}()

	return self.ByStopChan(stopChannel)
}
// Constructs a controller with an instance of "chan os.Signal"
//
// You can use this method to capture signals before "ServiceController",
// then re-send the signal to the channel.
//
// See: https://pkg.go.dev/os#Signal
func (self *IServiceControllerBuilder) BySignalChan(signalChan chan os.Signal) ServiceController {
	stopChannel := make(StopChannel, 0)
	go func() {
		signal := <-signalChan
		log.Printf("Received signal: %v.", signal)
		close(stopChannel)
	} ()

	return self.ByStopChan(stopChannel)
}
// Constructs a controller with an "StopChannel",
// which first message(buffered) or closed(un-buffered) would stop the running services.
func (self *IServiceControllerBuilder) ByStopChan(stopChannel StopChannel) ServiceController {
	context, cancel := context.WithCancel(context.Background())

	go func() {
		<-stopChannel
		cancel()
	}()

	return self.ByContext(context)
}
// Constructs a controller with custoimzed context.
//
// The "<-context.Done()" would be used to controler whether or not to stop the running services.
func (*IServiceControllerBuilder) ByContext(sourceContext context.Context) ServiceController {
	var waitingGroup sync.WaitGroup
	return &channelServiceController {
		mainContext: sourceContext,
		waitingGroup: &waitingGroup,
	}
}

// This controller is used to start multiple services and wait for all of them to be stopped.
//
// This object is used in "main()" function usually.
type ServiceController interface {
	// Starts a service. The service is not guaranteed ready.
	StartService(*ServiceRunner)
	// Waiting for all of the "Stop(context)" method of all of started services has get called.
	WaitForStop()
}

func init() {
	ServiceControllerBuilder = IServiceControllerBuilder(0)
}

type channelServiceController struct {
	mainContext context.Context
	waitingGroup *sync.WaitGroup
}
func (self *channelServiceController) StartService(service *ServiceRunner) {
	waitingGroup := self.waitingGroup
	waitingGroup.Add(1)

	serviceContext := self.mainContext
	serviceInfo := service.Info()
	go func() {
		if err := service.Start(serviceContext); err != nil {
			log.Printf("[%s] Starting service has error: %v", serviceInfo.Name, err)
		}
	}()
	go func() {
		<-serviceContext.Done()
		defer waitingGroup.Done()

		if err := service.Stop(serviceContext); err != nil {
			log.Printf("[%s] Stopping service has error: %v", serviceInfo.Name, err)
		}
	}()
}
func (self *channelServiceController) WaitForStop() {
	self.waitingGroup.Wait()
}
