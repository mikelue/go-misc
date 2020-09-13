/*
Built-in mechanism for loading configurations(properties) from various source.

ConfigLoader

You can use "DefaultLoader" to load configuration by default.

  evn := DefaultLoader.New().
    ParseFlags().
    Load()

Documentation: https://github.com/mikelue/go-misc/blob/master/ioc/frangipani/README.md

ConfigBuilder

You can use this object to customized the behavior of loading configurations.

  env := NewConfigBuilder().
    Priority(CL_ENVVAR, CL_ARGS, CL_CONFIG_FILE, CL_PWD).
    Prefix("dog").
    Build().
    ParseFlags().Load()
*/
package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	fg "github.com/mikelue/go-misc/ioc/frangipani"

	l4 "github.com/go-eden/slf4go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

const (
	// Default prefix of loading mechanism
	DEFAULT_PREFIX = "fgapp"

	// Default name of logger
	//
	// See: https://github.com/go-eden/slf4go-logrus
	LOGGER_NAME_CONFIG = "fgapp.config"

	// The configuration file
	PROP_CONFIG_FILES = ".config.files"
	// The profiles to be activated
	PROP_PROFILES_ACTIVE = ".profiles.active"
)

// Method space used to construct new instance of "ConfigLoader"
var DefaultLoader IDefaultLoader

type IDefaultLoader int
// With default setting for loading configurations
func (self *IDefaultLoader) New() ConfigLoader {
	return self.WithViper(nil)
}
// With some default values(as "*viper.Viper") and settings for loading configurations
func (*IDefaultLoader) WithViper(viper *viper.Viper) ConfigLoader {
	return NewConfigBuilder().
		DefaultWithViper(viper).
		Build()
}
// With some default values(as map) and settings for loading configurations
func (*IDefaultLoader) WithMap(properties map[string]interface{}) ConfigLoader {
	return NewConfigBuilder().
		DefaultWithMap(properties).
		Build()
}

// Consturcts a new builder so that you can custoimze the loading of configurations.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder {
		flags: pflag.CommandLine,
		prefix: DEFAULT_PREFIX,
		sources: []ConfigSource {
			CL_XDG, CL_ARGS, CL_ENVVAR, CL_CONFIG_FILE,
			CL_PWD, CL_CMDDIR,
		},
	}
}

// To load a bunch of sources as "Environment".
type ConfigLoader interface {
	// Loads the environment object
	Load() fg.Environment
	// Parse the flags
	ParseFlags() ConfigLoader
}

// CL - stands for [C]onfiguration [L]oading
type ConfigSource int
const (
	// The source comes from default names of files in $XDG_CONFIG_HOME
	CL_XDG ConfigSource = 1
	// The source comes from arguments:
	//   --fgapp.config.yaml
	//   --fgapp.config.json
	//   --fgapp.config.config.files
	//   --fgapp.config.files
	//   --fgapp.profiles.active
	CL_ARGS ConfigSource = 2
	// The source comes from environment variables:
	//   $FGAPP_CONFIG_YAML
	//   $FGAPP_CONFIG_JSON
	//   $FGAPP_CONFIG_FILES
	//   $FGAPP_PROFILES_ACTIVE
	CL_ENVVAR ConfigSource = 3
	// The source comes from provided file name
	CL_CONFIG_FILE ConfigSource = 4
	// The source comes default names of files in working directory
	//
	// See: os.Getwd()
	CL_PWD ConfigSource = 5
	// The source comes default names of files in directory of command
	//
	// See: os.Args[0]
	CL_CMDDIR ConfigSource = 6
)

var prefixPattern, _ = regexp.Compile("[a-zA-Z][a-zA-Z0-9-_]+")
// This object is the center of customizing loading mechanisms of configurations.
type ConfigBuilder struct {
	flags *pflag.FlagSet
	prefix prefixHolder
	defaultValues *viper.Viper
	sources []ConfigSource
	workers map[ConfigSource]loadingWorker
}
// Sets the priority of supported sources.
//
// The higher priority of sources should be put in front of others.
//
// See: ConfigSource
func (self *ConfigBuilder) Priority(sources ...ConfigSource) *ConfigBuilder {
	newSourcesMap := make(map[ConfigSource]bool, len(sources))
	newSources := make([]ConfigSource, 0, len(sources))

	for _, source := range sources {
		if _, ok := newSourcesMap[source]; ok {
			configLogger.Warnf("Duplicated config source: %v", source)
			continue
		}

		newSourcesMap[source] = true
		newSources = append(newSources, source)
	}

	newSources = eliminateSameDirForWdAndCmd(newSources)
	self.sources = newSources
	return self
}
// Sets the prefix of loading.
//
// This can affect naming of various sources.
//
// See detail: https://github.com/mikelue/go-misc/blob/master/ioc/frangipani/README.md
func (self *ConfigBuilder) Prefix(prefix string) *ConfigBuilder {
	if !prefixPattern.MatchString(prefix) {
		configLogger.Warnf("Prefix is not valid: %v", prefix)
		return self
	}

	self.prefix = prefixHolder(prefix)
	return self
}
// Sets-up customized flags for parsing of arguments
func (self *ConfigBuilder) Pflags(newFlags *pflag.FlagSet) *ConfigBuilder {
	self.flags = newFlags
	return self
}
// Sets up the default properties, this has the lowest priority set by "Priority".
func (self *ConfigBuilder) DefaultWithMap(properties map[string]interface{}) *ConfigBuilder {
	if properties == nil {
		return self.DefaultWithViper(nil)
	}

	newViper := viper.New()

	if err := newViper.MergeConfigMap(properties); err != nil {
		configLogger.Errorf("Unable to set default values by map. Error: %v", err)
		return self
	}
	self.DefaultWithViper(newViper)

	return self
}
// Sets up the default properties(by viper), this has the lowest priority set by "Priority".
//
// See spf13/Viper: https://github.com/spf13/viper
func (self *ConfigBuilder) DefaultWithViper(viper *viper.Viper) *ConfigBuilder {
	self.defaultValues = viper
	return self
}
// Builds the "ConfigLoader", which is used to load an instance of "Environment".
func (self *ConfigBuilder) Build() ConfigLoader {
	newLoader := &configLoaderImpl{ ConfigBuilder: self }
	return newLoader.init()
}

func init() {
	DefaultLoader = IDefaultLoader(0)
	l4.SetLevel(l4.InfoLevel)
}

type configLoaderImpl struct {
	*ConfigBuilder

	argsConfig *packedConfig

	hasConfigFile bool
}
func (self *configLoaderImpl) Load() fg.Environment {
	pass1Env := self.pass1Load()
	return self.pass2Load(pass1Env)
}
func (self *configLoaderImpl) ParseFlags() ConfigLoader {
	envArgs := make([]string, 0, 0)

	envPrefix := fmt.Sprintf("--%s.", self.prefix.turnDash())
	for _, v := range os.Args {
		if strings.HasPrefix(v, envPrefix) {
			if configLogger.IsDebugEnabled() {
				configLogger.Debugf("Found argv: \"%v\"",
					v[:strings.Index(v, "=")])
			}
			envArgs = append(envArgs, v)
		}
	}

	self.flags.Parse(envArgs)
	return self
}
func (self *configLoaderImpl) init() *configLoaderImpl {
	for _, source := range self.sources {
		if source == CL_ARGS {
			self.argsConfig = new(argsConfig).
				setPrefix(string(self.prefix)).
				bindByPflag(self.flags)
			break
		}
	}

	return self
}
func (self *configLoaderImpl) pass1Load() fg.Environment {
	workersMap := make(map[ConfigSource]loadingWorker, len(self.sources))
	allVipers := make(vipers, 0, len(self.sources))

	/**
	 * Builds the vipers
	 */
	for _, source := range self.sources {
		var worker loadingWorker

		switch source {
		case CL_XDG, CL_ENVVAR, CL_PWD, CL_CMDDIR:
			worker = workerFactories[source](string(self.prefix))
		case CL_ARGS:
			worker = self.argsConfig
		case CL_CONFIG_FILE:
			emptyLoader := emptyLoader(0)
			worker = &emptyLoader
			self.hasConfigFile = true
		default:
			configLogger.Warnf("Unsupported type of source: %v", source)
			continue
		}

		allVipers = append(allVipers, worker.load()...)
		workersMap[source] = worker
	}
	// :~)

	/**
	 * Appends the default values to last viper
	 */
	if self.defaultValues != nil {
		allVipers = append(allVipers, self.defaultValues)
	}
	// :~)

	/**
	 * Sets the property for active profiles
	 */
	profiles := self.getProfiles(allVipers)
	if profiles != "" {
		allVipers[0].Set(fg.PROP_ACITVE_PROFILES, profiles)
	}
	// :~)


	self.workers = workersMap
	return fg.EnvBuilder.NewByVipers(allVipers...)
}
func (self *configLoaderImpl) pass2Load(env fg.Environment) fg.Environment {
	/**
	 * Loads additional files of configuration
	 */
	if self.hasConfigFile {
		configFiles := env.Typed().GetStringSlice(
			string(self.prefix.withSuffix(PROP_CONFIG_FILES)),
		)
		if len(configFiles) > 0 {
			self.workers[CL_CONFIG_FILE] = workerBuilder.newFiles(configFiles...)
		}
	}
	// :~)

	profiles := env.GetActiveProfiles()

	/**
	 * Loads vipers with profile
	 */
	allVipers := make(vipers, 0, len(self.sources))
	for _, source := range self.sources {
		var loadedVipers vipers

		switch source {
		// Only these sources support profiles
		case CL_XDG, CL_PWD, CL_CMDDIR:
			loadedVipers = self.workers[source].loadWithProfiles(profiles...)
		default:
			loadedVipers = self.workers[source].load()
		}

		allVipers = append(allVipers, loadedVipers...)
	}
	// :~)

	if self.defaultValues != nil {
		allVipers = append(allVipers, self.defaultValues)
	}

	return fg.EnvBuilder.NewByVipers(allVipers...)
}
func (self *configLoaderImpl) getProfiles(sources []*viper.Viper) string {
	profilesProp := self.prefix.withSuffix(PROP_PROFILES_ACTIVE)

	for _, viper := range sources {
		if profiles := viper.GetString(string(profilesProp)); profiles != "" {
			return profiles;
		}
	}

	return ""
}

var configLogger = l4.NewLogger(LOGGER_NAME_CONFIG)

var workerFactories = map[ConfigSource]workerFactory {
	CL_XDG: workerBuilder.newXdg,
	CL_ENVVAR: func(prefix string) loadingWorker {
		return new(envConfig).
			setPrefix(prefix).
			load()
	},
	CL_PWD: workerBuilder.newWd,
	CL_CMDDIR: workerBuilder.newCmdDir,
}

type loadingWorker interface {
	// Used by 1-pass
	load() vipers
	// Used by 2-pass
	loadWithProfiles(profiles ...string) vipers
}

type emptyLoader int
func (*emptyLoader) load() vipers {
	return nil
}
func (self *emptyLoader) loadWithProfiles(...string) vipers {
	return self.load()
}

type workerFactory func(prefix string) loadingWorker

type vipers []*viper.Viper
func (self vipers) appendMore(tails ...*viper.Viper) vipers {
	appendedVipers := self

	for _, viper := range tails {
		appendedVipers = append(appendedVipers, viper)
	}

	return appendedVipers
}

type prefixHolder string
func (self *prefixHolder) withSuffix(suffix string) prefixHolder {
	return prefixHolder(fmt.Sprintf("%s%s", *self, suffix))
}
func (self *prefixHolder) turnDash() prefixHolder {
	return prefixHolder(strings.ReplaceAll(string(*self), "-", "_"))
}
func (self *prefixHolder) turnDashAndUpper() prefixHolder {
	stringValue := strings.ReplaceAll(
		strings.ToUpper(string(*self)),
		"-", "_",
	)

	return prefixHolder(stringValue)
}

func eliminateSameDirForWdAndCmd(sources []ConfigSource) []ConfigSource {
	var wdIndex, cmdDirIndex int = -1, -1

	for i, source := range sources {
		if source == CL_PWD {
			wdIndex = i
		}
		if source == CL_CMDDIR {
			cmdDirIndex = i
		}
	}

	/**
	 * Nothing changed
	 */
	if wdIndex == -1 || cmdDirIndex == -1 {
		return sources
	}

	// If the working directory is not same as directory of command
	wd, cmdDir := getWd(), getCmdDir()
	if wd != cmdDir {
		return sources
	}
	// :~)

	/**
	 * Removes the lowest priority between CL_PWD and CL_CMDDIR
	 */
	beingRemovedSource := CL_CMDDIR
	if cmdDirIndex < wdIndex {
		beingRemovedSource = CL_PWD
	}

	return funk.Filter(
		sources,
		func(s ConfigSource) bool {
			return s != beingRemovedSource
		},
	).([]ConfigSource)
	// :~)
}
