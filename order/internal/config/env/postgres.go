package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	Name     string `env:"DB_NAME,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func (pc *postgresConfig) URI() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		pc.raw.User,
		pc.raw.Password,
		pc.raw.Host,
		pc.raw.Port,
		pc.raw.Name,
	)
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing .env: %w", err)
	}

	return &postgresConfig{
		raw: raw,
	}, nil
}
