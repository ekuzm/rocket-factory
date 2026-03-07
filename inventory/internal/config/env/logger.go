package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJSON bool   `env:"LOGGER_AS_JSON,required"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func (lc *loggerConfig) Level() string {
	return lc.raw.Level
}

func (lc *loggerConfig) AsJSON() bool {
	return lc.raw.AsJSON
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing .env: %w", err)
	}

	return &loggerConfig{
		raw: raw,
	}, nil
}
