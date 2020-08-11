package frangipani

import (
	"fmt"
	"github.com/spf13/viper"
)

func ExampleIEnvBuilder_newByMap() {
	env := EnvBuilder.NewByMap(map[string]interface{} {
		"db.name": "irma",
		"db.password": "XxlXt_7A",
	})

	fmt.Printf("name: %s. password: %s.", env.GetProperty("db.name"), env.GetProperty("db.password"))
	// Output:
	// name: irma. password: XxlXt_7A.
}

func ExampleIEnvBuilder_newByVipers() {
	var v1, v2, v3 *viper.Viper

	// Priority: v1 overrides v2, and v2 overrides v3
	env := EnvBuilder.NewByVipers(v1, v2, v3)

	_ = env
}
