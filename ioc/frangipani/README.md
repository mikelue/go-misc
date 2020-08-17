Frangipani is a poor imitation of [Spring Framework](https://spring.io/projects/spring-framework).

**package**: `github.com/mikelue/go-misc/ioc/frangipani` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/frangipani)

## Environment

Liking [Environment](https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/core/env/Environment.html) of SpringFramwork,
Frangipani provides similar interface to access properties of your application.

```go
import (
	fg "github.com/mikelue/go-misc/ioc/frangipani"
)

env := fg.EnvBuilder.NewByMap(
	map[string]interface{} {
		"v1": 20, "v2": 40,
	},
)

v1 := env.GetProperty("v1")
```

Or by [Viper](https://github.com/spf13/viper):
```go
// Supports multiple objects of viper with priority by their order.
fg.EnvBuilder.NewByViper(
	viper1, viper2
)

```

### Profile

You can access current active profiles by `Environment.GetActiveProfiles()`.

## Profiles

The property used for active profiles is `fgapp.profiles.active`.
