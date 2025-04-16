package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

var AppConfig *Config

type Config struct{}

func LoadConfig() error {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		slg.Logger.Error("error loading config", "error", err)
		return err
	}

	AppConfig = config

	return nil
}
