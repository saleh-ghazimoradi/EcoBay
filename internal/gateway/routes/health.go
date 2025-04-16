package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
)

func healthCheckRoute(v1 fiber.Router, health *handlers.HealthCheckHandler) {
	v1.Get("/health", health.HealthCheck)
}
