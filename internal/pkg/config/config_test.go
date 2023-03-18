package config_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
)

const (
	keyAppPort = "APP_PORT"

	defaultAppPort = 8080
	customAppPort  = 12345
)

var _ = Describe("Config", func() {
	Describe("New", func() {
		It("has defaults", func() {
			instance := config.New()

			Expect(instance.AppPort).To(Equal(defaultAppPort))
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
					err := os.Setenv(keyAppPort, "non integer value")
					Expect(err).NotTo(HaveOccurred())

					instance := config.New(config.WithEnvVars())

					Expect(instance.AppPort).To(Equal(defaultAppPort))
				})
			})

			When("there are environment variables that could be parsed", func() {
				It("keeps the defaults", func() {
					err := os.Setenv(keyAppPort, strconv.Itoa(customAppPort))
					Expect(err).NotTo(HaveOccurred())

					instance := config.New(config.WithEnvVars())

					Expect(instance.AppPort).To(Equal(customAppPort))
				})
			})
		})

		When("WithAppPort is specified", func() {
			It("has a port having the value of the argument", func() {
				instance := config.New(config.WithAppPort(customAppPort))

				Expect(instance.AppPort).To(Equal(customAppPort))
			})
		})
	})
})
