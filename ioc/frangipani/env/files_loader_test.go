package env

import (
	"fmt"

	fg "github.com/mikelue/go-misc/ioc/frangipani"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader for file names", func() {
	Context("configFileNames", contextOfConfigFileNames)

	DescribeTable("readInByFile",
		func(filename string, expected bool) {
			_, testedResult := readInByFile(fmt.Sprintf("%s/%s", currentSrcDir, filename), true)

			Expect(testedResult).To(BeEquivalentTo(expected))
		},
		Entry("Success", "sample-1.yaml", true),
		Entry("Failed because it is not existing", "sample76.yaml", false),
		Entry("Failed because it is not correct", "sample-2-err.yaml", false),
	)
})

var contextOfConfigFileNames = func() {
	It("Load()", func() {
		files := []string {
			fmt.Sprintf("%s/%s", currentSrcDir, "sample-1.yaml"),
		}
		fileNames := &configFileNames{ files, false }

		testedEnv := fg.EnvBuilder.NewByVipers(fileNames.load()...)
		Expect(testedEnv.GetProperty("gdb.key")).
			To(BeEquivalentTo("G6Coi2S4"))
	})

	It("LoadByDir()", func() {
		files := []string { "sample-1.yaml" }
		fileNames := &configFileNames{ files, false }

		testedEnv := fg.EnvBuilder.NewByVipers(fileNames.loadByDir(currentSrcDir)...)
		Expect(testedEnv.GetProperty("gdb.key")).
			To(BeEquivalentTo("G6Coi2S4"))
	})

	Context("Priority of multiple files", func() {
		DescribeTable("Multiple different files",
			func(files []string, expected string) {
				fileNames := &configFileNames{ files, false }

				vipers := fileNames.loadByDir(currentSrcDir)
				testedEnv := fg.EnvBuilder.NewByVipers(vipers...)

				Expect(testedEnv.GetProperty("db.sample.key")).
					To(BeEquivalentTo(expected))
			},
			Entry("json file has priority", []string { "lime-config-p1.yaml", "lime-config-p2.yaml" }, "YZZnOcpw"),
			Entry("yaml file has priority", []string { "lime-config-p2.yaml", "lime-config-p1.yaml" }, "gAyqlZyq"),
		)

		DescribeTable("Multiple profiles",
			func(profiles []string, expected string) {
				fileNames := &configFileNames{ []string { "split-peas-config.yaml" }, false }

				vipers := fileNames.loadByDirWithProfiles(currentSrcDir, profiles...)
				testedEnv := fg.EnvBuilder.NewByVipers(vipers...)

				Expect(testedEnv.GetProperty("db.sample.host")).
					To(BeEquivalentTo(expected))
			},
			Entry("p1 has priority over p2", []string{ "p1", "p2" }, "192.186.6.50"),
			Entry("p2 has priority over p1", []string{ "p2", "p1" }, "192.186.6.70"),
		)
	})

	Context("Load with supported types of file", func() {
		DescribeTable("Every supported type",
			func(
				sampleFile string,
				name string, expected string,
			) {
				fileNames := &configFileNames{
					[]string{
						fmt.Sprintf("%s/%s", currentSrcDir, sampleFile),
					},
					true,
				}
				vipers := fileNames.load()

				Expect(vipers).To(HaveLen(1))
				Expect(vipers[0].GetString(name)).To(BeEquivalentTo(expected))
			},
			Entry(".properties file", "sample-1.properties", "gdb.key", "3IZg3eQZ"),
			Entry(".json file", "sample-1.json", "gdb.key", "RxTEHH4s"),
			Entry(".yaml file", "sample-1.yaml", "gdb.key", "G6Coi2S4"),
		)
	})
}
