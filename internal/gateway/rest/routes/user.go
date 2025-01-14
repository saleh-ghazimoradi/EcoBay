package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/rest/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
)

func UserRoutes(app *fiber.App) {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	app.Post("/v1/users/register", userHandler.Register)
	app.Post("/v1/users", userHandler.Login)

	app.Get("/verify", userHandler.GetVerificationCode)
	app.Post("/verify", userHandler.Verify)
	app.Post("/profile", userHandler.CreateProfile)
	app.Get("/profile", userHandler.GetProfile)

	app.Post("/cart", userHandler.AddToCart)
	app.Get("/cart", userHandler.GetCart)
	app.Get("/order", userHandler.GetOrders)
	app.Get("/order/:id", userHandler.GetOrder)
	app.Post("/become-seller", userHandler.BecomeSeller)
}
