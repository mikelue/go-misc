This project contains experimental frameworks/libraries used in my work over GoLang.

# Projects

## IoC(Inverse of Control) Related

### ioc/frangipani

A poor imitation of SpringFramework of Java.

**package**: `github.com/mikelue/go-misc/ioc/frangipani` [README.md](./ioc/frangipani/README.md)

```go
env := EnvBuilder.NewByMap(map[string]interface{} {
    "global.v1": 20,
    "global.v2": 40,
})
```

### ioc/gin

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
### ioc/gorm

**package**: `github.com/mikelue/go-misc/ioc/gorm` [README.md](./ioc/gorm/README.md)

Some enhancements for [Gorm](http://gorm.io/).

Error-free coding style:
```go
// Initializes DbTemplate
tmpl := NewDbTemplate(db)

// Panic with "DbException" if the creation of object has failed
tmpl.Create(newObject)
```

## slf4go-logrus

**package**: `github.com/mikelue/go-misc/slf4go-logrus` [README.md](./slf4go-logrus/README.md)

This packages contains driver of [slf4go](https://github.com/go-eden/slf4go) with **"named logger"** by [logrus](https://github.com/sirupsen/logrus).

```go
UseLogrus.WithConfig(LogrousConfig{
    DEFAULT_LOGGER: yourDefaultLogger,
    "log.name.1": yourLogger1,
    "log.name.2": yourLogger2,
})
```

## utils

**package**: `github.com/mikelue/go-misc/utils` [README.md](./utils/README.md)

The utilities of GoLang.

### reflect

Packages:
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
