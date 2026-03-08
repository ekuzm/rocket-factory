package config

import (
	"fmt"

	"github.com/ekuzm/rocket-factory/order/internal/config/env"
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

func Setup() error {
	http, err := env.NewHTTPConfig()
	if err != nil {
		return fmt.Errorf("create http config: %w", err)
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("create logger config: %w", err)
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return fmt.Errorf("create postgres config: %w", err)
	}

	app = &config{
		HTTP:     http,
		Logger:   logger,
		Postgres: postgres,
	}

	return nil
}
