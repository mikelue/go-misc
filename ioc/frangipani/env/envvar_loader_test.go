package env

import (
	"github.com/mikelue/go-misc/utils"

	fg "github.com/mikelue/go-misc/ioc/frangipani"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loader by environment variables", func() {
	Context("envConfig.load()", contextOfLoadByEnv)

	It("Set prefix", func() {
		testedConfig := new(envConfig).
			setPrefix("jj-o90")

		Expect(testedConfig.prefix).
			To(BeEquivalentTo("JJ_O90"))
		Expect(testedConfig.ordinaryPrefix).
			To(BeEquivalentTo("jj-o90"))
	})
})

var contextOfLoadByEnv = func() {
	var testedPackedConfig *packedConfig

	envContainer := utils.RollbackContainerBuilder.NewEnv(
		map[string]string {
			"KKAPP_CONFIG_YAML": `{ srv.lion.port: 981 }`,
			"KKAPP_CONFIG_JSON": `{ "srv.deer.port": 651 }`,
			"KKAPP_CONFIG_FILES": "o-sample-9.json,o-sample-10.json",
			"KKAPP_PROFILES_ACTIVE": "h1,h2",
		},
	)

	BeforeEach(func() {
		envContainer.Setup()

		envConfig := new(envConfig)
		envConfig.setPrefix("kkapp")
		testedPackedConfig = envConfig.load()
	})
	AfterEach(func() {
		envContainer.TearDown()
	})

	Context("JSON/YAML as properties", func() {
		var testedEnv fg.Environment

		BeforeEach(func() {
			testedEnv = fg.EnvBuilder.NewByVipers(testedPackedConfig.loadFormattedProps()...)
		})

		It("yaml", func() {
			Expect(testedEnv.Typed().GetInt("srv.lion.port")).
				To(BeEquivalentTo(981))
		})

		It("json", func() {
			Expect(testedEnv.Typed().GetInt("srv.deer.port")).
				To(BeEquivalentTo(651))
		})
	})

	It("file name", func() {
		testedViper := testedPackedConfig.configFileProp()

		Expect(testedViper.GetStringSlice("kkapp.config.files")).
			To(ConsistOf("o-sample-9.json", "o-sample-10.json"))
	})

	It("active profiles", func() {
		testedViper := testedPackedConfig.loadActiveProfiles()

		Expect(testedViper.GetString("kkapp.profiles.active")).
			To(BeEquivalentTo("h1,h2"))
	})
}
