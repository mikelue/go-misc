package env

import (
	"strings"

	"github.com/spf13/viper"
)

// As collection of supported source
//
// From a content as JSON format
// From a content as YAML format
// From a name(path) of file
// From a content as active profiles
type packedConfig struct {
	prefix prefixHolder
	jsonProps string
	yamlProps string
	externalFiles string
	activeProfiles string

	loadedVipers vipers
}
func (self *packedConfig) loadFormattedProps() vipers {
	loadedVipers := make(vipers, 0)

	if self.yamlProps != "" {
		configLogger.Debugf("Reading properties by JSON")
		if viperObj, ok := readInByString(
			FLAG_CONFIG_YAML, "yaml", self.yamlProps,
		); ok {
			loadedVipers = append(loadedVipers, viperObj)
		}
	}

	if self.jsonProps != "" {
		configLogger.Debugf("Reading properties by YAML")
		if viperObj, ok := readInByString(
			FLAG_CONFIG_JSON, "json", self.jsonProps,
		); ok {
			loadedVipers = append(loadedVipers, viperObj)
		}
	}

	return loadedVipers
}
// This only build the property value "kyc.config.file" to value of "externalFiles".
func (self *packedConfig) configFileProp() *viper.Viper {
	if self.externalFiles == "" {
		return nil
	}

	filesAsSlice := make([]string, 0, 1)
	for _, fileName := range strings.Split(self.externalFiles, ",") {
		fileName = strings.TrimSpace(fileName)
		if fileName != "" {
			filesAsSlice = append(filesAsSlice, fileName)
		}
	}

	viperObj := viper.New()
	viperObj.Set(
		self.prefixKey(PROP_CONFIG_FILES),
		filesAsSlice,
	)
	return viperObj
}
func (self *packedConfig) load() vipers {
	if self.loadedVipers != nil {
		return self.loadedVipers
	}

	finalVipers := make(vipers, 0)

	if activeProfiles := self.loadActiveProfiles(); activeProfiles != nil {
		finalVipers = append(finalVipers, activeProfiles)
	}

	finalVipers = append(finalVipers, self.loadFormattedProps()...)

	if externalFilesViper := self.configFileProp(); externalFilesViper != nil {
		finalVipers = append(finalVipers, externalFilesViper)
	}

	self.loadedVipers = finalVipers
	return self.loadedVipers
}
func (self *packedConfig) loadWithProfiles(profiles ...string) vipers {
	return self.load()
}
func (self *packedConfig) loadActiveProfiles() *viper.Viper {
	if self.activeProfiles == "" {
		return nil
	}

	configLogger.Debugf("Got active profiles of viable")
	activeProfiles := viper.New()

	/**
	 * Puts for both of "prefixed" property, and profile property of Frangipani
	 */
	activeProfiles.Set(
		self.prefixKey(PROP_PROFILES_ACTIVE),
		self.activeProfiles,
	)
	// :~)

	return activeProfiles
}
func (self *packedConfig) prefixKey(key string) string {
	return string(self.prefix.withSuffix(key))
}

// Initializes viper by content of string
func readInByString(prop string, contentType string, content string) (*viper.Viper, bool) {
	viperObj := viper.New()
	viperObj.SetConfigType(contentType)

	if err := viperObj.ReadConfig(strings.NewReader(content)); err != nil {
		configLogger.Warnf(`Read property "%s" has error(should be YAML): %v`, prop, err)
		return nil, false
	}

	return viperObj, true
}
