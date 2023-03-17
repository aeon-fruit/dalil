package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	keyAppPort = "APP_PORT"

	defaultAppPort = 8080
)

type AppConfig struct {
	AppPort int
}

type AppConfigOption func(*AppConfig)

func New(options ...AppConfigOption) AppConfig {
	instance := AppConfig{
		AppPort: defaultAppPort,
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

			appConfig.AppPort = getEnvVarInt(keyAppPort, defaultAppPort)
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

func getEnvVarInt(key string, defaultValue int) int {
	if value, found := os.LookupEnv(key); found {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
