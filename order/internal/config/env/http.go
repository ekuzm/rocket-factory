package env

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type serviceConfig struct {
	Host string `env:"SERVICE_HOST,required"`
	Port string `env:"SERVICE_PORT,required"`
}

type httpEnvConfig struct {
	Order             serviceConfig `envPrefix:"ORDER_"`
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT,required"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT,required"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func (hc *httpConfig) OrderAddress() string {
	return net.JoinHostPort(hc.raw.Order.Host, hc.raw.Order.Port)
}

func (hc *httpConfig) ShutdownTimeout() time.Duration {
	return hc.raw.ShutdownTimeout
}

func (hc *httpConfig) RequestTimeout() time.Duration {
	return hc.raw.RequestTimeout
}

func (hc *httpConfig) ReadHeaderTimeout() time.Duration {
	return hc.raw.ReadHeaderTimeout
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
