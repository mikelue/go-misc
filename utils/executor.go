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

See examples of this documentation.

RollbackContainerBuilder

See examples of this documentation.
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
var RollbackExecutor IRollbackExecutor

type IRollbackExecutor int
// Run callback with simple containers
func (*IRollbackExecutor) Run(callback func(), containers ...RollbackContainer) error {
	containersP := make([]RollbackContainerP, 0, len(containers))

	for _, container := range containers {
		containersP = append(containersP, RollbackContainerBuilder.ToContainerP(container))
	}

	return rollbackExecutorImpl(containersP).Run(callback)
}
// Run callback with containers having parameters
func (*IRollbackExecutor) RunP(callback func(Params), containers ...RollbackContainerP) error {
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

// Empty rollback container, do nothing.
const EmptyRollbackContainer emptyRollbackContainer = 0

type emptyRollbackContainer int
func (emptyRollbackContainer) Setup() error {
	return nil
}
func (emptyRollbackContainer) TearDown() error {
	return nil
}

// Empty rollback container with parameters, do nothing.
const EmptyRollbackContainerP emptyRollbackContainerP = 0

type emptyRollbackContainerP int
func (emptyRollbackContainerP) Setup() (Params, error) {
	return nil, nil
}
func (emptyRollbackContainerP) TearDown(Params) (error) {
	return nil
}

// Method space to build new instances of executors
var RollbackContainerBuilder IRollbackContainerBuilder

type IRollbackContainerBuilder int

// Converts a "RollbackExecutorP" to a "RollbackExecutor"
func (*IRollbackContainerBuilder) ToContainer(containerP RollbackContainerP) RollbackContainer {
	return &fromPContainer{ containerP }
}

// Converts a "RollbackExecutor" to a "RollbackExecutorP"
func (*IRollbackContainerBuilder) ToContainerP(container RollbackContainer) RollbackContainerP {
	return &simpleContainerP{ container }
}

// Concatenate multiple "RollbackContainer"s as a single one.
//
// The following containers would not be Setup() if the setup of prior container hsa failed
//
// Only the containers of successful setting-up would be tear-down in reversed order.
func (self *IRollbackContainerBuilder) Concate(containers ...RollbackContainer) RollbackContainer {
	containersP := make([]RollbackContainerP, 0, len(containers))

	for _, container := range containers {
		containersP = append(containersP, self.ToContainerP(container))
	}

	concatContainersP := self.ConcateP(containersP...).(*concatRollbackContainersP)
	return &concatRollbackContainers {
		containersP: concatContainersP,
	}
}
// Concatenate multiple "RollbackContainerP"s as a single one.

// The following containers would not be Setup() if the setup of prior container hsa failed
//
// Only the containers of successful setting-up would be tear-down in reversed order.
//
// The parameters would be merged as arguments.
func (*IRollbackContainerBuilder) ConcateP(containers ...RollbackContainerP) RollbackContainerP {
	return &concatRollbackContainersP{ containers: containers }
}

// Constructs an new container with copy/removal of files.
//
// The copied files would be remove by the rollback.
func (*IRollbackContainerBuilder) NewCopyFiles(dir string, files ...string) RollbackContainer {
	return &copyFilesExecutorImpl { destDir: dir, srcFiles: files }
}

// Constructs an new container with creation/removal temp directory.
//
// The key of parameters for the temp directory is "PKEY_TEMP_DIR".
//
// The temp directory would be removed after the execution.
func (*IRollbackContainerBuilder) NewTmpDir(tempName string) RollbackContainerP {
	return &tempDirExecutorImpl{ tempName: tempName }
}

// Constructs an new container with creation/removal a directory.
//
// The directory would be removed after the execution.
func (*IRollbackContainerBuilder) NewDir(dir string) RollbackContainer {
	return &dirExecutorImpl{ dir }
}

// Constructs an new container with modified environment.
//
// The "os.Environ()" would be reverted to original status after the execution.
func (*IRollbackContainerBuilder) NewEnv(changedEnvVars map[string]string) RollbackContainer {
	return &envExecutorImpl{
		changedEnvVars,
		map[string]string{},
	}
}

// Constructs an new container with modified working directory.
//
// The "os.Chdir()" would be used to revert the working dictionary.
func (*IRollbackContainerBuilder) NewChdir(targetDir string) RollbackContainer {
	return &chdirExecutorImpl{
		targetDir: targetDir,
	}
}

func init() {
	RollbackExecutor = 0
	RollbackContainerBuilder = 0
}

type concatRollbackContainers struct {
	containersP *concatRollbackContainersP
	params Params
}
func (self *concatRollbackContainers) Setup() (err error) {
	self.params, err = self.containersP.Setup()
	return
}
func (self *concatRollbackContainers) TearDown() (err error) {
	err = self.containersP.TearDown(self.params)
	return
}

type concatRollbackContainersP struct {
	containers []RollbackContainerP
	workingContainers []RollbackContainerP
}
func (self *concatRollbackContainersP) Setup() (params Params, err error) {
	self.workingContainers = nil
	workingContainers := make([]RollbackContainerP, 0, len(self.containers))

	params = make(Params, len(self.containers))

	/**
	 * Stops setting-up of following containers if the prior one has failed
	 */
	for _, container := range self.containers {
		currentParams, setupErr := container.Setup()
		if setupErr != nil {
			err = setupErr
			break
		}

		for k, v := range currentParams {
			params[k] = v
		}

		// Keeps the containers of successful setting-up
		workingContainers = append(workingContainers, container)
	}
	// :~)

	self.workingContainers = workingContainers

	return
}
func (self *concatRollbackContainersP) TearDown(params Params) (err error) {
	workingContainers := self.workingContainers

	for i := len(workingContainers) - 1; i >= 0; i-- {
		if tearDownErr := workingContainers[i].TearDown(params); tearDownErr != nil {
			err = tearDownErr
		}
	}

	return
}

type copyFilesExecutorImpl struct {
	srcFiles []string
	destDir string
	copiedFiles []string
}
func (self *copyFilesExecutorImpl) Setup() (err error) {
	self.copiedFiles = make([]string, 0, len(self.srcFiles))

	for _, srcFile := range self.srcFiles {
		dstFile, copyErr := self.copyTo(srcFile)
		if copyErr != nil {
			err = copyErr
			continue
		}

		self.copiedFiles = append(self.copiedFiles, dstFile)
	}

	return
}
func (self *copyFilesExecutorImpl) TearDown() (err error) {
	for _, file := range self.copiedFiles {
		// Try to remove all of the files
		if rmErr := os.Remove(file); rmErr != nil {
			err = rmErr
		}
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
	err = os.MkdirAll(self.dirName, os.ModePerm)
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
func (self *envExecutorImpl) Setup() (err error) {
	/**
	 * Set environment variable and backup old value
	 */
	for k, v := range self.neededVars {
		if oldValue, ok := os.LookupEnv(k); ok {
			self.oldVars[k] = oldValue
		}

		if setErr := os.Setenv(k, v); setErr != nil {
			err = setErr
		}
	}
	// :~)

	return
}
func (self *envExecutorImpl) TearDown() (err error) {
	/**
	 * Reverts the old value of environment variable
	 */
	for k, v := range self.oldVars {
		// Try to revert each of modified env-variables
		if setErr := os.Setenv(k, v); setErr != nil {
			err = setErr
		}
	}
	// :~)

	/**
	 * Un-set the environment variables
	 */
	for k := range self.neededVars {
		if _, existedInOld := self.oldVars[k]; !existedInOld {
			// Try to un-set each of introduced env-variables
			if unsetErr := os.Unsetenv(k); unsetErr != nil {
				err = unsetErr
			}
		}
	}
	// :~)

	self.oldVars = map[string]string{}
	return
}

type rollbackExecutorPImpl []RollbackContainerP
func (self rollbackExecutorPImpl) Run(callback func(params Params)) (err error) {
	allParams := make(Params)
	lastSetupIndex := 0

	defer func() {
		if p := recover(); p != nil {
			if panicAsErr, ok := p.(error); ok {
				err = panicAsErr
			} else {
				err = fmt.Errorf("Panic of callback: %v", p)
			}
		}

		for i := lastSetupIndex - 1; i >= 0; i-- {
			if tearDownErr := self[i].TearDown(allParams); tearDownErr != nil {
				err = tearDownErr
				return
			}
		}
	}()

	for _, container := range self {
		currentParams := make(Params)
		currentParams, err = container.Setup()

		if err != nil {
			return err
		}

		for k, v := range currentParams {
			allParams[k] = v
		}
		lastSetupIndex++
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
