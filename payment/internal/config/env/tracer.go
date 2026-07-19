package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type tracerEnvConfig struct {
	ServiceName           string        `env:"OTEL_SERVICE_NAME,required"`
	Endpoint              string        `env:"OTEL_EXPORTER_OTLP_ENDPOINT,required"`
	ServiceVersion        string        `env:"OTEL_SERVICE_VERSION,required"`
	DeploymentEnvironment string        `env:"DEPLOYMENT_ENVIRONMENT,required"`
	Compressor            string        `env:"OTEL_EXPORTER_OTLP_COMPRESSION,required"`
	RetryEnabled          bool          `env:"OTEL_EXPORTER_OTLP_RETRY_ENABLED,required"`
	RetryInitialInterval  time.Duration `env:"OTEL_EXPORTER_OTLP_RETRY_INITIAL_INTERVAL,required"`
	RetryMaxInterval      time.Duration `env:"OTEL_EXPORTER_OTLP_RETRY_MAX_INTERVAL,required"`
	RetryMaxElapsedTime   time.Duration `env:"OTEL_EXPORTER_OTLP_RETRY_MAX_ELAPSED_TIME,required"`
	Timeout               time.Duration `env:"OTEL_EXPORTER_OTLP_TIMEOUT,required"`
}

type tracerConfig struct {
	raw tracerEnvConfig
}

func (tc *tracerConfig) ServiceName() string {
	return tc.raw.ServiceName
}

func (tc *tracerConfig) Endpoint() string {
	return tc.raw.Endpoint
}

func (tc *tracerConfig) ServiceVersion() string {
	return tc.raw.ServiceVersion
}

func (tc *tracerConfig) DeploymentEnvironment() string {
	return tc.raw.DeploymentEnvironment
}

func (tc *tracerConfig) Compressor() string {
	return tc.raw.Compressor
}

func (tc *tracerConfig) RetryEnabled() bool {
	return tc.raw.RetryEnabled
}

func (tc *tracerConfig) RetryInitialInterval() time.Duration {
	return tc.raw.RetryInitialInterval
}

func (tc *tracerConfig) RetryMaxInterval() time.Duration {
	return tc.raw.RetryMaxInterval
}

func (tc *tracerConfig) RetryMaxElapsedTime() time.Duration {
	return tc.raw.RetryMaxElapsedTime
}

func (tc *tracerConfig) Timeout() time.Duration {
	return tc.raw.Timeout
}

func NewTracerConfig() (*tracerConfig, error) {
	var raw tracerEnvConfig

	err := env.Parse(&raw)
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &tracerConfig{raw: raw}, nil
}
