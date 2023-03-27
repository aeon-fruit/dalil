package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppEnv string

const (
	AppEnvLocal   = AppEnv("local")
	AppEnvDev     = AppEnv("dev")
	AppEnvNonProd = AppEnv("nonprod")
	AppEnvProd    = AppEnv("prod")
)

const (
	keyAppEnv     = "APP_ENV"
	defaultAppEnv = AppEnvLocal

	keyAppPort     = "APP_PORT"
	defaultAppPort = 8080

	keyAppLoggingVerbosityGlobal     = "APP_LOGGING_VERBOSITY_GLOBAL"
	keyAppLoggingVerbosityModules    = "APP_LOGGING_VERBOSITY_MODULES"
	defaultAppLoggingVerbosityGlobal = 0
)

type LoggingConfig struct {
	globalVerbosity  int
	modulesVerbosity map[string]int
}

func (lc LoggingConfig) GetGlobalVerbosity() int {
	return lc.globalVerbosity
}

func (lc LoggingConfig) GetVerbosity(moduleName string) int {
	if verbosity, found := lc.modulesVerbosity[moduleName]; found {
		return verbosity
	}
	return lc.globalVerbosity
}

func (lc LoggingConfig) GetModules() (modules []string) {
	for module := range lc.modulesVerbosity {
		modules = append(modules, module)
	}
	return
}

type AppConfig struct {
	AppEnv  AppEnv
	AppPort int
	Logging LoggingConfig
}

type AppConfigOption func(*AppConfig)

func New(options ...AppConfigOption) AppConfig {
	instance := AppConfig{
		AppEnv:  defaultAppEnv,
		AppPort: defaultAppPort,
		Logging: LoggingConfig{
			globalVerbosity: defaultAppLoggingVerbosityGlobal,
		},
	}

	for _, option := range options {
		if option != nil {
			option(&instance)
		}
	}

	return instance
}

func WithEnvVars() AppConfigOption {
	return func(appConfig *AppConfig) {
		if appConfig != nil {
			godotenv.Load()

			appConfig.AppEnv = getAppEnv()
			appConfig.AppPort = getEnvVarInt(keyAppPort, defaultAppPort)
			appConfig.Logging.globalVerbosity = getEnvVarInt(keyAppLoggingVerbosityGlobal, defaultAppLoggingVerbosityGlobal)
			appConfig.Logging.modulesVerbosity = getEnvVarInts(keyAppLoggingVerbosityModules)
		}
	}
}

func WithAppEnv(appEnv AppEnv) AppConfigOption {
	return func(appConfig *AppConfig) {
		if appConfig != nil {
			appConfig.AppEnv = appEnv
		}
	}
}

func WithAppPort(appPort int) AppConfigOption {
	return func(appConfig *AppConfig) {
		if appConfig != nil {
			appConfig.AppPort = appPort
		}
	}
}

func WithLoggingGlobalVerbosity(loggingGlobalVerbosity int) AppConfigOption {
	return func(appConfig *AppConfig) {
		if appConfig != nil {
			appConfig.Logging.globalVerbosity = loggingGlobalVerbosity
		}
	}
}

func WithLoggingModulesVerbosity(loggingModulesVerbosity map[string]int) AppConfigOption {
	return func(appConfig *AppConfig) {
		if appConfig != nil {
			appConfig.Logging.modulesVerbosity = loggingModulesVerbosity
		}
	}
}

func getAppEnv() AppEnv {
	if value, found := os.LookupEnv(keyAppEnv); found {
		switch value {
		case string(AppEnvLocal), string(AppEnvDev), string(AppEnvNonProd), string(AppEnvProd):
			return AppEnv(value)
		}
	}
	return defaultAppEnv
}

func getEnvVarInt(key string, defaultValue int) int {
	if value, found := os.LookupEnv(key); found {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvVarInts(key string) map[string]int {
	if envVar, found := os.LookupEnv(key); found {
		pairs := strings.Split(envVar, ",")

		intsMap := map[string]int{}
		for _, pair := range pairs {
			if key, intValue, err := parsePair(pair); err == nil {
				intsMap[key] = intValue
			}
		}

		if len(intsMap) > 0 {
			return intsMap
		}
	}
	return nil
}

func parsePair(pair string) (string, int, error) {
	err := fmt.Errorf("error")
	pairTokens := strings.Split(pair, "=")
	if len(pairTokens) != 2 {
		return "", 0, err
	}

	key, value := strings.TrimSpace(pairTokens[0]), strings.TrimSpace(pairTokens[1])
	if len(key) == 0 {
		return "", 0, err
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return key, intValue, nil
	}
	return "", 0, err
}
