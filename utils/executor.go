/*
Rollback

While we are testing, there are something we like to set-up. For example:

	1. make a temporary directory
	2. copy some test files to that directory
	3. sets some environment variable

After we finish our test, the environment should be as same as nothing happened.

The "RollbackExecutor" can be used to execute your routine with "RollbackContainer"s.

This package provides some built-in containers for temp directory, environment variables, etc.

RollbackExecutor

Using this space to perform execution surrounding by some containers.

RollbackContainerBuilder

Using this space to constructs built-in containers.

  RollbackExecutor.Run(
    RollbackContainerBuilder.NewEnv(
      map[string]string {
		"XDG_CONFIG_HOME": "your_dir",
	  }
	),
    RollbackContainerBuilder.NewDir("another_dir"),
  )
*/
package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	PKEY_TEMP_DIR = "_u_temp_dir_"
)

type Params map[string]interface{}


// Method space to run callback with 1:M containers
var RollbackExecutor IRollbackExecutor = 0

type IRollbackExecutor int
// Run callback with simple containers
func (IRollbackExecutor) Run(callback func(), containers ...RollbackContainer) error {
	containersP := make([]RollbackContainerP, 0, len(containers))

	for _, container := range containers {
		containersP = append(containersP, RollbackContainerBuilder.ToContainerP(container))
	}

	return rollbackExecutorImpl(containersP).Run(callback)
}
// Run callback with containers having parameters
func (IRollbackExecutor) RunP(callback func(Params), containers ...RollbackContainerP) error {
	return rollbackExecutorPImpl(containers).
		Run(callback)
}

// The container is used to perform setup/teardown logic
type RollbackContainer interface {
	// This method sets-up your fabricated container
	Setup() error
	// This method gets called as defer way after the callback gets finised
	TearDown() error
}

// The container is used to perform setup/teardown logic with parameters
type RollbackContainerP interface {
	// This method sets-up your fabricated container
	Setup() (Params, error)
	// This method gets called as defer way after the callback gets finised
	TearDown(Params) error
}

// Method space to build new instances of executors
var RollbackContainerBuilder IRollbackExecBuilder = 0

type IRollbackExecBuilder int

// Converts a "RollbackExecutorP" to a "RollbackExecutor"
func (IRollbackExecBuilder) ToContainer(containerP RollbackContainerP) RollbackContainer {
	return &fromPContainer{ containerP }
}

// Converts a "RollbackExecutor" to a "RollbackExecutorP"
func (IRollbackExecBuilder) ToContainerP(container RollbackContainer) RollbackContainerP {
	return &simpleContainerP{ container }
}

// Constructs an new container with copy/removal of files.
//
// The copied files would be remove by the rollback.
func (self IRollbackExecBuilder) NewCopyFiles(dir string, files ...string) RollbackContainer {
	return &copyFilesExecutorImpl { destDir: dir, srcFiles: files }
}

// Constructs an new container with creation/removal temp directory.
//
// The key of parameters for the temp directory is "PKEY_TEMP_DIR".
//
// The temp directory would be removed after the execution.
func (self IRollbackExecBuilder) NewTmpDir(tempName string) RollbackContainerP {
	return &tempDirExecutorImpl{ tempName: tempName }
}

// Constructs an new container with creation/removal a directory.
//
// The directory would be removed after the execution.
func (self IRollbackExecBuilder) NewDir(dir string) RollbackContainer {
	return &dirExecutorImpl{ dir }
}

// Constructs an new container with modified environment.
//
// The "os.Environ()" would be reverted to original status after the execution.
func (self IRollbackExecBuilder) NewEnv(changedEnvVars map[string]string) RollbackContainer {
	return &envExecutorImpl{
		changedEnvVars,
		map[string]string{},
	}
}

// Constructs an new container with modified working directory.
//
// The "os.Chdir()" would be used to revert the working dictionary.
func (self IRollbackExecBuilder) NewChdir(targetDir string) RollbackContainer {
	return &chdirExecutorImpl{
		targetDir: targetDir,
	}
}

type copyFilesExecutorImpl struct {
	srcFiles []string
	destDir string
	copiedFiles []string
}
func (self *copyFilesExecutorImpl) Setup() (err error) {
	self.copiedFiles = make([]string, 0, len(self.srcFiles))

	for _, srcFile := range self.srcFiles {
		var dstFile string
		dstFile, err = self.copyTo(srcFile)
		if err == nil {
			self.copiedFiles = append(self.copiedFiles, dstFile)
		}
	}

	return
}
func (self *copyFilesExecutorImpl) TearDown() (err error) {
	for _, file := range self.copiedFiles {
		// Try to remove all of the files
		err = os.Remove(file)
	}

	self.copiedFiles = nil
	return
}
func (self *copyFilesExecutorImpl) copyTo(srcFile string) (dstFile string, err error) {
	baseFileName := filepath.Base(srcFile)
	dstFile = fmt.Sprintf("%s/%s", self.destDir, baseFileName)

	source, srcErr := os.Open(srcFile)
	if srcErr != nil {
		err = srcErr
		return
	}
	defer source.Close()

	destination, destErr := os.Create(dstFile)
	if destErr != nil {
		err = destErr
		return
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return
}

type dirExecutorImpl struct {
	dirName string
}
func (self *dirExecutorImpl) Setup() (err error) {
	err = os.MkdirAll(self.dirName, os.ModeDir)
	return
}
func (self *dirExecutorImpl) TearDown() (err error) {
	err = os.RemoveAll(self.dirName)
	return
}

type tempDirExecutorImpl struct {
	tempName string
	tempDir string
}
func (self *tempDirExecutorImpl) Setup() (params Params, err error) {
	params = make(Params)

	self.tempDir, err = ioutil.TempDir(os.TempDir(), self.tempName)
	if err != nil {
		return
	}

	params[PKEY_TEMP_DIR] = self.tempDir

	return
}
func (self *tempDirExecutorImpl) TearDown(params Params) (err error) {
	err = os.RemoveAll(self.tempDir)
	self.tempDir = ""
	return
}

type chdirExecutorImpl struct {
	targetDir string
	oldDir string
}
func (self *chdirExecutorImpl) Setup() (err error) {
	self.oldDir, err = os.Getwd()

	if self.oldDir != self.targetDir {
		err = os.Chdir(self.targetDir)
	}
	return
}
func (self *chdirExecutorImpl) TearDown() (err error) {
	if self.oldDir == self.targetDir {
		return
	}

	err = os.Chdir(self.oldDir)
	self.oldDir = ""
	return
}

type envExecutorImpl struct {
	neededVars map[string]string
	oldVars map[string]string
}
func (self *envExecutorImpl) Setup() error {
	/**
	 * Set environment variable and backup old value
	 */
	for k, v := range self.neededVars {
		if oldValue, ok := os.LookupEnv(k); ok {
			self.oldVars[k] = oldValue
		}

		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	// :~)

	return nil
}
func (self *envExecutorImpl) TearDown() (err error) {
	/**
	 * Reverts the old value of environment variable
	 */
	for k, v := range self.oldVars {
		// Try to revert each of modified env-variables
		err = os.Setenv(k, v)
	}
	// :~)

	/**
	 * Un-set the environment variables
	 */
	for k := range self.neededVars {
		if _, existedInOld := self.oldVars[k]; !existedInOld {
			// Try to un-set each of introduced env-variables
			err = os.Unsetenv(k)
		}
	}
	// :~)

	self.oldVars = map[string]string{}
	return
}

type rollbackExecutorPImpl []RollbackContainerP
func (self rollbackExecutorPImpl) Run(callback func(params Params)) (err error) {
	allParams := make(Params)

	defer func() {
		for i := len(self) - 1; i >= 0; i-- {
			err = self[i].TearDown(allParams)
			if err != nil {
				return
			}
		}

		if p := recover(); p != nil {
			if panicAsErr, ok := p.(error); ok {
				err = panicAsErr
			} else {
				err = fmt.Errorf("Panic of callback: %v", p)
			}
		}
	}()

	for _, container := range self {
		currentParams := make(Params)
		currentParams, err = container.Setup()

		if err != nil {
			return
		}

		for k, v := range currentParams {
			allParams[k] = v
		}
	}

	callback(allParams)
	return
}

// Re-uses the implementation of "rollbackExecutorPImpl"
type rollbackExecutorImpl []RollbackContainerP
func (self rollbackExecutorImpl) Run(callback func()) error {
	return rollbackExecutorPImpl(self).Run(func(Params) {
		callback()
	})
}

// Converts a "RollbackContainer" to a "RollbackContainerP"
type simpleContainerP struct {
	container RollbackContainer
}
func (self *simpleContainerP) Setup() (Params, error) {
	return nil, self.container.Setup()
}
func (self *simpleContainerP) TearDown(Params) error {
	return self.container.TearDown()
}

type fromPContainer struct {
	container RollbackContainerP
}
func (self *fromPContainer) Setup() (err error) {
	_, err = self.container.Setup()
	return
}
func (self *fromPContainer) TearDown() error {
	return self.container.TearDown(nil)
}
