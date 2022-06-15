package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	fg "github.com/mikelue/go-misc/ioc/frangipani"
	"github.com/mikelue/go-misc/utils"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Load configuration and active profiles(properties)", func() {
	Context("ConfigBuilder", contextOfConfigBuilder)

	Context("1st Pass", contextOf1stPass)
	Context("Load", contextOfLoad)

	Context("eliminateSameDirForWdAndCmd", contextOfeliminateSameDirOfWdAndCmd)
})

var contextOfeliminateSameDirOfWdAndCmd = func() {
	toAnySlice := func(sources []ConfigSource) []interface{} {
		return funk.Map(sources,
			func(v interface{}) interface{} { return v },
		).
			([]interface{})
	}

	When("Nothing changed", func() {
		DescribeTable("The sources as same as input",
			func(sampleSources []ConfigSource) {
				testedSources := eliminateSameDirForWdAndCmd(sampleSources)

				Expect(testedSources).
					To(ConsistOf(toAnySlice(sampleSources)...))
			},
			Entry("Only working directory", []ConfigSource{ CL_PWD, CL_ARGS, CL_ENVVAR }),
			Entry("Only directory of command", []ConfigSource{ CL_ARGS, CL_ENVVAR, CL_CMDDIR }),
			Entry("Neither of both are used", []ConfigSource{ CL_XDG, CL_CONFIG_FILE }),
		)
	})

	When("The current directory and the directory of command is not same", func() {
		var tmpDir string
		var params utils.Params
		var tmpDirContainer utils.RollbackContainerP

		BeforeEach(func() {
			tmpDirContainer = utils.RollbackContainerBuilder.NewTmpDir("temp-wd-*")
			params, _ = tmpDirContainer.Setup()
			tmpDir = params[utils.PKEY_TEMP_DIR].(string)
		})
		AfterEach(func() {
			tmpDirContainer.TearDown(params)
		})

		DescribeTable("Nothing changed",
			func(sampleSources []ConfigSource) {
				chDirToTempDir := utils.RollbackContainerBuilder.NewChdir(tmpDir)

				err := utils.RollbackExecutor.Run(
					func() {
						testedSources := eliminateSameDirForWdAndCmd(sampleSources)

						Expect(testedSources).
							To(ConsistOf(toAnySlice(sampleSources)...))
					},
					chDirToTempDir,
				)
				Expect(err).To(Succeed())
			},
			Entry("CL_PWD has higher priority", []ConfigSource{ CL_PWD, CL_ARGS, CL_CMDDIR }),
			Entry("CL_CMDDIR has higher priority", []ConfigSource{ CL_CMDDIR, CL_ARGS, CL_PWD }),
		)
	})

	When("The current directory and the directory of command is same", func() {
		cmdDir := getCmdDir()
		chDirToCmdDir := utils.RollbackContainerBuilder.NewChdir(cmdDir)

		DescribeTable("Remove least priority source",
			func(sampleSources []ConfigSource) {
				err := utils.RollbackExecutor.Run(
					func() {
						testedSources := eliminateSameDirForWdAndCmd(sampleSources)

						Expect(testedSources).
							To(ConsistOf(toAnySlice(sampleSources[:2])...))
					},
					chDirToCmdDir,
				)
				Expect(err).To(Succeed())
			},
			Entry("Remove CL_CMDDIR", []ConfigSource{ CL_PWD, CL_ARGS, CL_CMDDIR }),
			Entry("Remove CL_PWD", []ConfigSource{ CL_CMDDIR, CL_ARGS, CL_PWD }),
		)
	})
}

var contextOfConfigBuilder = func() {
	When("Duplicated sources", func() {
		It("Should show the warning log", func() {
			NewConfigBuilder().Priority(
				CL_ARGS, CL_ARGS,
			)
		})
	})
}

var contextOf1stPass = func() {
	var (
		tempDir string
		params utils.Params
		builder *configLoaderImpl
		counter = 0
	)

	tmpDirContainer := utils.RollbackContainerBuilder.
		NewTmpDir("1pass-xdg-*")

	BeforeEach(func() {
		params, _ = tmpDirContainer.Setup()
		tempDir = params[utils.PKEY_TEMP_DIR].(string)
	})
	AfterEach(func() {
		tmpDirContainer.TearDown(params)
	})

	newFlagSet := func() *pflag.FlagSet {
		counter++
		flagSetName := fmt.Sprintf("test-pass-1-%d", counter)
		return pflag.NewFlagSet(flagSetName, pflag.ExitOnError)
	}

	loader := func() *configLoaderImpl {
		builder = NewConfigBuilder().
			Prefix("lime").
			Pflags(newFlagSet()).
			Build().(*configLoaderImpl)

		return builder
	}
	runFirstPass := func() fg.Environment {
		return loader().pass1Load()
	}

	// Since "ginkgo run" is differ from "go test" for the
	// current working directory while running tests,
	// the former way to run tests would remove "lime-config...yaml" files.
	copyFiles := func(dir string) utils.RollbackContainer {
		_, err := os.Stat(fmt.Sprintf("%s/lime-config.yaml", dir))

		if errors.Is(err, os.ErrNotExist) {
			return utils.RollbackContainerBuilder.NewCopyFiles(
				dir,
				fmt.Sprintf("%s/lime-config.yaml", currentSrcDir),
				fmt.Sprintf("%s/lime-config-p1.yaml", currentSrcDir),
				fmt.Sprintf("%s/lime-config-p2.yaml", currentSrcDir),
			)
		}

		return utils.EmptyRollbackContainer
	}

	assertEnvByFile := func(env fg.Environment) {
		testedProps := env.Typed()

		Expect(testedProps.GetString("db.sample.key")).
			To(BeEquivalentTo("LdTB83tK"))
	}

	It("Default values", func() {
		testedEnv := NewConfigBuilder().
			Prefix("lime").
			DefaultWithMap(map[string]interface{} {
				"kc.weight": 279,
			}).
			Pflags(newFlagSet()).
			Build().(*configLoaderImpl).
			pass1Load().
			Typed()

		Expect(testedEnv.GetInt("kc.weight")).To(BeEquivalentTo(279))
	})

	It("Loaded by $XDG_CONFIG_HOME/lime-config.yaml", func() {
		finalDir := fmt.Sprintf("%s/%s", tempDir, "lime")
		setupContainers := newXdgSetup(tempDir, "lime")
		setupContainers = append(setupContainers, copyFiles(finalDir))

		err := utils.RollbackExecutor.Run(
			func() {
				/**
				 * Asserts the final environment
				 */
				assertEnvByFile(runFirstPass())
				// :~)
			},
			setupContainers...,
		)
		Expect(err).To(Succeed())
	})

	It("Load by environment", func() {
		envContainer := utils.RollbackContainerBuilder.NewEnv(
			map[string]string {
				"LIME_CONFIG_YAML": `{ cv_1: 91 }`,
				"LIME_CONFIG_JSON": `{ "cv_1": 92, "cv_2": 104 }`,
				"LIME_CONFIG_FILES": "g1.yaml,g2.yaml",
				"LIME_PROFILES_ACTIVE": "ky1,ky2",
			},
		)

		err := utils.RollbackExecutor.Run(
			func() {
				/**
				 * Asserts the result environment
				 */
				testedProps := runFirstPass().Typed()
				Expect(testedProps.GetInt("cv_1")).To(BeEquivalentTo(91))
				Expect(testedProps.GetInt("cv_2")).To(BeEquivalentTo(104))
				Expect(testedProps.GetStringSlice("lime.config.files")).
					To(ConsistOf("g1.yaml", "g2.yaml"))
				Expect(testedProps.GetString("lime.profiles.active")).
					To(BeEquivalentTo("ky1,ky2"))
				Expect(testedProps.GetString("fgapp.profiles.active")).
					To(BeEquivalentTo("ky1,ky2"))
				// :~)
			},
			envContainer,
		)
		Expect(err).To(Succeed())
	})

	It("Load by current directory", func() {
		err := utils.RollbackExecutor.Run(
			func() {
				/**
				 * Asserts the final environment
				 */
				assertEnvByFile(runFirstPass())
				// :~)
			},
			copyFiles(tempDir),
			utils.RollbackContainerBuilder.NewChdir(tempDir),
		)
		Expect(err).To(Succeed())
	})

	It("Load by directory of executable", func() {
		executableDir, _ := filepath.Abs(os.Args[0])
		executableDir = filepath.Dir(executableDir)

		err := utils.RollbackExecutor.Run(
			func() {
				/**
				 * Asserts the final environment
				 */
				assertEnvByFile(runFirstPass())
				// :~)
			},
			copyFiles(executableDir),
			utils.RollbackContainerBuilder.NewChdir(tempDir),
		)
		Expect(err).To(Succeed())
	})

	Context("Parse arguments", func() {
		newLoader := loader()
		newLoader.flags.Parse([]string{
			`--lime.config.yaml={ a1: 20 }`,
			`--lime.config.json={ "a1": 21, "a2": 30 }`,
			fmt.Sprintf("--lime.config.files=%s", "z1.yaml,z2.yaml"),
			`--lime.profiles.active=u1,u2`,
		})

		It("Parsed configurations by arguments", func() {
			testedEnv := newLoader.pass1Load().Typed()

			/**
			 * Asserts the final environment
			 */
			Expect(testedEnv.GetInt("a1")).To(BeEquivalentTo(20))
			Expect(testedEnv.GetInt("a2")).To(BeEquivalentTo(30))
			Expect(testedEnv.GetStringSlice("lime.config.files")).
				To(ConsistOf("z1.yaml", "z2.yaml"))
			Expect(testedEnv.GetString("lime.profiles.active")).
				To(BeEquivalentTo("u1,u2"))
			Expect(testedEnv.GetString("fgapp.profiles.active")).
				To(BeEquivalentTo("u1,u2"))
			// :~)
		})
	})
}

var contextOfLoad = func() {
	var oldOsArgs []string

	BeforeEach(func() {
		oldOsArgs = os.Args
	})
	AfterEach(func() {
		os.Args = oldOsArgs
	})

	Context("DefaultLoader", func() {
		BeforeEach(func() {
			os.Args = []string {
				`--fgapp.config.yaml={ aksrv.host: 60.81.70.145 }`,
			}
			pflag.CommandLine = pflag.NewFlagSet("test-DefaultLoader", pflag.ExitOnError)
		})

		It("New()", func() {
			testedEnv := DefaultLoader.New().
				ParseFlags().
				Load()

			Expect(testedEnv.GetProperty("aksrv.host")).
				To(BeEquivalentTo("60.81.70.145"))
		})

		It("WithMap()", func() {
			testedEnv := DefaultLoader.WithMap(
				map[string]interface{} {
					"gzsrv.host": "98.71.104.176",
				},
			).
				ParseFlags().
				Load()

			Expect(testedEnv.GetProperty("aksrv.host")).
				To(BeEquivalentTo("60.81.70.145"))
			Expect(testedEnv.GetProperty("gzsrv.host")).
				To(BeEquivalentTo("98.71.104.176"))
		})

		It("WithViper()", func() {
			viper := viper.New()
			viper.Set("ucsrv.host", "71.44.71.81")
			testedEnv := DefaultLoader.WithViper(viper).
				ParseFlags().
				Load()

			Expect(testedEnv.GetProperty("aksrv.host")).
				To(BeEquivalentTo("60.81.70.145"))
			Expect(testedEnv.GetProperty("ucsrv.host")).
				To(BeEquivalentTo("71.44.71.81"))
		})
	})

	Context("Default priority of loading", func() {
		const prefix = "lime"

		tmpDirContainer := utils.RollbackContainerBuilder.
			NewTmpDir("fake-xdg-dp-*")
		tmpWorkingDir := utils.RollbackContainerBuilder.
			NewTmpDir("fake-wd-dp-*")
		var prepareAllEnv utils.RollbackContainer
		var params utils.Params
		var wdParams utils.Params
		var testedEnv fg.Environment

		BeforeEach(func() {
			var err error

			/**
			 * Creates direcotries:
			 * 1. For $XDG_CONFIG_HOME
			 * 2. For working direcotry
			 */
			params, _ := tmpDirContainer.Setup()
			tmpDir := params[utils.PKEY_TEMP_DIR].(string)
			GinkgoT().Logf("Temporary \"$XDG_CONFIG_HOME\": %s", tmpDir)
			wdParams, _ = tmpWorkingDir.Setup()

			viper := viper.New()
			viper.Set("mail.admin", "test-1@fgapp.org")
			newWd := wdParams[utils.PKEY_TEMP_DIR].(string)
			filename := fmt.Sprintf("%s/lime-config.yaml", newWd)
			err = viper.WriteConfigAs(filename)
			if err != nil {
				GinkgoT().Errorf("Unable to write YAML file[%s]: %v", filename, err)
			}
			// :~)

			/**
			 * 1. Copy sample files in $XDG_CONFIG_HOME
			 * 2. Prepares environment variables
			 * 3. Changes working folder
			 */
			allPrepareContainers := newXdgSetup(tmpDir, prefix)
			allPrepareContainers = append(allPrepareContainers,
				utils.RollbackContainerBuilder.NewCopyFiles(
					fmt.Sprintf("%s/%s", tmpDir, prefix),
					fmt.Sprintf("%s/lime-config.yaml", currentSrcDir),
					fmt.Sprintf("%s/lime-config-p1.yaml", currentSrcDir),
				),
				utils.RollbackContainerBuilder.NewEnv(
					map[string]string {
						"LIME_CONFIG_YAML": `{
							db.sample.cluster_name: asia-east-1
						}`,
						"LIME_CONFIG_FILES": fmt.Sprintf("%s/split-peas-config.yaml", currentSrcDir),
					},
				),
				utils.RollbackContainerBuilder.NewChdir(newWd),
			)
			prepareAllEnv = utils.RollbackContainerBuilder.Concate(allPrepareContainers...)
			prepareAllEnv.Setup()
			if err != nil {
				GinkgoT().Errorf("Setup all of the enviornments has error: %v", err)
			}
			// :~)

			/**
			 * Prepares arguments
			 */
			os.Args = []string {
				`--lime.config.yaml={ db.sample.host: 103.75.223.120 }`,
				`--lime.profiles.active=p1`,
			}
			pflag.CommandLine = pflag.NewFlagSet("test-default-priority", pflag.ExitOnError)
			// :~)

			testedEnv = NewConfigBuilder().Prefix(prefix).
				Build().ParseFlags().
				Load()
		})
		AfterEach(func() {
			prepareAllEnv.TearDown()

			// Removes working direcotry(temporary created)
			tmpWorkingDir.TearDown(wdParams)

			// This would remove all of the files recursively
			tmpDirContainer.TearDown(params)
		})

		It("XDG properties(Includes profile)", func() {
			Expect(testedEnv.GetProperty("db.sample.key")).
				To(BeEquivalentTo("YZZnOcpw"))
		})

		It("Arguments", func() {
			Expect(testedEnv.GetProperty("db.sample.host")).
				To(BeEquivalentTo("103.75.223.120"))
		})

		It("Environment variables", func() {
			Expect(testedEnv.GetProperty("db.sample.cluster_name")).
				To(BeEquivalentTo("asia-east-1"))
		})

		It("Config files", func() {
			Expect(testedEnv.Typed().GetInt("cassandra.port")).
				To(BeEquivalentTo(9891))
		})

		It("Load configuration file by working folder", func() {
			Expect(testedEnv.GetProperty("mail.admin")).
				To(BeEquivalentTo("test-1@fgapp.org"))
		})
	})

	Context("Customized priority", func() {
		BeforeEach(func() {
			os.Args = []string {
				`--lime.config.yaml={ db.sample.host: 160.181.27.90 }`,
				fmt.Sprintf(`--lime.config.files=%s/split-peas-config.yaml`, currentSrcDir),
			}
			pflag.CommandLine = pflag.NewFlagSet("test-CustomizedPriority", pflag.ExitOnError)
		})

		It("The config file should have priority", func() {
			testedEnv := NewConfigBuilder().
				Prefix("lime").
				Priority(CL_CONFIG_FILE, CL_ARGS).
				Build().
				ParseFlags().
				Load()

			Expect(testedEnv.GetProperty("db.sample.host")).
				To(BeEquivalentTo("192.186.21.50"))
		})
	})
}
