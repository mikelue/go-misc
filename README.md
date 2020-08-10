This project contains experimental frameworks/libraries used in my work over GoLang.

# Projects

## ./ioc/gin

Some enhances for [Gin Web Framework](https://onsi.github.io/ginkgo/). - See [README.md](./ioc/gin/README.md)

package:
```go
import(
	igin github.com/mikelue/go-misc/ioc/gin
)
```

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

## ./utils

The utilities of GoLang. See [README.md](./utils/README.md)

package:
```go
import(
	github.com/mikelue/go-misc/utils
)
```

* `github.com/mikelue/go-misc/utils/reflect` - See [README.md](./utils/reflect/README.md)
