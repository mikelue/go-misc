package service

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

// Builder space to construct "*ServiceRunner"
var ServiceBuilder IServiceBuilder

type IServiceBuilder int
// Constructs a "*ServiceRunner" with instance of "Service".
//
// The name of service would be the memory address of the object.
func (self *IServiceBuilder) New(service Service) *ServiceRunner {
	return self.NewWithInfo(
		service,
		ServiceInfo {
			Name: fmt.Sprintf("%x", reflect.ValueOf(service).Pointer()),
		},
	)
}
// Constructs a "*ServiceRunner" with instance of "Service" and information.
func (*IServiceBuilder) NewWithInfo(service Service, info ServiceInfo) *ServiceRunner {
	return &ServiceRunner {
		serviceInfo: &info,
		targetService: service,
	}
}

// Contract of "service", which is started/stopped by "ServiceController".
type Service interface {
	// Starts service with provided context, this method should block current thread.
	Start(context.Context) error
	// Stops service with provided context, this method should block current thread.
	Stop(context.Context) error
}

// Additional information to describe information about service.
type ServiceInfo struct {
	// Name of service
	Name string
}

// Executed object used by "ServiceController".
type ServiceRunner struct {
	serviceInfo *ServiceInfo
	targetService Service
}
func (self *ServiceRunner) Info() ServiceInfo {
	return *self.serviceInfo
}
func (self *ServiceRunner) Start(sourceContext context.Context) error {
	log.Printf("[%s] Starting service.", self.serviceInfo.Name)
	return self.targetService.Start(sourceContext)
}
func (self *ServiceRunner) Stop(sourceContext context.Context) error {
	log.Printf("[%s] Stopping service.", self.serviceInfo.Name)
	return self.targetService.Stop(sourceContext)
}

func init() {
	ServiceBuilder = IServiceBuilder(0)
}
