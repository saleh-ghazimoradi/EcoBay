package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	v1 := app.Group("/v1")

	healthCheck := handlers.NewHealthCheckHandler()
	healthCheckRoute(v1, healthCheck)
}
