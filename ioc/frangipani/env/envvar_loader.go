package env

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// Name of environment variable for packed properties as JSON format
	ENVVAR_JSON = "_CONFIG_JSON"
	// Name of environment variable for packed properties as YAML format
	ENVVAR_YAML = "_CONFIG_YAML"
	// Name of environment variable for external file
	ENVVAR_FILE = "_CONFIG_FILES"
	// Name of environment variable for activated profiles(see frangipani)
	ENVVAR_PROFILES_ACTIVE = "_PROFILES_ACTIVE"
)

type envConfig struct {
	ordinaryPrefix prefixHolder
	prefix prefixHolder
}
func (self *envConfig) setPrefix(prefix string) *envConfig {
	self.ordinaryPrefix = prefixHolder(prefix)
	prefix = strings.ReplaceAll(
		strings.ToUpper(prefix),
		"-", "_",
	)

	self.prefix = self.ordinaryPrefix.turnDashAndUpper()
	return self
}
func (self *envConfig) load() *packedConfig {
	envViper := viper.New()

	self.bind(envViper, ENVVAR_JSON)
	self.bind(envViper, ENVVAR_YAML)
	self.bind(envViper, ENVVAR_FILE)
	self.bind(envViper, ENVVAR_PROFILES_ACTIVE)

	packedConfigByEnv := &packedConfig{}
	packedConfigByEnv.prefix = self.ordinaryPrefix
	packedConfigByEnv.jsonProps = self.getBySuffix(envViper, ENVVAR_JSON)
	packedConfigByEnv.yamlProps = self.getBySuffix(envViper, ENVVAR_YAML)
	packedConfigByEnv.externalFiles = self.getBySuffix(envViper, ENVVAR_FILE)
	packedConfigByEnv.activeProfiles = self.getBySuffix(envViper, ENVVAR_PROFILES_ACTIVE)

	return packedConfigByEnv
}
func (self *envConfig) prefixWith(name string) string {
	return fmt.Sprintf("%s%s", self.prefix, name)
}
func (self *envConfig) getBySuffix(viper *viper.Viper, suffix string) string {
	value := viper.GetString(suffix)

	if value != "" && configLogger.IsDebugEnabled() {
		envVarName := fmt.Sprintf("%s%s", self.prefix, suffix)
		configLogger.Debugf(`env "$%s" is read: `, envVarName)
	}

	return value
}
func (self *envConfig) bind(viper *viper.Viper, suffix string) {
	envVarName := fmt.Sprintf("%s%s", self.prefix, suffix)
	viper.BindEnv(suffix, envVarName)
}
