package env

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"

	"github.com/spf13/viper"
)

type fileNamesBuilder struct {
	prefix string
}
func (self *fileNamesBuilder) getFiles() []string {
	files := make([]string, 0, len(default_suffix_names))

	prefix := self.prefix
	for _, suffix := range default_suffix_names {
		files = append(files, fmt.Sprintf("%s%s", prefix, suffix))
	}

	return files
}
func (self *fileNamesBuilder) toConfigFileNames(logErrorNotFound bool) *configFileNames {
	return &configFileNames {
		names: self.getFiles(),
		logErrorNotFound: logErrorNotFound,
	}
}

func newConfigFiles(logError bool, filenames ...string) *configFileNames {
	return &configFileNames {
		filenames, logError,
	}
}

// With values of names, provides method to load these files
type configFileNames struct {
	names []string
	logErrorNotFound bool
}
func (self *configFileNames) load() vipers {
	return self.loadByDir("")
}
func (self *configFileNames) loadWithProfiles(profiles ...string) vipers {
	return self.loadByDirWithProfiles("", profiles...)
}
func (self *configFileNames) loadByDir(directory string) vipers {
	propsByFiles := make(vipers, 0)

	for _, fileName := range self.names {
		fileName = filenameDir(fileName, directory)
		if viperObj, ok := readInByFile(fileName, self.logErrorNotFound); ok {
			propsByFiles = append(propsByFiles, viperObj)
		}
	}

	return propsByFiles
}
func (self *configFileNames) loadByDirWithProfiles(directory string, profiles ...string) vipers {
	propsByFiles := make(vipers, 0)

	for _, fileName := range self.names {
		for _, profiledFile := range profiledFiles(fileName, profiles...) {
			profiledFile = filenameDir(profiledFile, directory)
			if viperObj, ok := readInByFile(profiledFile, self.logErrorNotFound); ok {
				propsByFiles = append(propsByFiles, viperObj)
			}
		}
	}

	return propsByFiles
}

// Builds <dir>/<file> if the directory is not empty
func filenameDir(fileName string, directory string) string {
	if directory == "" {
		return fileName
	}

	return fmt.Sprintf("%s/%s", directory, fileName)
}

// Builds file names for for profiles(name>-<profile>.<ext>)
func profiledFiles(fileName string, profiles ...string) []string {
	result := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		fileExt := filepath.Ext(fileName)
		baseNameWithoutExt := filepath.Base(fileName)

		if fileExt != "" {
			lastIndexOfExt := strings.LastIndex(baseNameWithoutExt, ".")
			baseNameWithoutExt = baseNameWithoutExt[:lastIndexOfExt]
		}

		result = append(result, fmt.Sprintf("%s-%s%s", baseNameWithoutExt, profile, fileExt))
	}

	return result
}

// Initializes viper by external file
func readInByFile(file string, warnNotFound bool) (*viper.Viper, bool) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			if warnNotFound {
				configLogger.Warnf("file could not be found: [%s]", file)
			} else {
				configLogger.Debugf("file could not be found: [%s]", file)
			}
		} else {
			configLogger.Warn(fmt.Errorf("Stat file has error: %w", err))
		}

		return nil, false
	}

	viperByFile := viper.New()
	viperByFile.SetConfigFile(file)
	if err = viperByFile.ReadInConfig(); err != nil {
		configLogger.Warn(fmt.Errorf("Read file has error: %v", err))
		return viperByFile, false
	} else {
		configLogger.Infof("Config file loaded: [%s]", file)
	}

	return viperByFile, true
}
