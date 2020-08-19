package utils

import (
	"os"
	"fmt"
)

func ExampleIRollbackContainerBuilder_newTempDir() {
	tempDir := RollbackContainerBuilder.NewTmpDir("proj1-*")

	var createdDir string
	RollbackExecutor.RunP(
		func(params Params) {
			// Access name of create directory of temporary by "params[PKEY_TEMP_DIR]"
			createdDir = params[PKEY_TEMP_DIR].(string)
			stat, err := os.Stat(createdDir)
			fmt.Printf("Name of temporary directory: %v %v\n", stat.IsDir(), err == nil)
		},
		tempDir,
	)

	_, err := os.Stat(createdDir)
	fmt.Printf("Temporary directory is removed: %v", os.IsNotExist(err))

	// Output:
	// Name of temporary directory: true true
	// Temporary directory is removed: true
}

func ExampleIRollbackContainerBuilder_newEnv() {
	changeEnv := RollbackContainerBuilder.NewEnv(
		map[string]string {
			"MYENV_1": "hello",
			"MYENV_2": "world",
		},
	)

	RollbackExecutor.Run(
		func() {
			v1, _ := os.LookupEnv("MYENV_1")
			v2, _ := os.LookupEnv("MYENV_2")
			fmt.Printf("Environment: %s %s\n", v1, v2)
		},
		changeEnv,
	)

	_, ok1 := os.LookupEnv("MYENV_1")
	_, ok2 := os.LookupEnv("MYENV_2")
	fmt.Printf("Reset: %v %v", ok1, ok2)

	// Output:
	// Environment: hello world
	// Reset: false false
}

func ExampleIRollbackContainerBuilder_newDir() {
	tempDirName := fmt.Sprintf("%s/example-newdir", os.TempDir())
	newDir := RollbackContainerBuilder.NewDir(tempDirName)

	RollbackExecutor.Run(
		func() {
			stat, err := os.Stat(tempDirName)
			fmt.Printf("Name of new directory: %v %v\n", stat.IsDir(), err == nil)
		},
		newDir,
	)

	_, err := os.Stat(tempDirName)
	fmt.Printf("Directory is removed: %v\n", os.IsNotExist(err))

	// Output:
	// Name of new directory: true true
	// Directory is removed: true
}

func ExampleIRollbackContainerBuilder_newChdir() {
	tempDirName := fmt.Sprintf("%s/chdir-temp", os.TempDir())
	chDir := RollbackContainerBuilder.NewChdir(tempDirName)

	RollbackExecutor.Run(
		func() {
			workingDir, _ := os.Getwd()
			fmt.Printf("Working directory is same as created: %v\n",
				workingDir == tempDirName,
			)
		},
		// Creates the sample directory first
		RollbackContainerBuilder.NewDir(tempDirName),
		chDir,
	)

	workingDir, _ := os.Getwd()
	fmt.Printf("Working directory is changed back: %v\n",
		workingDir != tempDirName,
	)

	// Output:
	// Working directory is same as created: true
	// Working directory is changed back: true
}

func ExampleIRollbackContainerBuilder_newCopyfiles() {
	tempDirName := fmt.Sprintf("%s/copy-files-temp", os.TempDir())

	destFileName1 := fmt.Sprintf("%s/sample-1.txt", tempDirName)
	destFileName2 := fmt.Sprintf("%s/sample-2.txt", tempDirName)

	copyFiles := RollbackContainerBuilder.NewCopyFiles(
		tempDirName,
		fmt.Sprintf("%s/sample-1.txt", testSourceDir),
		fmt.Sprintf("%s/sample-2.txt", testSourceDir),
	)

	RollbackExecutor.Run(
		func() {
			_, err1 := os.Stat(destFileName1)
			_, err2 := os.Stat(destFileName1)

			fmt.Printf("[Existing] File 1: %v. File 2: %v\n",
				err1 == nil, err2 == nil,
			)
		},
		// Creates the sample directory first
		RollbackContainerBuilder.NewDir(tempDirName),
		copyFiles,
	)

	_, err1 := os.Stat(destFileName1)
	_, err2 := os.Stat(destFileName2)
	fmt.Printf("[Removed] File 1: %v. File 2: %v\n",
		os.IsNotExist(err1), os.IsNotExist(err2),
	)

	// Output:
	// [Existing] File 1: true. File 2: true
	// [Removed] File 1: true. File 2: true
}
