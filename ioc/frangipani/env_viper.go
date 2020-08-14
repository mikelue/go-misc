package frangipani

import (
	"github.com/spf13/viper"
)

func newEnvViperImpl(vipers ...*viper.Viper) Environment {
	props := make(map[string]interface{})

	/**
	 * Use the overriding method to take priority of properties
	 */
	for i := len(vipers) - 1; i >= 0; i-- {
		for _, k := range vipers[i].AllKeys() {
			props[k] = vipers[i].Get(k)
		}
	}
	// :~)

	return EnvBuilder.NewByMap(props)
}
