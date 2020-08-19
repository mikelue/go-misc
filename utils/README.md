**package**: `github.com/mikelue/go-misc/utils` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/utils)

This package contains language libraries over GoLang.

[./reflect](./reflect) - Some convenient libraries to manipulate [reflect](https://pkg.go.dev/reflect?tab=doc) in GoLang.

# RollbackContainer

`RollbackContainer` is an interface with **"Setup()"**/**"TearDown"** life-cycle to wrap certain
function with fabricated environment.

The `RollbackExecutor` can be used to run a function(`func() {}`) wrapped by **"1:M"** `RollbackContainer`s.

```go
tempDir := RollbackContainerBuilder.NewTmpDir("my-temp-*")

RollbackExecutor.RunP(
    func(params Params) {
        /* The directory of created is params[PKEY_TEMP_DIR] */
    },
    tempDir,
)
```

`RollbackContainerBuilder` provides some methods to constructs useful containers:
* `NewTmpDir()` - Builds temporary directory and remove it afterward, see [RollbackContainerP](https://pkg.go.dev/github.com/mikelue/go-misc/utils?tab=doc#RollbackContainerP).
* `NewEnv()` - Sets up environment variable and reset them afterward.
* `NewChdir()` - Changes current directory and change-back afterward.

See [GoDoc go-misc/utils](https://pkg.go.dev/github.com/mikelue/go-misc/utils)

<!-- vim: expandtab tabstop=4 shiftwidth=4
-->
