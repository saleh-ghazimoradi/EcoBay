package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	v1 := app.Group("/v1")

	userRepository := repository.NewUserRepository(db, db)
	userService := service.NewUserService(userRepository)

	healthCheck := handlers.NewHealthCheckHandler()
	userHandler := handlers.NewUserHandler(userService)
	healthCheckRoute(v1, healthCheck)
	userRoutes(v1, userHandler)
}
