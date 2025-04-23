package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/notification"
	"github.com/saleh-ghazimoradi/EcoBay/pkg/payment"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	v1 := app.Group("/v1")

	/*---------- Dependencies ----------*/
	email, _ := notification.NewEmail(config.AppConfig.SMTP.Host, config.AppConfig.SMTP.Port, config.AppConfig.SMTP.UserName, config.AppConfig.SMTP.Password, config.AppConfig.SMTP.Sender)
	authentication := helper.NewAuth(config.AppConfig.AuthConfig.Secret)
	paymentClient := payment.NewPaymentsClient(config.AppConfig.Stripe.Secret, config.AppConfig.Stripe.SuccessUrl, config.AppConfig.Stripe.CancelUrl)

	/*---------- Repositories ----------*/
	catalogRepository := repository.NewCatalogRepository(db, db)
	userRepository := repository.NewUserRepository(db, db)
	transactionRepository := repository.NewTransactionRepository(db, db)

	/*---------- Services ----------*/
	catalogService := service.NewCatalogService(catalogRepository)
	userService := service.NewUserService(userRepository, catalogRepository, authentication, email)
	transactionService := service.NewTransactionService(transactionRepository, authentication)

	/*---------- Handlers ----------*/
	healthCheck := handlers.NewHealthCheckHandler()
	catalogHandler := handlers.NewCatalogHandler(catalogService, authentication)
	userHandler := handlers.NewUserHandler(userService, authentication)
	transactionHandler := handlers.NewTransactionHandler(transactionService, authentication, paymentClient, userService)

	healthCheckRoute(v1, healthCheck)
	catalogRoutes(v1, catalogHandler, authentication)
	userRoutes(v1, userHandler, authentication)
	transactionRoute(v1, transactionHandler, authentication, userHandler)
}
