package gin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_TIMEOUT = 5
)

// Constructs a new instance of "*StoppableServer" with "*gin.Engine"
func NewStoppableServer(newEngine *gin.Engine) *StoppableServer {
	return &StoppableServer{
		ginEngine: newEngine,
	}
}

// Provides methods to shutdown server gracefully.
type StoppableServer struct {
	ginEngine *gin.Engine
	httpServer *http.Server
}

// Starts server
func (self *StoppableServer) ListenAndServeAsync(addr string) {
	self.httpServer = &http.Server {
		Addr: addr,
		Handler: self.ginEngine,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("Listen address[%s] has failed: %w", addr, err))
	}

	go self.httpServer.Serve(listener)
}

// Shutdowns server with 5 seconds timeout by default.
func (self *StoppableServer) Shutdown() {
	self.ShutdownWithTimeout(DEFAULT_TIMEOUT * time.Second)
}

// Shutdowns server with timeout.
func (self *StoppableServer) ShutdownWithTimeout(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()

    if err := self.httpServer.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("Shutdown server has failed: %w", err))
    }

	self.httpServer = nil
}
