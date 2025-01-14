package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
)

func Server() error {
	app := fiber.New()

	if err := app.Listen(config.AppConfig.ServerConfig.Port); err != nil {
		return err
	}
	
	return nil
}
