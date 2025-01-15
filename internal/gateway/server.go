package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/rest/routes"
	"github.com/saleh-ghazimoradi/EcoBay/utils"
)

func Server() error {
	app := fiber.New()

	db, err := utils.DBConnection(utils.DBMigrator)
	if err != nil {
		return err
	}

	routes.Health(app)
	routes.UserRoutes(app, db)

	if err = app.Listen(config.AppConfig.ServerConfig.Port); err != nil {
		return err
	}

	return nil
}
