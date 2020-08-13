/*
Environment

An environment a container of properties and profile for your application.

You can construct "Environment" by "EnvBuilder.NewByMap" or "EnvBuilder.NewByViper".

Profile

The profile of application is defined by property of "frangipani.profiles.active".

For example(by "viper"):

  your_app -frangipani.profiles.active=cus1,cus2

*/
package frangipani

import (
	"strings"
	"github.com/spf13/viper"
)

// Property name for "frangipani.profiles.active"
const PROP_ACITVE_PROFILES = "frangipani.profiles.active"

// "default" value of profile
const DEFAULT_PROFILE = "default"

// Global space to construct "Environment"
var EnvBuilder IEnvBuilder = IEnvBuilder(0)

type IEnvBuilder int
// Constructs environment by multiple objects of "*viper.Viper".
//
// The order of variadic arguments would be the priority of loading of configurations.
func (IEnvBuilder) NewByVipers(vipers ...*viper.Viper) Environment {
	return newEnvViperImpl(vipers...)
}

// Constructs environment by a "*viper.Viper".
func (self IEnvBuilder) NewByViper(viper *viper.Viper) Environment {
	return self.NewByVipers(viper)
}

// Constructs new environment by "map[string]interface{}".
func (IEnvBuilder) NewByMap(props map[string]interface{}) Environment {
	newEnv := &mapBasedEnv{ PropertyResolver: PropertyResolverBuilder.NewByMap(props) }
	newEnv.activeProfiles = newEnv.processActiveProfiles()
	return newEnv
}

// A conceptual container to gain properties and profiles of a application.
type Environment interface {
	// Getter of properties
	PropertyResolver

	// Checks whether or not the profiles is enabled in current application.
	AcceptsProfiles(profiles Profiles) bool
	// Gets current active profiles.
	GetActiveProfiles() []string
}

type mapBasedEnv struct {
	PropertyResolver

	activeProfiles []string
}
func (self *mapBasedEnv) AcceptsProfiles(profiles Profiles) bool {
	return profiles.Matches(self.matchProfile)
}
func (self *mapBasedEnv) GetActiveProfiles() []string {
	profiles := make([]string, len(self.activeProfiles))
	copy(profiles, self.activeProfiles)
	return profiles
}
func (self *mapBasedEnv) processActiveProfiles() []string {
	profiles := make([]string, 1)
	profiles[0] = DEFAULT_PROFILE

	profilesByProperty := strings.Split(self.GetProperty(PROP_ACITVE_PROFILES), ",")
	for _, profile := range profilesByProperty {
		profile = strings.TrimSpace(profile)

		if profile != "" {
			profiles = append(profiles, profile)
		}
	}

	return profiles
}
func (self *mapBasedEnv) matchProfile(profile string) bool {
	for _, activeProfile := range self.activeProfiles {
		if activeProfile == profile {
			return true
		}
	}

	return false
}
