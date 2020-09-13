package env

import (
	"fmt"

	"github.com/spf13/pflag"
)

const (
	// Properties packed as JSON format
	FLAG_CONFIG_JSON = ".config.json"
	// Properties packed as YAML format
	FLAG_CONFIG_YAML = ".config.yaml"
	// Properties from external file
	FLAG_CONFIG_FILES = PROP_CONFIG_FILES
	// For activated profiles(see frangipani)
	FLAG_ACITVE_PROFILES = PROP_PROFILES_ACTIVE
)

type argsConfig struct {
	ordinaryPrefix prefixHolder
	prefix prefixHolder
}
func (self *argsConfig) setPrefix(prefix string) *argsConfig {
	self.ordinaryPrefix = prefixHolder(prefix)
	self.prefix = self.ordinaryPrefix.turnDash()
	return self
}
func (self *argsConfig) bindByPflag(flagSet *pflag.FlagSet) *packedConfig {
	packedConfigByArgs := &packedConfig{}

	/**
	 * Loads configuration by arguments of command
	 */
	flagSet.StringVar(&packedConfigByArgs.jsonProps,
		self.prefixWith(FLAG_CONFIG_JSON), "",
		`'{ "key": value }'`,
	)
	flagSet.StringVar(&packedConfigByArgs.yamlProps,
		self.prefixWith(FLAG_CONFIG_YAML), "",
		`'{ key: value }'`,
	)
	flagSet.StringVar(&packedConfigByArgs.externalFiles,
		self.prefixWith(FLAG_CONFIG_FILES), "",
		`'<file path>'`,
	)
	flagSet.StringVar(&packedConfigByArgs.activeProfiles,
		self.prefixWith(FLAG_ACITVE_PROFILES), "",
		`'profile1,profile2'`,
	)
	// :~)

	packedConfigByArgs.prefix = self.ordinaryPrefix
	return packedConfigByArgs
}
func (self *argsConfig) prefixWith(suffix string) string {
	return fmt.Sprintf("%s%s", self.prefix, suffix)
}
