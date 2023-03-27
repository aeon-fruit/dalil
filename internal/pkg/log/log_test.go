package log_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
	"github.com/aeon-fruit/dalil.git/internal/pkg/log"
	logrMock "github.com/aeon-fruit/dalil.git/test/mocks/logr"
)

var _ = Describe("Log", func() {

	const (
		logMsg           = "log message"
		loggerNameParent = "parent"
		loggerNameModule = "module"
		fileName         = "log_test.go"
	)

	var (
		customErr     error
		keysAndValues []any

		appConfig config.AppConfig
		output    *bytes.Buffer

		mockCtrl *gomock.Controller
	)

	BeforeEach(func() {
		customErr = fmt.Errorf("custom error")
		keysAndValues = []any{"keyInt", 10., "keyStr", "value", "keyBool", true}

		appConfig = config.New()
		output = &bytes.Buffer{}

		mockCtrl = gomock.NewController(GinkgoT())
	})

	Describe("New", func() {

		var logger log.Logger

		BeforeEach(func() {
			logger = log.New(appConfig, output)
		})

		When("called with the default config", func() {
			It("returns a non-nil Logger instance", func() {
				Expect(logger).NotTo(BeZero())
			})

			It("is enabled", func() {
				Expect(logger.Enabled()).To(BeTrue())
			})

			It("logs error", func() {
				logger.Error(customErr, logMsg, keysAndValues...)

				out := output.String()
				Expect(out).To(ContainSubstring(fileName))
				Expect(out).To(ContainSubstring(logMsg))
				Expect(out).To(ContainSubstring(customErr.Error()))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[0]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[1]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[2]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[3]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[4]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[5]))
			})

			It("has a non-nil LogSink", func() {
				Expect(logger.GetSink()).NotTo(BeNil())
			})

			It("logs info", func() {
				logger.Info(logMsg, keysAndValues...)

				out := output.String()
				Expect(out).To(ContainSubstring(fileName))
				Expect(out).To(ContainSubstring(logMsg))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[0]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[1]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[2]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[3]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[4]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[5]))
			})

			It("is disabled for higher values of Verbosity threshold", func() {
				l := logger.V(3)

				Expect(l.Enabled()).To(BeFalse())

				l.Info(logMsg, keysAndValues...)

				Expect(output.String()).To(BeEmpty())
			})

			It("sets the specified call depth", func() {
				l := logger.WithCallDepth(1)

				Expect(l.GetSink()).To(BeAssignableToTypeOf(logger.GetSink()))
			})

			It("adds the specified name to the log entry", func() {
				l := logger.WithName(loggerNameParent)
				l.Info(logMsg, keysAndValues...)

				out := output.String()
				Expect(out).To(ContainSubstring(fileName))
				Expect(out).To(ContainSubstring(logMsg))
				Expect(out).To(ContainSubstring(loggerNameParent))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[0]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[1]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[2]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[3]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[4]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[5]))

				output.Reset()
				l.WithName(loggerNameModule).Info(logMsg, keysAndValues...)

				Expect(output.String()).To(ContainSubstring("%v/%v", loggerNameParent, loggerNameModule))
			})

			It("uses the specified LogSink", func() {
				sink := logrMock.NewMockLogSink(mockCtrl)
				sink.EXPECT().Enabled(gomock.Any()).Return(true)
				sink.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

				l := logger.WithSink(sink)
				l.Info(logMsg, keysAndValues...)
			})

			It("adds the specified values to the log entry", func() {
				l := logger.WithValues(keysAndValues...)
				l.Info(logMsg)

				out := output.String()
				Expect(out).To(ContainSubstring(fileName))
				Expect(out).To(ContainSubstring(logMsg))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[0]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[1]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[2]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[3]))
				Expect(out).To(ContainSubstring("%v=", keysAndValues[4]))
				Expect(out).To(ContainSubstring("%v", keysAndValues[5]))
			})

			It("prints formatted strings", func() {
				expected := fmt.Sprintf("%v: %v=%v, %v=%v, %v=%v", logMsg, keysAndValues[0], keysAndValues[1],
					keysAndValues[2], keysAndValues[3],
					keysAndValues[4], keysAndValues[5])
				logger.Print(expected)

				Expect(output.String()).To(ContainSubstring(expected))
			})
		})

		When("called with custom verbosity config", func() {

			When("only the global verbosity is specified", func() {
				It("uses the global verbosity by default", func() {
					appConfig = config.New(config.WithLoggingGlobalVerbosity(1))
					logger = log.New(appConfig, output)

					logger.Info(logMsg)

					Expect(output.String()).NotTo(BeEmpty())

					output.Reset()
					logger.V(2).Info(logMsg)

					Expect(output.String()).To(BeEmpty())
				})

				It("uses the global verbosity for modules", func() {
					appConfig = config.New(config.WithLoggingGlobalVerbosity(1))
					logger = log.New(appConfig, output)
					l := logger.WithName(loggerNameParent)

					l.Info(logMsg)

					Expect(output.String()).NotTo(BeEmpty())

					output.Reset()
					l.V(2).Info(logMsg)

					Expect(output.String()).To(BeEmpty())
				})
			})

			When("the module verbosity is specified", func() {
				It("uses the specific verbosity", func() {
					appConfig = config.New(config.WithLoggingModulesVerbosity(map[string]int{loggerNameParent: 1}))
					logger = log.New(appConfig, output)
					l := logger.WithName(loggerNameParent)

					l.Info(logMsg)

					Expect(output.String()).NotTo(BeEmpty())

					output.Reset()
					l.V(2).Info(logMsg)

					Expect(output.String()).To(BeEmpty())
				})
			})

		})

		When("called with non-local app environment", func() {
			It("outputs JSON log entries", func() {
				appConfig = config.New(config.WithAppEnv(config.AppEnvDev))
				logger = log.New(appConfig, output)

				logger.Info(logMsg, keysAndValues...)

				var out map[string]any
				err := json.Unmarshal(output.Bytes(), &out)

				Expect(err).ToNot(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out["message"]).To(Equal(logMsg))
				Expect(out["level"]).To(Equal("info"))
				Expect(out["v"]).To(Equal(0.))
				Expect(out["caller"]).To(ContainSubstring(fileName))
				Expect(out["time"]).NotTo(BeEmpty())
				Expect(out[keysAndValues[0].(string)]).To(Equal(keysAndValues[1]))
				Expect(out[keysAndValues[2].(string)]).To(Equal(keysAndValues[3]))
				Expect(out[keysAndValues[4].(string)]).To(Equal(keysAndValues[5]))
			})
		})

	})

})
