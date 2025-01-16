package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/rest/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/notification"
	"gorm.io/gorm"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
	userRepository := repository.NewUserRepository(db)
	authService := helper.NewAuth(config.AppConfig.AppSecret.Secret)
	notificationService := notification.NewNotificationsClient()
	userService := service.NewUserService(userRepository, authService, notificationService)
	userHandler := handlers.NewUserHandler(userService)

	pubRoutes := app.Group("/users")
	pubRoutes.Post("/register", userHandler.Register)
	pubRoutes.Post("/login", userHandler.Login)

	// Private routes
	pvtRoutes := pubRoutes.Group("/", authService.Authorize)
	pvtRoutes.Get("/verify", userHandler.GetVerificationCode)
	pvtRoutes.Post("/verify", userHandler.Verify)
	pvtRoutes.Post("/profile", userHandler.CreateProfile)
	pvtRoutes.Get("/profile", userHandler.GetProfile)
	pvtRoutes.Post("/cart", userHandler.AddToCart)
	pvtRoutes.Get("/cart", userHandler.GetCart)
	pvtRoutes.Get("/order", userHandler.GetOrders)
	pvtRoutes.Get("/order/:id", userHandler.GetOrder)
	pvtRoutes.Post("/become-seller", userHandler.BecomeSeller)
}
