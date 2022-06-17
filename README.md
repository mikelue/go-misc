This project contains experimental frameworks/libraries used in my work over GoLang. [![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://www.gnu.org/licenses/lgpl-3.0)

[![Tests all of the modules](https://github.com/mikelue/go-misc/actions/workflows/test-workflow.yaml/badge.svg)](https://github.com/mikelue/go-misc/actions/workflows/test-workflow.yaml) [![codecov](https://codecov.io/gh/mikelue/go-misc/branch/master/graph/badge.svg?token=5C7MJP5G6D)](https://codecov.io/gh/mikelue/go-misc)

Table of Contents
=================

* [IoC(Inverse of Control) Related](#iocinverse-of-control-related)
  * [ioc/frangipani](#iocfrangipani)
  * [ioc/gin](#iocgin)
  * [ioc/gorm](#iocgorm)
  * [ioc/service](#iocservice)
* [slf4go-logrus](#slf4go-logrus)
* [utils](#utils)
  * [reflect](#reflect)
* [ginkgo](#ginkgo)

# IoC(Inverse of Control) Related

## ioc/frangipani [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/ioc/frangipani.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/frangipani)

A poor imitation of SpringFramework of Java.

**package**: `github.com/mikelue/go-misc/ioc/frangipani` [README.md](./ioc/frangipani/README.md)

```go
env := EnvBuilder.NewByMap(map[string]interface{} {
    "global.v1": 20,
    "global.v2": 40,
})
```

----

**package**: `github.com/mikelue/go-misc/ioc/frangipani/env` [README.md](./ioc/frangipani/README.md)

With `env.DefaultLoader()`, you can use the out-of-box loading mechanisms of configuration:

```go
import (
    "github.com/mikelue/go-misc/ioc/frangipani/env"
)

env := env.DefaultLoader().
    ParseFlags().
    Load()
```

In above example, the environment comes from combination of configuration files(`fgapp-config.yaml`), environment variables(`$FGAPP_CONFIG_YAML`), and
arguments(`--fgapp.config.yaml`), which may be from various sources.

See [usage](./ioc/frangipani/README.md#usage)

## ioc/gin [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/ioc/gin.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/gin)

**package**: `github.com/mikelue/go-misc/ioc/gin` [README.md](./ioc/gin/README.md)

Some enhancements for [Gin Web Framework](https://onsi.github.io/ginkgo/).

The IoC of Gin can builds handler with **injected parameter** of supported types:

```go
handler := NewMvcConfig().ToBuilder().
    WrapToGinHandler(yourHandler)

ginEngine.POST("/some-resource", handler)

func yourHandler(
    params &struct {
        Id int `json:id`
        String int `json:id`
    },
) igin.OutputHandler {
    // Use params

    return igin.JsonOutputHandler(http.StatusOK, &yourResult{})
}
```
## ioc/gorm [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/ioc/gorm.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/gorm)

**package**: `github.com/mikelue/go-misc/ioc/gorm` [README.md](./ioc/gorm/README.md)

Some enhancements for [Gorm](http://gorm.io/).

Error-free coding style:
```go
// Initializes DbTemplate
tmpl := NewDbTemplate(db)

// Panic with "DbException" if the creation of object has failed
tmpl.Create(newObject)
```

## ioc/service [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/ioc/service.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/service)

**Packages**:
* `github.com/mikelue/go-misc/ioc/service` [README.md](./ioc/service/README.md)
* `github.com/mikelue/go-misc/ioc/service/http` [README.md](./ioc/service/README.md)

This package provides controller to start/stop multiple services with trapping of [os signal](https://pkg.go.dev/os#Signal).

You can start multiple services in your `main()` and stop them by trapping desired [signals(IPC)](https://en.wikipedia.org/wiki/Signal_(IPC)).

```go

// The objects implements "service.Service" interface.

func main() {
  srv1 := service.ServiceBuilder.New(initSrv1())
  srv2 := service.ServiceBuilder.New(initSrv2())

  signalChan := make(chan os.Signal, 1)
  defer close(signalChan)

  ctrl := service.ServiceControllerBuilder.ByTrapSignals(service.DEFAULT_STOP_SIGNALS...)
  ctrl.StartService(srv1)
  ctrl.StartService(srv2)

  ctrl.WaitForStop()
}
```

`github.com/mikelue/go-misc/ioc/service/http` has some wrapping method to construct a service by [http.Server](https://pkg.go.dev/net/http#Server).

# slf4go-logrus [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/slf4go-logrus.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/slf4go-logrus)

**package**: `github.com/mikelue/go-misc/slf4go-logrus` [README.md](./slf4go-logrus/README.md)

This packages contains driver of [slf4go](https://github.com/go-eden/slf4go) with **"named logger"** by [logrus](https://github.com/sirupsen/logrus).

```go
UseLogrus.WithConfig(LogrousConfig{
    DEFAULT_LOGGER: yourDefaultLogger,
    "log.name.1": yourLogger1,
    "log.name.2": yourLogger2,
})
```

# utils [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/utils.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/utils)

**package**: `github.com/mikelue/go-misc/utils` [README.md](./utils/README.md)

The utilities of GoLang.

## reflect

**Packages:**
* `github.com/mikelue/go-misc/utils/reflect` [README.md](./utils/reflect/README.md)
* `github.com/mikelue/go-misc/utils/reflect/types` [README.md](./utils/reflect/README.md)

Some convenient methods to manipulate [reflect](https://pkg.go.dev/reflect) of GoLang.

```go
valueExt := TypeExtBuilder.NewByAny(int32(0))

valueExt.IsArrayOrSlice()
valueExt.IsPointer()

concreteValue := valueExt.RecursiveIndirect()
```

<!-- vim: expandtab tabstop=4 shiftwidth=4
-->

# ginkgo [![Go Reference](https://pkg.go.dev/badge/github.com/mikelue/go-misc/ginkgo.svg)](https://pkg.go.dev/github.com/mikelue/go-misc/ginkgo)

**package**: `github.com/mikelue/go-misc/ginkgo` [README.md](./ginkgo/README.md)

Some [Ginkgo](https://onsi.github.io/ginkgo/) utilities with integration with `github.com/mikelue/go-misc/utils`.
