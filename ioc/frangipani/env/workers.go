package env

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var workerBuilder = workerBuilderI(0)

type workerBuilderI int
func (*workerBuilderI) newXdg(prefix string) loadingWorker {
	newWorker := newDefaultFilesWorker(prefix)

	dir, err := os.UserConfigDir()
	if err != nil {
		return newWorker
	}

	configDirInXdg := fmt.Sprintf("%s/%s", dir, prefix)
	if _, err = os.Stat(configDirInXdg); err != nil {
		return newWorker
	}

	configLogger.Debugf("Found XDG dir: %s", configDirInXdg)
	newWorker.targetDir = configDirInXdg
	return newWorker
}
func (*workerBuilderI) newWd(prefix string) loadingWorker {
	newWorker := newDefaultFilesWorker(prefix)

	workingDir := getWd()
	if workingDir != "" {
		configLogger.Debugf("Found working directory: %s", workingDir)
	}

	newWorker.targetDir = workingDir
	return newWorker
}
func (*workerBuilderI) newCmdDir(prefix string) loadingWorker {
	dirOfCmd := getCmdDir()

	if dirOfCmd != "" {
		configLogger.Debugf("Found directory of command: %s", dirOfCmd)
	}

	newWorker := newDefaultFilesWorker(prefix)
	newWorker.targetDir = dirOfCmd
	return newWorker
}
func (*workerBuilderI) newFiles(fileNames ...string) loadingWorker {
	return &filesWorker {
		files: fileNames,
	}
}

var default_suffix_names = []string {
	"-config.properties", "-config.yaml", "-config.json",
}

func newDefaultFilesWorker(prefix string) *filesWorker {
	fileNames := make([]string, 0, len(default_suffix_names))

	for _, fileSuffix := range default_suffix_names {
		fileNames = append(fileNames, fmt.Sprintf("%s%s", prefix, fileSuffix))
	}

	return &filesWorker {
		files: fileNames,
	}
}

type filesWorker struct {
	files []string
	targetDir string
	// Cached vipers
	loadedFiles []*viper.Viper
}
func (self *filesWorker) load() vipers {
	if self.loadedFiles != nil {
		return self.loadedFiles
	}

	configFiles := newConfigFiles(false, self.files...)
	self.loadedFiles = configFiles.loadByDir(self.targetDir)

	configLogger.Debugf("Found [%d] files.", len(self.loadedFiles))
	return self.loadedFiles
}
func (self *filesWorker) loadWithProfiles(profiles ...string) vipers {
	loadedFiles := self.load()

	if self.targetDir == "" || len(profiles) == 0 {
		return loadedFiles
	}

	filesWithProfiles := newConfigFiles(false, self.files...).
		loadByDirWithProfiles(self.targetDir, profiles...).
		appendMore(loadedFiles...)

	configLogger.Debugf("Found [%d] files(with profile).", len(filesWithProfiles))
	return filesWithProfiles
}
