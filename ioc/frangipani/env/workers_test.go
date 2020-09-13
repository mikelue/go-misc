package env

import (
	"github.com/mikelue/go-misc/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workers", func() {
	Context("workerBuilderI", func() {
		Context("newXdg", func() {
			When("Found a XDG dir", func() {
				const prefix = "test-c1"
				var tmpDir string
				var params utils.Params
				tmpDirContainer := utils.RollbackContainerBuilder.
					NewTmpDir("fake-xdg-*")
				var xdgContainers []utils.RollbackContainer

				BeforeEach(func() {
					params, _ := tmpDirContainer.Setup()
					tmpDir = params[utils.PKEY_TEMP_DIR].(string)
					GinkgoT().Logf("Temporary \"$XDG_CONFIG_HOME\": %s", tmpDir)

					xdgContainers = newXdgSetup(tmpDir, prefix)
				})
				AfterEach(func() {
					tmpDirContainer.TearDown(params)
				})

				It("The existing of XDG directory", func() {
					utils.RollbackExecutor.Run(
						func() {
							testedWorker := workerBuilder.newXdg(prefix).(*filesWorker)
							Expect(testedWorker.targetDir).ToNot(BeEmpty())
						},
						xdgContainers...,
					)
				})
			})
			It("Not found a XDG dir", func() {
				testedWorker := workerBuilder.newXdg("not-existing").(*filesWorker)

				Expect(testedWorker.targetDir).To(BeEmpty())
				Expect(testedWorker.load()).To(BeEmpty())
				Expect(testedWorker.loadWithProfiles()).To(BeEmpty())
			})
		})

		Context("newWd", func() {
			var tmpDir string
			var params utils.Params
			tmpDirContainer := utils.RollbackContainerBuilder.
				NewTmpDir("fake-wd-*")

			BeforeEach(func() {
				params, _ := tmpDirContainer.Setup()
				tmpDir = params[utils.PKEY_TEMP_DIR].(string)
				GinkgoT().Logf("Temporary \"working directory\": %s", tmpDir)
			})
			AfterEach(func() {
				tmpDirContainer.TearDown(params)
			})

			It("Under working folder", func() {
				chdirContainer := utils.RollbackContainerBuilder.NewChdir(tmpDir)

				utils.RollbackExecutor.Run(
					func() {
						testedWorker := workerBuilder.newWd("wd-test").(*filesWorker)
						Expect(testedWorker.targetDir).To(BeEquivalentTo(tmpDir))
					},
					chdirContainer,
				)
			})
		})

		Context("newCmdDir", func() {
			It("By cmd dir", func() {
				testedWorker := workerBuilder.newCmdDir("cmd-test").(*filesWorker)

				GinkgoT().Logf("Cmd dir: %s", testedWorker.targetDir)
				Expect(testedWorker.targetDir).ToNot(BeEmpty())
			})
		})
	})

	Context("filesWorker", func() {
		testedWorker := newDefaultFilesWorker("split-peas")
		testedWorker.targetDir = currentSrcDir
		testedWorker.files = []string{ "split-peas-config.yaml", "lime-config.yaml" }

		It("load", func() {
			testedVipers := testedWorker.load()
			Expect(testedVipers).To(HaveLen(2))

			Expect(testedVipers[0].GetString("db.sample.host")).
				To(BeEquivalentTo("192.186.21.50"))
			Expect(testedVipers[1].GetString("db.sample.key")).
				To(BeEquivalentTo("LdTB83tK"))
		})

		It("loadWithProfiles", func() {
			testedVipers := testedWorker.loadWithProfiles("p1", "p2")
			Expect(testedVipers).To(HaveLen(6))

			Expect(testedVipers[0].GetString("db.sample.host")).
				To(BeEquivalentTo("192.186.6.50"))
			Expect(testedVipers[2].GetString("db.sample.key")).
				To(BeEquivalentTo("YZZnOcpw"))
		})
	})
})
