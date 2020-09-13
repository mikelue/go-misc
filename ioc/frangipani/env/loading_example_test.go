package env

import (
	"os"
	"fmt"

	"github.com/spf13/viper"
	"github.com/spf13/pflag"
)

func ExampleIDefaultLoader_New() {
	// This is just to make testing re-runnable
	pflag.CommandLine = pflag.NewFlagSet("ExampleIDefaultLoader_New", pflag.ExitOnError)

	os.Setenv(
		"FGAPP_CONFIG_YAML",
		`{ db.host: dev-linux, db.port: 1980 }`,
	)

	typedProperties := DefaultLoader.New().
		ParseFlags().Load().
		Typed()

	fmt.Printf("%s:%d",
		typedProperties.GetString("db.host"),
		typedProperties.GetInt("db.port"),
	)

	// Output:
	// dev-linux:1980
}

func ExampleIDefaultLoader_WithMap() {
	// This is just to make testing re-runnable
	pflag.CommandLine = pflag.NewFlagSet("ExampleIDefaultLoader_New", pflag.ExitOnError)

	os.Setenv(
		"FGAPP_CONFIG_YAML",
		`{ db.host: dev-linux }`,
	)

	typedProperties := DefaultLoader.
		WithMap(map[string]interface{} {
			"db.port": 87,
		}).
		ParseFlags().Load().
		Typed()

	fmt.Printf("%s:%d",
		typedProperties.GetString("db.host"),
		typedProperties.GetInt("db.port"),
	)

	// Output:
	// dev-linux:87
}

func ExampleIDefaultLoader_WithViper() {
	// This is just to make testing re-runnable
	pflag.CommandLine = pflag.NewFlagSet("ExampleIDefaultLoader_New", pflag.ExitOnError)

	os.Setenv(
		"FGAPP_CONFIG_YAML",
		`{ db.host: dev-linux }`,
	)

	viper := viper.New()
	viper.Set("db.port", 998)

	typedProperties := DefaultLoader.
		WithViper(viper).
		ParseFlags().Load().
		Typed()

	fmt.Printf("%s:%d",
		typedProperties.GetString("db.host"),
		typedProperties.GetInt("db.port"),
	)

	// Output:
	// dev-linux:998
}

func ExampleConfigBuilder() {
	oldOsArgs := os.Args

	defer func() {
		os.Args = oldOsArgs
	}()

	// This is just to make testing re-runnable
	pflag.CommandLine = pflag.NewFlagSet("ExampleIDefaultLoader_New", pflag.ExitOnError)

	os.Args = []string {
		`--guava.config.yaml={ db.host: apple-linux, db.port: 8871 }`,
	}
	os.Setenv(
		"GUAVA_CONFIG_YAML",
		`{ db.host: guava-linux }`,
	)

	// Makes the source of environment variables having higher priority than arguments
	typedProperties := NewConfigBuilder().
		Prefix("guava").
		Priority(CL_ENVVAR, CL_ARGS).
		Build().
		ParseFlags().Load().
		Typed()

		fmt.Printf("%s:%d", typedProperties.GetString("db.host"), typedProperties.GetInt("db.port"))

	// Output:
	// guava-linux:8871
}
