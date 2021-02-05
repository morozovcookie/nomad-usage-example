package config

import (
	env "github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type HTTPServerConfig struct {
	Address string `env:"HTTP_SERVER_ADDRESS"`
}

type Config struct {
	HTTPServerConfig
}

func New() *Config {
	return &Config{
		HTTPServerConfig: HTTPServerConfig{
			Address: "0.0.0.0:8080",
		},
	}
}

func (cfg *Config) Parse() error {
	if err := env.Parse(&cfg.HTTPServerConfig); err != nil {
		return errors.Wrap(err, "parse HTTPServerConfig error")
	}

	return nil
}
