package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/notification"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	v1 := app.Group("/v1")

	/*---------- Dependencies ----------*/
	email, _ := notification.NewEmail(config.AppConfig.SMTP.Host, config.AppConfig.SMTP.Port, config.AppConfig.SMTP.UserName, config.AppConfig.SMTP.Password, config.AppConfig.SMTP.Sender)
	authentication := helper.NewAuth(config.AppConfig.AuthConfig.Secret)

	/*---------- Repositories ----------*/
	userRepository := repository.NewUserRepository(db, db)

	/*---------- Services ----------*/
	userService := service.NewUserService(userRepository, authentication, email)

	/*---------- Handlers ----------*/
	healthCheck := handlers.NewHealthCheckHandler()
	userHandler := handlers.NewUserHandler(userService, authentication)

	healthCheckRoute(v1, healthCheck)
	userRoutes(v1, userHandler, authentication)
}
