**package**: `github.com/mikelue/go-misc/ioc/frangipani` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/frangipani)

Frangipani is a poor imitation of [Spring Framework](https://spring.io/projects/spring-framework).

Table of Contents
=================

* [Table of Contents](#table-of-contents)
* [Environment](#environment)
* [Loading of configurations](#loading-of-configurations)
  * [Usage](#usage)
    * [Logging](#logging)
  * [Default names of configuration files](#default-names-of-configuration-files)
  * [Default priorities of loading](#default-priorities-of-loading)
    * [About XDG](#about-xdg)
  * [Profiles](#profiles)
    * [Current active profiles](#current-active-profiles)
  * [Customized loading behavior](#customized-loading-behavior)
    * [Change Prefix](#change-prefix)
    * [Change priority](#change-priority)

# Environment

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

----

# Loading of configurations

**package**: `github.com/mikelue/go-misc/ioc/frangipani/env` [GoDoc](https://pkg.go.dev/github.com/mikelue/go-misc/ioc/frangipani/env)

This package provides loading configurations by built-in mechanisms.

## Usage

With `Default().Load()`, your application would load properties(as `Env` object)
by following rules:

```go
import (
    "github.com/mikelue/go-misc/ioc/frangipani/env"
)

env := env.DefaultLoader().
    ParseFlags().
    Load()

// Default value by viper (lowest priority) for your application
configLoader := env.DefaultLoader().WithViper(yourViper)
// Default value by map(lowest priority) for your application
configLoader := env.DefaultLoader().WithMap(
    map[string]interface{} {
        "db.host": "20.10.29.189",
        "db.port": 8871,
    },
)
```

### Logging

This library uses [slf4go](https://github.com/go-eden/slf4go), you may like to set level of **`fgapp.config`** logger.

By default driver:
```go
import (
    sl "github.com/go-eden/slf4go"
)

func init() {
    sl.SetLevel(sl.InfoLevel)
}
```

By [logrus driver](https://github.com/go-eden/slf4go-logrus):
```go
import (
    sl "github.com/go-eden/slf4go-logrus"
    "github.com/sirupsen/logrus"
)

func init() {
    sl.Init()
    logrus.SetLevel(logrus.InfoLevel)
}
```

## Default names of configuration files

following list are ordered by priority for loading:
1. `fgapp-config-<profile>.properties`
1. `fgapp-config-<profile>.yaml`
1. `fgapp-config-<profile>.json`
1. `fgapp-config-default.<ext>`
1. `fgapp-config.<ext>`

`<ext>` are all of the supported formats: [properties](https://en.wikipedia.org/wiki/.properties), [YAML](https://en.wikipedia.org/wiki/YAML), [JSON](https://en.wikipedia.org/wiki/JSON)

## Default priorities of loading

The loading of configurations is as following rules(higher priority is listed first):
1. Loading configuration files in `$XDG_CONFIG_HOME/fgapp`(default: `$HOME/.config/fgapp/`)
1. Command line arguments of aggregated(ordered by priority):
    1. **`--fgapp.config.yaml=`** - YAML format of configurations
    ```sh
    your_app --fgapp.config.yaml='
    prop1: 10
    prop2: "u1:p1@someconn"
    '
    ```
    1. **`--fgapp.config.json`** - JSON format of configurations
    ```sh
    your_app --fgapp.config.json='{ "prop1": 20, "prop2": "u1:p1@someconn" }'
    ```
    1. **`--fgapp.config.files`** - Additional file of configurations, as alias of **`fgap.config.files`** property.
    ```sh
    your_app --fgapp.config.files=custom-test.yaml,custom.yaml
    ```
    1. **`--fgapp.profiles.active`** - Active profiles, as alias of **`fgapp.profiles.active`** property.
    ```sh
    your_app --fgapp.profiles.active='p1,p2'
    ```
1. Environment variables of aggregated(ordered by priority):
    1. **`$FGAPP_CONFIG_YAML`** - YAML format of configurations
    1. **`$FGAPP_CONFIG_JSON`** - JSON format of configurations
    1. **`$FGAPP_CONFIG_FILES`** - Additional file of configurations, as alias of **`fgapp.config.files`** property.
    1. **`$FGAPP_PROFILES_ACTIVE`** - Active profiles, as alias of **`fgapp.profiles.active`** property.
1. The file path from property: **`fgapp.config.file`**(ordered by priority).
    For example, if the value of `fgapp.config.files` is `my-config-test.yaml,my-config-default.yaml`:
    1. `my-config-test.yaml`
    1. `my-config-default.yaml`
1. Loading configuration files(by default names) in `$(pwd)/`
1. Loading configuration files(by default names) in same directory of [os.Args[0]](https://golang.org/pkg/os/#pkg-variables)
    * If and only if the **[working directory](https://pkg.go.dev/os?tab=doc#Getwd)** is not same as directory of [`os.Args[0]`](https://pkg.go.dev/os?tab=doc#pkg-variables).

### About XDG

Related environment variables:
* `$XDG_CONFIG_HOME`(default: `$HOME/.config`) - see [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
* GoLang: [os.UserConfigDir()](https://golang.org/pkg/os/#UserConfigDir)

## Profiles

Profiles is loaded by following property:
* `fgapp.profiles.active` - Activated profiles

The profiles are ruled by **OVERRIDING**:
```sh
your_app -fgapp.profiles.active=dev1,dev2
```

In above example, the properties of `dev1` would override the ones of `dev2`.

Additionaly, there is always a **"default"** profile appended automatically.
```properties
fgapp.profiles.active=<your profiles>,default
```

### Current active profiles

You can access current active profiles by `Environment.GetActiveProfiles()`.

```go
env := EnvBuilder.NewByMap(map[string]interface{}{
    PROP_ACITVE_PROFILES: "pf1,pf2",
})

env.GetActiveProfiles()
```

## Customized loading behavior

### Change Prefix

You can change the prefix for loading of properties:
* XDG directory: `$XDG_CONFIG_HOME/<prefix>/`
* Environment variables: `$<prefix>_CONFIG_XXX`
    * Prefix would be converted to **upper case**.
    * Character `-` would be covnerted to `_`.
* Arguments: `--<prefix>.config.xxx=`
    * Character `-` would be covnerted to `_`.
* Default names of configuration file: `<prefix>-config.<ext>`

Accepted pattern of prefix:
* First character: `[a-zA-Z]`
* Rest characters: `[a-zA-Z0-9-_]`

For example(with prefix `cherry`):
```go
env := NewConfigBuilder().
     // Default()
    Prefix("cherry").
    Build().
    ParseFlags().Load()
```

The customized loading with prefix:
* The default names of configuration files:
    1. `cherry-config-<profile>.properties`
    1. `cherry-config-<profile>.yaml`
    1. `cherry-config-<profile>.json`
    1. `cherry-config-default.<ext>`
    1. `cherry-config.<ext>`
* XDG directory: `$XDG_CONFIG_HOME/cherry/`
* Environment variable:
    1. `$CHERRY_CONFIG_YAML`
    1. `$CHERRY_CONFIG_JSON`
    1. `$CHERRY_CONFIG_FILES`
    1. `$CHERRY_PROFILES_ACTIVE`
* Arguments:
    1. `--cherry.config.yaml=''`
    1. `--cherry.config.json=''`
    1. `--cherry.config.files=''`
    1. `--cherry.profiles.active=''`

### Change priority

You can change the priority of loading of properties by `ConfigBuilder`:

```go
env := NewConfigBuilder().
    // Default()
    Priority(CL_XDG, CL_ARGS, CL_ENVVAR, CL_CONFIG_FILE, CL_PWD, CL_CMDDIR).
    Build().
    ParseFlags().Load()
```

<!-- vim: expandtab tabstop=4 shiftwidth=4
-->
