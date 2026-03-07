package config

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/config/env"
)

type GRPC interface {
	Address() string
}

type Logger interface {
	Level() string
	AsJSON() bool
}

type Mongo interface {
	Name() string
	URI() string
	IsInit() bool
}

type config struct {
	GRPC   GRPC
	Logger Logger
	Mongo  Mongo
}

var app *config

func Setup() error {
	grpc, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongo, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	app = &config{
		GRPC:   grpc,
		Logger: logger,
		Mongo:  mongo,
	}

	return nil
}

func App() *config {
	return app
}
