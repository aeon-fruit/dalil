package config_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
)

var _ = Describe("Config", func() {

	const (
		keyAppEnv                     = "APP_ENV"
		keyAppPort                    = "APP_PORT"
		keyAppLoggingVerbosityGlobal  = "APP_LOGGING_VERBOSITY_GLOBAL"
		keyAppLoggingVerbosityModules = "APP_LOGGING_VERBOSITY_MODULES"

		defaultAppEnv = config.AppEnvLocal
		customAppEnv  = config.AppEnvNonProd

		defaultAppPort = 8080
		customAppPort  = 12345

		defaultAppLoggingVerbosityGlobal = 0
		customAppLoggingVerbosityGlobal  = 1

		moduleParent = "parent"
		moduleNode   = "node"
		moduleLeaf   = "leaf"

		customAppLoggingVerbosityModulesEnvVar = moduleParent + "=2," + moduleNode + "=1," + moduleLeaf + "=3"
	)

	var (
		customAppLoggingVerbosityModules = map[string]int{
			moduleParent: 2,
			moduleNode:   1,
			moduleLeaf:   3,
		}
	)

	Describe("New", func() {
		It("has defaults", func() {
			instance := config.New()

			Expect(instance.AppEnv).To(Equal(defaultAppEnv))
			Expect(instance.AppPort).To(Equal(defaultAppPort))
			Expect(instance.Logging.GetGlobalVerbosity()).To(Equal(defaultAppLoggingVerbosityGlobal))
			Expect(instance.Logging.GetModules()).To(BeEmpty())
		})

		Context("WithEnvVars is specified", func() {
			BeforeEach(func() {
				err := os.Unsetenv(keyAppPort)
				Expect(err).NotTo(HaveOccurred())
			})

			When("there is no environment variables", func() {
				It("keeps the defaults", func() {
					instance := config.New(config.WithEnvVars())

					Expect(instance.AppPort).To(Equal(defaultAppPort))
				})
			})

			When("there are environment variables but the values cannot be parsed", func() {
				It("keeps the defaults", func() {
					err := os.Setenv(keyAppEnv, "a random env")
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppPort, "non integer value")
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppLoggingVerbosityGlobal, "debug")
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppLoggingVerbosityModules, "parent:1,  =3,node=  ,child=2;leaf=trace")
					Expect(err).NotTo(HaveOccurred())

					instance := config.New(config.WithEnvVars())

					Expect(instance.AppEnv).To(Equal(defaultAppEnv))
					Expect(instance.AppPort).To(Equal(defaultAppPort))
					Expect(instance.Logging.GetGlobalVerbosity()).To(Equal(defaultAppLoggingVerbosityGlobal))
					Expect(instance.Logging.GetModules()).To(BeEmpty())
				})
			})

			When("there are environment variables that could be parsed", func() {
				It("uses the parsed values from the environment variables", func() {
					err := os.Setenv(keyAppEnv, string(customAppEnv))
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppPort, strconv.Itoa(customAppPort))
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppLoggingVerbosityGlobal, strconv.Itoa(customAppLoggingVerbosityGlobal))
					Expect(err).NotTo(HaveOccurred())

					err = os.Setenv(keyAppLoggingVerbosityModules, customAppLoggingVerbosityModulesEnvVar)
					Expect(err).NotTo(HaveOccurred())

					instance := config.New(config.WithEnvVars())

					Expect(instance.AppEnv).To(Equal(customAppEnv))
					Expect(instance.AppPort).To(Equal(customAppPort))
					Expect(instance.Logging.GetGlobalVerbosity()).To(Equal(customAppLoggingVerbosityGlobal))
					Expect(len(instance.Logging.GetModules())).To(Equal(len(customAppLoggingVerbosityModules)))
					for module, verbosity := range customAppLoggingVerbosityModules {
						Expect(instance.Logging.GetVerbosity(module)).To(Equal(verbosity))
					}
				})
			})
		})

		When("WithAppEnv is specified", func() {
			It("has an env having the value of the argument", func() {
				instance := config.New(config.WithAppEnv(customAppEnv))

				Expect(instance.AppEnv).To(Equal(customAppEnv))
			})
		})

		When("WithAppPort is specified", func() {
			It("has a port having the value of the argument", func() {
				instance := config.New(config.WithAppPort(customAppPort))

				Expect(instance.AppPort).To(Equal(customAppPort))
			})
		})

		When("WithLoggingGlobalVerbosity is specified", func() {
			It("has a global verbosity having the value of the argument", func() {
				instance := config.New(config.WithLoggingGlobalVerbosity(customAppLoggingVerbosityGlobal))

				Expect(instance.Logging.GetGlobalVerbosity()).To(Equal(customAppLoggingVerbosityGlobal))
				Expect(instance.Logging.GetVerbosity(moduleParent)).To(Equal(customAppLoggingVerbosityGlobal))
			})
		})

		When("WithLoggingModulesVerbosity is specified", func() {
			It("has a set of module verbosity having the value of the argument", func() {
				instance := config.New(config.WithLoggingModulesVerbosity(customAppLoggingVerbosityModules))

				Expect(len(instance.Logging.GetModules())).To(Equal(len(customAppLoggingVerbosityModules)))

				for module, verbosity := range customAppLoggingVerbosityModules {
					Expect(instance.Logging.GetVerbosity(module)).To(Equal(verbosity))
				}
			})
		})

	})

})
