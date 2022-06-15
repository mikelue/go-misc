package env

import (
	"os"

	"github.com/spf13/pflag"

	fg "github.com/mikelue/go-misc/ioc/frangipani"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("arguments loader", func() {
	Context("argsConfig.bindByPflag", contextOfBindByPflag)

	It("setPrefix", func() {
		testedArgsConfig := new(argsConfig).
			setPrefix("ok-01")

		Expect(testedArgsConfig.prefix).
			To(BeEquivalentTo("ok_01"))
		Expect(testedArgsConfig.ordinaryPrefix).
			To(BeEquivalentTo("ok-01"))
	})
})

var contextOfBindByPflag = func() {
	var testedPackedConfig *packedConfig

	BeforeEach(func() {
		flagSet := pflag.NewFlagSet("test-args-loader", pflag.ExitOnError)

		argsConfig := &argsConfig{}
		argsConfig.setPrefix("myapp")
		testedPackedConfig = argsConfig.bindByPflag(flagSet)

		flagSet.Parse(
			[]string {
				`--myapp.config.json={ "apple.size": 98 }`,
				`--myapp.config.yaml={ "tamarind.size": 62 }`,
				`--myapp.config.files=a1.yaml,a2.yaml`,
				`--myapp.profiles.active=g1,g2`,
			},
		)
	})
	AfterEach(func() {
		pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	})

	Context("bindByPflag", func() {
		Context("JSON/YAML as properties", func() {
			var testedEnv fg.Environment

			BeforeEach(func() {
				testedEnv = fg.EnvBuilder.NewByVipers(testedPackedConfig.loadFormattedProps()...)
			})

			It("yaml", func() {
				Expect(testedEnv.Typed().GetInt("tamarind.size")).
					To(BeEquivalentTo(62))
			})

			It("json", func() {
				Expect(testedEnv.Typed().GetInt("apple.size")).
					To(BeEquivalentTo(98))
			})
		})

		It("file name", func() {
			testedViper := testedPackedConfig.configFileProp()

			Expect(testedViper.GetStringSlice("myapp.config.files")).
				To(ConsistOf("a1.yaml", "a2.yaml"))
		})

		It("active profiles", func() {
			testedViper := testedPackedConfig.loadActiveProfiles()

			Expect(testedViper.GetString("myapp.profiles.active")).
				To(BeEquivalentTo("g1,g2"))
		})
	})
}
