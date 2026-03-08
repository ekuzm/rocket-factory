package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host string `env:"SERVICE_HOST,required"`
	Port string `env:"SERVICE_PORT,required"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func (gc *grpcConfig) Address() string {
	return net.JoinHostPort(gc.raw.Host, gc.raw.Port)
}

func NewGRPCConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing envs: %w", err)
	}

	return &grpcConfig{
		raw: raw,
	}, nil
}
