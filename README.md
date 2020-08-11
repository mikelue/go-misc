This project contains experimental frameworks/libraries used in my work over GoLang.

# Projects

## ./ioc/frangipani

A poor imitation of SpringFramework of Java.

**package**: `github.com/mikelue/go-misc/ioc/frangipani` [README.md](./ioc/frangipani/README.md)

## ./ioc/gin

**package**: `github.com/mikelue/go-misc/ioc/gin` [README.md](./ioc/gin/README.md)

Some enhancements for [Gin Web Framework](https://onsi.github.io/ginkgo/).

The IOC of Gin can builds handler with **injected parameter** of supported types:

```go
ginHandler := builder.WrapToGinHandler()

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

## ./ioc/gorm

**package**: `github.com/mikelue/go-misc/ioc/gorm` [README.md](./ioc/gorm/README.md)

Some enhancements for [Gorm](http://gorm.io/).

Error-free coding style:
```go
// Initializes DbTemplate
tmpl := NewDbTemplate(db)

// Panic with "DbException" if the creation of object has failed
tmpl.Create(newObject)
```

## ./utils

**package**: `github.com/mikelue/go-misc/utils` [README.md](./utils/README.md)

The utilities of GoLang.

Sub-packages:
* `github.com/mikelue/go-misc/utils/reflect` [README.md](./utils/reflect/README.md)
