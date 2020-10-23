package http

import (
	"time"
	"context"
	"net/http"

	s "github.com/mikelue/go-misc/ioc/service"
)

const DEFAULT_SHUTDOWN_TIMEOUT = 8 * time.Second

// Since "*http.Server" has various way to serve HTTP,
// this type is the callback function while "Service.Start(context)" gets called.
//
// See: https://pkg.go.dev/net/http#Server
type HttpServiceStarter func(server *http.Server) error

// Builder space to construct a "Service" by "*http.Server" and "HttpServiceStarter".
//
// Start Service
//
//	The callback of "HttpServiceStarter" is the entrypoint(block) to perform starting.
//
// Stop Service
//
//	The "Shutdown(context)" method of "*http.Service" is used to perform stopping.
var HttpServiceBuilder IHttpServiceBuilder

type IHttpServiceBuilder int
// Constructs a service by "*http.Server" starter.
//
// See "DEFAULT_SHUTDOWN_TIMEOUT" for default duration of timeout.
func (self *IHttpServiceBuilder) New(server *http.Server, starter HttpServiceStarter) s.Service {
	return self.NewWithShutdownTimeout(server, starter, DEFAULT_SHUTDOWN_TIMEOUT)
}
// Constructs a service by "*http.Server" with starter and timeout value while shutdowning.
func (*IHttpServiceBuilder) NewWithShutdownTimeout(server *http.Server, starter HttpServiceStarter, shutdownTimeout time.Duration) s.Service {
	return &httpService {
		server, starter, shutdownTimeout,
	}
}

func init () {
	HttpServiceBuilder = IHttpServiceBuilder(0)
}

// Implements service.Service
type httpService struct {
	httpServer *http.Server
	startFunc HttpServiceStarter
	shutdownTimeout time.Duration
}
func (self *httpService) Start(context.Context) error {
	if err := self.startFunc(self.httpServer); err != http.ErrServerClosed {
		return err
	}

	return nil
}
func (self *httpService) Stop(sourceContext context.Context) error {
	shutdownContext, cancel := context.WithTimeout(
		context.Background(), self.shutdownTimeout,
	)
	defer cancel()

	return self.httpServer.Shutdown(shutdownContext)
}
