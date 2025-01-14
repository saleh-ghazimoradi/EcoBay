package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

var AppConfig *Config

type Config struct {
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Port    string `env:"SERVER_PORT,required"`
	Version string `env:"SERVER_VERSION,required"`
}

func LoadConfig() error {
	if err := godotenv.Load("app.env"); err != nil {
		log.Fatal("Error loading app.env file")
	}

	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatal("Error parsing config")
	}

	serverConfig := &ServerConfig{}
	if err := env.Parse(serverConfig); err != nil {
		log.Fatal("Error parsing config")
	}
	config.ServerConfig = *serverConfig

	AppConfig = config

	return nil
}
