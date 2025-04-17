package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	v1 := app.Group("/v1")

	userRepository := repository.NewUserRepository(db, db)

	authentication := helper.NewAuth(config.AppConfig.AuthConfig.Secret)
	userService := service.NewUserService(userRepository, authentication)

	healthCheck := handlers.NewHealthCheckHandler()
	userHandler := handlers.NewUserHandler(userService, authentication)
	healthCheckRoute(v1, healthCheck)
	userRoutes(v1, userHandler, authentication)
}
