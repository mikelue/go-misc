package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mrt "github.com/mikelue/go-misc/utils/runtime"
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
		currentDir := mrt.CallerUtils.GetDirOfSource()
		file1 := fmt.Sprintf("%s/%s", currentDir, "sample-1.txt")
		file2 := fmt.Sprintf("%s/%s", currentDir, "sample-2.txt")

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
			testedContainer := &sampleContainer{ false, false }

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

type sampleContainer struct {
	setup bool
	tearDown bool
}
func(self *sampleContainer) Setup() (Params, error) {
	self.setup = true
	return nil, nil
}
func(self *sampleContainer) TearDown(Params) error {
	/**
	 * Only if the setup gets called
	 */
	if self.setup {
		self.tearDown = true
	}
	// :~)
	return nil
}
