package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type GRPCServerConfig struct {
	Address string `env:"address"`
}

type Config struct {
	GRPCServerConfig
}

func New() *Config {
	return &Config{
		GRPCServerConfig: GRPCServerConfig{
			Address: "0.0.0.0:8080",
		},
	}
}

func (cfg *Config) Parse() error {
	if err := env.Parse(&cfg.GRPCServerConfig); err != nil {
		return errors.Wrap(err, "parse GRPCServerConfig error")
	}

	return nil
}
