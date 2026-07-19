package config

import (
	"fmt"
	"time"

	"github.com/ekuzm/rocket-factory/payment/internal/config/env"
)

type GRPC interface {
	Address() string
	ShutdownTimeout() time.Duration
}

type Logger interface {
	Level() string
	AsJSON() bool
}

type Tracer interface {
	ServiceName() string
	Endpoint() string
	ServiceVersion() string
	DeploymentEnvironment() string
	Compressor() string
	RetryEnabled() bool
	RetryInitialInterval() time.Duration
	RetryMaxInterval() time.Duration
	RetryMaxElapsedTime() time.Duration
	Timeout() time.Duration
}

type config struct {
	GRPC   GRPC
	Logger Logger
	Tracer Tracer
}

var app *config

func App() *config {
	return app
}

func Setup() error {
	grpc, err := env.NewGRPCConfig()
	if err != nil {
		return fmt.Errorf("create grpc config: %w", err)
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("create logger config: %w", err)
	}

	tracer, err := env.NewTracerConfig()
	if err != nil {
		return fmt.Errorf("create tracer config: %w", err)
	}

	app = &config{
		GRPC:   grpc,
		Logger: logger,
		Tracer: tracer,
	}

	return nil
}
