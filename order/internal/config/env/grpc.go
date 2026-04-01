package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type adapterConfig struct {
	Host string `env:"ADAPTER_HOST,required"`
	Port string `env:"ADAPTER_PORT,required"`
}

type grpcEnvConfig struct {
	Inventory adapterConfig `envPrefix:"INVENTORY_"`
	Payment   adapterConfig `envPrefix:"PAYMENT_"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func (gc *grpcConfig) InventoryAddress() string {
	return net.JoinHostPort(gc.raw.Inventory.Host, gc.raw.Inventory.Port)
}

func (gc *grpcConfig) PaymentAddress() string {
	return net.JoinHostPort(gc.raw.Payment.Host, gc.raw.Payment.Port)
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
