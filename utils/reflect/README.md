**package**: `github.com/mikelue/go-misc/utils` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/utils)

This package contains some enhancements for manipulating [reflect](https://pkg.go.dev/reflect?tab=doc) of GoLang.

## reflect/

**package**: `github.com/mikelue/go-misc/utils/reflect` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/utils/reflect)

This package provides out-of-box methods to manipulate instances of **reflect** easily.

`TypeExt` - Some convenient methods to manipulate [`reflect.Type`](https://pkg.go.dev/reflect?tab=doc#Type).

`ValueExt` - Some convenient methods to manipulate [`reflect.Value`](https://pkg.go.dev/reflect?tab=doc#Value).

`AnyValue` - Some convenient methods to manipulate [`interface{}`](https://golang.org/ref/spec#Interface_types).

### reflect/types
**package**: `github.com/mikelue/go-misc/utils/reflect/types` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/utils/reflect/types)

`BasicTypes` - Instance can be used directly without calling `reflect.TypeOf(<type>)` for builtin types provided by GoLang.
This package provides instances of [reflect.Type](https://pkg.go.dev/reflect?tab=doc#Type) or [reflect.Value](https://pkg.go.dev/reflect?tab=doc#Value) for builtin types of GoLang.

`PointerTypes` - Instance can be used directly without calling `reflect.PtrTo(<type>)` for builtin types provided by GoLang.

`SliceTypes` - Instance can be used directly without calling `reflect.SliceOf(<type>)` for builtin types provided by GoLang.

`ArrayTypes` - Instance can be used directly without calling `reflect.ArrayOf(<type>)` for builtin types provided by GoLang.

```go
import (
	"reflect"
	t "github.com/mikelue/go-misc/utils/reflect/types"
)

// New a pointer to "uint64"
pointerValueOfInt64 := reflect.New(t.BasicTypes.OfUint64())
```
