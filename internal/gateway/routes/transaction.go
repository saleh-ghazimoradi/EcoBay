package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
)

func transactionRoute(v1 fiber.Router, transactionHandler *handlers.TransactionHandler, authentication helper.Authentication, userHandler *handlers.UserHandler) {
	secRoute := v1.Group("/buyer", authentication.Authorize)
	secRoute.Get("/payment", transactionHandler.MakePayment)
	secRoute.Get("/verify", transactionHandler.VerifyPayment)

	sellerRoute := v1.Group("/seller", authentication.AuthorizeSeller)
	sellerRoute.Get("/orders", userHandler.GetOrders)
}
