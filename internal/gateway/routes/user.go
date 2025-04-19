package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
)

func userRoutes(v1 fiber.Router, handler *handlers.UserHandler, helper helper.Authentication) {
	publicRoutes := v1.Group("/")
	publicRoutes.Post("/register", handler.Register)
	publicRoutes.Post("/login", handler.Login)

	privateRoutes := publicRoutes.Group("/users", helper.Authorize)
	privateRoutes.Get("/profile", handler.GetProfile)
	privateRoutes.Get("/verify", handler.GetVerificationCode)
	privateRoutes.Post("/verify", handler.Verify)
	privateRoutes.Post("/become-seller", handler.BecomeSeller)
}
