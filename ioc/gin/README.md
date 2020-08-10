This MVC framework is based on [Gin](https://onsi.github.io/ginkgo/).

See [godoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/gin?tab=doc) for API references.

Importing:
```go
import(
	igin "github.com/mikelue/go-misc/ioc/gin"
)
```

## Configuration and MvcBuilder

```go
config := igin.NewMvcConfig().
	RegisterParamResolvers(...).
	RegisterErrorHandlers(...)

builder := config.ToBuilder()
```

## Wrap customized handler

You can use [struct tag](https://golang.org/ref/spec#Struct_types) to **inject** parameter for your hander.

```go

type YourParams struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

func yourHandler(params &YourParams) OutputBody {
	// Use params to process your function

	return OutputBody(http.StatusOK, yourData)
}
```