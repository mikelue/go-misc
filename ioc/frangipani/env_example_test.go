package frangipani

import (
	"fmt"
	"github.com/spf13/viper"
)

func ExampleEnvironment_getActiveProfiles() {
	env := EnvBuilder.NewByMap(map[string]interface{} {
		PROP_ACITVE_PROFILES: "pf1,pf2",
	})

	fmt.Printf("Acitve profiles: %v", env.GetActiveProfiles())
	// Output:
	// Acitve profiles: [pf1 pf2 default]
}

func ExampleEnvironment_acceptsProfiles() {
	env := EnvBuilder.NewByMap(map[string]interface{} {
		PROP_ACITVE_PROFILES: "pf1,pf2,pf98,pf99",
	})

	fmt.Printf("Accept profiles[pf2 pf98]: %v\n",
		env.AcceptsProfiles(OfProfiles("pf98", "pf2")))
	fmt.Printf("Accept profiles[pf77 pf98]: %v\n",
		env.AcceptsProfiles(OfProfiles("pf77", "pf98")))
	// Output:
	// Accept profiles[pf2 pf98]: true
	// Accept profiles[pf77 pf98]: false
}

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
	var v1, v2, v3 *viper.Viper =
		newViper(testProps { "g1": 20 }),
		newViper(testProps { "g1": 21, "g2": 30 }),
		newViper(testProps { "g1": 22, "g2": 31, "g3": 40 })

	// Priority: v1 overrides v2, and v2 overrides v3
	env := EnvBuilder.NewByVipers(v1, v2, v3)

	fmt.Printf("g1: %v. g2: %v. g3: %v.",
		env.GetProperty("g1"),
		env.GetProperty("g2"),
		env.GetProperty("g3"),
	)

	// Output:
	// g1: 20. g2: 30. g3: 40.
}

type testProps map[string]interface{}
func newViper(props testProps) *viper.Viper {
	newViper := viper.New()

	for k, v := range props {
		newViper.Set(k, v)
	}

	return newViper
}
