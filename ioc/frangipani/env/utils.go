package env

import (
	"os"
	"path/filepath"
)

func getWd() string {
	workingDir, err := os.Getwd()

	if err != nil {
		configLogger.Errorf("Get working directory has failed: %v", err)
		return ""
	}

	return workingDir
}
func getCmdDir() string {
	dirOfCmd, err := filepath.Abs(os.Args[0])

	if err != nil {
		configLogger.Errorf("Get directory of command has failed: %v", err)
		return ""
	}

	return filepath.Dir(dirOfCmd)
}
