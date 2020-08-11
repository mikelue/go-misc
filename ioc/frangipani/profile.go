package frangipani

import (
	"strings"
)

// Provides matching with callback of predicate.
type Profiles interface {
	// Iterates each of checked profiles by fed callback
	Matches(func(string) bool) bool
}

// Initializes the profiles as "Profiles"
func OfProfiles(profiles ...string) Profiles {
	newProfiles := make([]string, 0, len(profiles))

	for _, profile := range profiles {
		profile = strings.TrimSpace(profile)
		if profile != "" {
			newProfiles = append(newProfiles, profile)
		}
	}

	return ofProfilesImpl(newProfiles)
}

type ofProfilesImpl []string
func (self ofProfilesImpl) Matches(matchFunc func(string) bool) bool {
	if len(self) == 0 {
		return true
	}

	for _, profile := range self {
		if !matchFunc(profile) {
			return false
		}
	}

	return true
}
