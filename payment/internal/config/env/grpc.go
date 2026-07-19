package env

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host            string        `env:"SERVICE_HOST,required"`
	Port            string        `env:"SERVICE_PORT,required"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func (gc *grpcConfig) Address() string {
	return net.JoinHostPort(gc.raw.Host, gc.raw.Port)
}

func (gc *grpcConfig) ShutdownTimeout() time.Duration {
	return gc.raw.ShutdownTimeout
}

func NewGRPCConfig() (*grpcConfig, error) {
	raw, err := env.ParseAs[grpcEnvConfig]()
	if err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &grpcConfig{
		raw: raw,
	}, nil
}
