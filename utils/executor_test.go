package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rollback executors", func() {
	Describe("By IRollbackExecBuilder", describeOfIRollbackExecBuilder)
})

var describeOfIRollbackExecBuilder = func() {
	It("NewDir", func() {
		testedDir := fmt.Sprintf("%s/%s", os.TempDir(), "nrd-1")
		testedContainer := RollbackContainerBuilder.NewDir(testedDir)

		err := RollbackExecutor.Run(func() {
			/**
			 * Asserts the creation of temp directory
			 */
			Expect(testedDir).To(BeADirectory())
			// :~)

			/**
			 * Asserts the mode(can read/write/execute) of directory
			 */
			info, _ := os.Stat(testedDir)
			Expect(info.Mode() & 0700).To(BeEquivalentTo(0700))
			// :~)
		}, testedContainer)
		Expect(err).To(Succeed())

		/**
		 * Asserts the removal of temp directory
		 */
		Expect(testedDir).ToNot(BeADirectory())
		// :~)
	})
	It("NewTmpDir", func() {
		testedContainer := RollbackContainerBuilder.NewTmpDir("rollback-executor-*")

		var tempDir string
		err := RollbackExecutor.RunP(func(params Params) {
			tempDir = params[PKEY_TEMP_DIR].(string)

			/**
			 * Asserts the creation of temp directory
			 */
			Expect(tempDir).To(BeADirectory())
			// :~)
		}, testedContainer)
		Expect(err).To(Succeed())

		/**
		 * Asserts the removal of temp directory
		 */
		Expect(tempDir).ToNot(BeADirectory())
		// :~)
	})

	Context("NewChdir", func() {
		var tmpWorkingDir string

		BeforeEach(func() {
			tmpWorkingDir, _ = ioutil.TempDir(os.TempDir(), "test-utils-*")
		})
		AfterEach(func() {
			os.RemoveAll(tmpWorkingDir)
			tmpWorkingDir = ""
		})

		It("Working directory should not be changed", func() {
			testedContainer := RollbackContainerBuilder.NewChdir(tmpWorkingDir)

			oldWd, _ := os.Getwd()

			err := RollbackExecutor.Run(func() {
				/**
				 * Asserts the changed working directory
				 */
				testedWd, _ := os.Getwd()
				Expect(testedWd).To(BeEquivalentTo(tmpWorkingDir))
				// :~)
			}, testedContainer)
			Expect(err).To(Succeed())

			/**
			 * Asserts the working directory of rollbacked
			 */
			testedWd, _ := os.Getwd()
			Expect(testedWd).To(BeEquivalentTo(oldWd))
			// :~)
		})
	})
	Context("NewEnv", func() {
		const sampleEnvVar = "GMT_VAR_1"
		const oldValue = "old_value"

		BeforeEach(func() {
			os.Setenv(sampleEnvVar, "old_value")
		})
		AfterEach(func() {
			os.Unsetenv(sampleEnvVar)
		})

		It("Env-variables are unchanged for Run()", func() {
			const (
				newValue = "new_value"
				newVar = "GMT_VAR_78"
				newVarValue = "new_value_2"
			)

			sampleFunc := func() {
				/**
				 * Asserts the changed env-variables
				 */
				Expect(os.Getenv(sampleEnvVar)).
					To(BeEquivalentTo(newValue))
				Expect(os.Getenv(newVar)).
					To(BeEquivalentTo(newVarValue))
				// :~)
			}

			testedContainer := RollbackContainerBuilder.NewEnv(map[string]string {
				sampleEnvVar: newValue,
				newVar: newVarValue,
			})
			RollbackExecutor.Run(sampleFunc, testedContainer)

			/**
			 * Asserts the rollback of env-variables and introduced(removed) env-variable
			 */
			Expect(os.Getenv(sampleEnvVar)).
				To(BeEquivalentTo(oldValue))

			_, newOne := os.LookupEnv(newVar)
			Expect(newOne).To(BeFalse())
			// :~)
		})
	})
	Context("NewCopyFiles", func() {
		var (
			params Params
			tempDir string
			tempFile1, tempFile2 string
		)
		tempDirContainer := RollbackContainerBuilder.NewTmpDir("rcf-*")
		file1 := fmt.Sprintf("%s/%s", testSourceDir, "sample-1.txt")
		file2 := fmt.Sprintf("%s/%s", testSourceDir, "sample-2.txt")

		/**
		 * Setup temporary directory
		 */
		BeforeEach(func() {
			params, _ = tempDirContainer.Setup()
			tempDir = params[PKEY_TEMP_DIR].(string)
			tempFile1 = fmt.Sprintf("%s/%s", tempDir, "sample-1.txt")
			tempFile2 = fmt.Sprintf("%s/%s", tempDir, "sample-2.txt")
		})
		AfterEach(func() {
			tempDirContainer.TearDown(params)
		})
		// :~)

		It("Ensures the copying/removal of files", func() {
			testedContainer := RollbackContainerBuilder.NewCopyFiles(
				tempDir, file1, file2,
			)
			err := RollbackExecutor.Run(func() {
				Expect(tempFile1).To(BeAnExistingFile())
				Expect(tempFile2).To(BeAnExistingFile())
			}, testedContainer)

			/**
			 * Asserts the removal of files
			 */
			Expect(err).To(Succeed())
			Expect(tempFile1).ToNot(BeAnExistingFile())
			Expect(tempFile2).ToNot(BeAnExistingFile())
			// :~)
		})
	})

	Context("rollbackExecutorPImpl", func() {
		It("Setup has error for containers", func() {
			err := RollbackExecutor.RunP(
				func(Params) {},
				&errorContainer{ setupError: false, tearDownError: false },
				&errorContainer{ setupError: true, tearDownError: false },
			)

			Expect(err).To(MatchError("Setup-Error"))
		})

		It("Teardown has error for containers", func() {
			err := RollbackExecutor.RunP(
				func(Params) {},
				&errorContainer{ setupError: false, tearDownError: false },
				&errorContainer{ setupError: false, tearDownError: true },
			)

			Expect(err).To(MatchError("TearDown-Error"))
		})

		Describe("Sequence of multiple containers", func() {
			/**
			 * Runs the result
			 */
			seqContainer := &numberedContainer{ 0, make([]int, 0), make([]int, 0) }
			err := RollbackExecutor.RunP(
				func(Params) {},
				seqContainer, seqContainer, seqContainer,
			)
			// :~)

			It("The sequence of setup should be as arguments", func() {
				Expect(err).To(Succeed())
				Expect(seqContainer.setupSeq[0]).
					To(BeEquivalentTo(1))
				Expect(seqContainer.setupSeq[1]).
					To(BeEquivalentTo(2))
				Expect(seqContainer.setupSeq[2]).
					To(BeEquivalentTo(3))
			})
			It("The sequence of teardown should be as arguments(reverse)", func() {
				Expect(err).To(Succeed())
				Expect(seqContainer.tearDownSeq[0]).
					To(BeEquivalentTo(3))
				Expect(seqContainer.tearDownSeq[1]).
					To(BeEquivalentTo(2))
				Expect(seqContainer.tearDownSeq[2]).
					To(BeEquivalentTo(1))
			})
		})

		It("Ensure the calling of methods and callback", func() {
			touchCallback := false
			testedContainer := &sampleContainerP{ false, false }

			/**
			 * Prepares the executor
			 */
			RollbackExecutor.RunP(func(Params) {
				/**
				 * Only if the setup gets called, tear down is not get called
				 */
				if testedContainer.setup && !testedContainer.tearDown {
					touchCallback = true
				}
				// :~)
			}, testedContainer)
			// :~)

			/**
			 * Asserts the calling of every method
			 */
			Expect(touchCallback).To(BeTrue())
			Expect(testedContainer.setup).To(BeTrue())
			Expect(testedContainer.tearDown).To(BeTrue())
			// :~)
		})
	})

	Context("Concate", func() {
		toC := func(containerP RollbackContainerP) RollbackContainer {
			return RollbackContainerBuilder.ToContainer(containerP)
		}

		assertCalled := func(containerP *sampleContainerP, expectedSetup bool, expectedTearDown bool) {
			ExpectWithOffset(1, containerP.setup).To(BeEquivalentTo(expectedSetup), "Setup should be: %v", expectedSetup)
			ExpectWithOffset(1, containerP.tearDown).To(BeEquivalentTo(expectedTearDown), "Tear down should be: %v", expectedTearDown)
		}

		When("Every container is successfully set-up/tear-down", func() {
			It("Setup/TearDown gets called for every container", func() {
				sampleContainers := []*sampleContainerP {
					new(sampleContainerP), new(sampleContainerP), new(sampleContainerP),
				}

				testedContainer := RollbackContainerBuilder.Concate(
					toC(sampleContainers[0]), toC(sampleContainers[1]), toC(sampleContainers[2]),
				)

				testedContainer.Setup()
				testedContainer.TearDown()

				assertCalled(sampleContainers[0], true, true)
				assertCalled(sampleContainers[1], true, true)
				assertCalled(sampleContainers[2], true, true)
			})
		})
		When("Some containers are not set-up/tear-down", func() {
			It("Some Setup/TearDown don't get called", func() {
				errContainer := &errorContainer{
					sampleContainerP: new(sampleContainerP),
					setupError: true,
				}
				sampleContainers := []*sampleContainerP {
					new(sampleContainerP), errContainer.sampleContainerP,
					new(sampleContainerP),
				}

				testedContainer := RollbackContainerBuilder.Concate(
					toC(sampleContainers[0]), toC(errContainer), toC(sampleContainers[2]),
				)

				testedContainer.Setup()
				testedContainer.TearDown()

				assertCalled(sampleContainers[0], true, true)
				assertCalled(sampleContainers[1], true, false)
				assertCalled(sampleContainers[2], false, false)
			})
		})
	})
}

type errorContainer struct {
	*sampleContainerP
	params Params

	setupError bool
	tearDownError bool
}
func(self *errorContainer) Setup() (Params, error) {
	if (self.sampleContainerP != nil) {
		self.sampleContainerP.Setup()
	}

	var err error
	if self.setupError {
		err = fmt.Errorf("Setup-Error")
	}

	return nil, err
}
func(self *errorContainer) TearDown(Params) error {
	if self.sampleContainerP != nil {
		self.sampleContainerP.TearDown(nil)
	}

	var err error
	if self.tearDownError {
		err = fmt.Errorf("TearDown-Error")
	}

	return err
}

type numberedContainer struct {
	id int

	setupSeq []int
	tearDownSeq []int
}
func(self *numberedContainer) Setup() (Params, error) {
	self.id++
	self.setupSeq = append(self.setupSeq, self.id)
	return nil, nil
}
func(self *numberedContainer) TearDown(Params) error {
	self.tearDownSeq = append(self.tearDownSeq, self.id)
	self.id--
	return nil
}

type sampleContainerP struct {
	setup bool
	tearDown bool
}
func(self *sampleContainerP) Setup() (Params, error) {
	self.setup = true
	return nil, nil
}
func(self *sampleContainerP) TearDown(Params) error {
	/**
	 * Only if the setup gets called
	 */
	if self.setup {
		self.tearDown = true
	}
	// :~)
	return nil
}
