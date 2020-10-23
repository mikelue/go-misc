**package**: `github.com/mikelue/go-misc/ioc/service` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/service)

This package provides easy way to start multiple background services and stop them by desired [signal(IPC)](https://en.wikipedia.org/wiki/Signal_(IPC)).

Table of Contents
=================

* [Service](#service)
	* [ServiceRunner](#servicerunner)
* [ServiceController](#servicecontroller)
* [HTTP](#http)

# Service

`Service` is the main interface you should implement to **start/stop** your service.

For method of `Start(context) and Stop(context)`, you just block the thread in your code.

```go
type yourService struct {}
func (*yourService) Start(context.Context) error {
}
func (*yourService) Stop(context.Context) error {
}
```

## ServiceRunner

The needed object to be run is the struct of `ServiceRunner`,
it is responsible for additional information of a service.

```go
yourServiceRunner := service.ServiceBuilder.NewWithInfo(
  &yourService{},
  service.ServiceInfo {
    Name: "Your-srv-1",
  },
)
```

# ServiceController

This object is responsible for controlling multiple services.

In most situation, the controller is used in `main()`:

```go

func main() {
  signalChan := make(chan os.Signal, 1)
  defer close(signalChan)

  testedController := ServiceControllerBuilder.ByTrapSignals(os.Interrupt)
  testedController.StartService(yourServiceRunner)

  // Blocks until all of the services' "Stop(context)" method get called.
  testedController.WaitForStop()
}

```

# HTTP

**package**: `github.com/mikelue/go-misc/ioc/service/http` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/service/http)

You can use `HttpServiceBuilder` to build a `Service` with [*http.Server](https://pkg.go.dev/net/http#Server).

```go
httpService := HttpServiceBuilder.New(
  yourHttpServer,
  func(server *http.Server) error {
    return server.Serve()
  },
)

httpServiceRunner := service.ServiceBuilder.NewWithInfo(
  httpService,
  service.ServiceInfo {
    Name: "restful",
  },
)
```
