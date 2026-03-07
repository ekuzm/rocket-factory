package config

import (
	"fmt"

	"github.com/ekuzm/rocket-factory/order/internal/config/env"
	"github.com/joho/godotenv"
)

type Logger interface {
	Level() string
	AsJSON() bool
}

type Postgres interface {
	URI() string
}

type HTTP interface {
	OrderAddress() string
	InventoryAddress() string
	PaymentAddress() string
}

type config struct {
	HTTP     HTTP
	Logger   Logger
	Postgres Postgres
}

var app *config

func App() *config {
	return app
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil {
		return fmt.Errorf("load env file from %v: %w", path, err)
	}

	http, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	app = &config{
		HTTP:     http,
		Logger:   logger,
		Postgres: postgres,
	}

	return nil
}
