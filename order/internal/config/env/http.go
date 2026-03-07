package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type serviceConfig struct {
	Host string `env:"SERVICE_HOST, required"`
	Port string `env:"SERVICE_PORT, required"`
}

type httpEnvConfig struct {
	Order     serviceConfig `envPrefix:"ORDER_"`
	Inventory serviceConfig `envPrefix:"INVENTORY_"`
	Payment   serviceConfig `envPrefix:"PAYMENT_"`
}

type httpConfig struct {
	raw httpEnvConfig
}

func (hc *httpConfig) OrderAddress() string {
	return net.JoinHostPort(hc.raw.Order.Host, hc.raw.Order.Port)
}

func (hc *httpConfig) InventoryAddress() string {
	return net.JoinHostPort(hc.raw.Inventory.Host, hc.raw.Inventory.Port)
}

func (hc *httpConfig) PaymentAddress() string {
	return net.JoinHostPort(hc.raw.Payment.Host, hc.raw.Payment.Port)
}

func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("parsing .env: %w", err)
	}

	return &httpConfig{
		raw: raw,
	}, nil
}
