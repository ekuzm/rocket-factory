package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type serviceConfig struct {
	Host string `env:"SERVICE_HOST,required"`
	Port string `env:"SERVICE_PORT,required"`
}

type httpEnvConfig struct {
	Order serviceConfig `envPrefix:"ORDER_"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func (hc *httpConfig) OrderAddress() string {
	return net.JoinHostPort(hc.raw.Order.Host, hc.raw.Order.Port)
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &httpConfig{
		raw: raw,
	}, nil
}
