package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/rest/handlers"
)

func Health(app *fiber.App) {
	health := handlers.NewHealthHandler()
	app.Get("/v1/health", health.Health)
}
