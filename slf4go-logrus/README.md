**package**: `github.com/mikelue/go-misc/slf4go-logrus`

This packages contains driver of [slf4go](https://github.com/go-eden/slf4go) with **"named logger"** by [logrus](https://github.com/sirupsen/logrus).

## Usage

Besides `UseLogrus.Default()`, you can customized your logger with specific name.

```go
UseLogrus.WithConfig(LogrousConfig{
    DEFAULT_LOGGER: yourDefaultLogger,
    "log.name.1": yourLogger1,
    "log.name.2": yourLogger2,
})
```

<!-- vim: expandtab tabstop=4 shiftwidth=4
-->
