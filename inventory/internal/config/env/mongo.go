package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	Name     string `env:"DB_NAME,required"`
	Auth     string `env:"DB_AUTH,required"`
	IsInit   bool   `env:"DB_IS_INIT,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func (mc *mongoConfig) Name() string {
	return mc.raw.Name
}

func (mc *mongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", mc.raw.User, mc.raw.Password, mc.raw.Host, mc.raw.Port, mc.raw.Name, mc.raw.Auth)
}

func (mc *mongoConfig) IsInit() bool {
	return mc.raw.IsInit
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &mongoConfig{
		raw: raw,
	}, nil
}
