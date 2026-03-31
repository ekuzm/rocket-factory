package config

import (
	"fmt"

	"github.com/ekuzm/rocket-factory/payment/internal/config/env"
	platformLogger "github.com/ekuzm/rocket-factory/platform/pkg/logger"
)

type GRPC interface {
	Address() string
}

type Logger interface {
	Level() string
	AsJSON() bool
}

type config struct {
	GRPC   GRPC
	Logger Logger
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

	app = &config{
		GRPC:   grpc,
		Logger: logger,
	}

	platformLogger.Init(logger.Level(), logger.AsJSON())

	return nil
}
