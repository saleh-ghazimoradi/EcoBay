package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
)

func userRoutes(v1 fiber.Router, handler *handlers.UserHandler) {
	v1.Post("/register", handler.Register)
	//v1.Post("/login", handler.Login)
}
