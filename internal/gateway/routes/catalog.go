package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
)

func catalogRoutes(v1 fiber.Router, handler *handlers.CatalogHandler, auth helper.Authentication) {
	publicRoutes := v1.Group("/")
	publicRoutes.Get("/products", handler.GetProducts)
	publicRoutes.Get("/products/:id", handler.GetProduct)
	publicRoutes.Get("/categories", handler.GetCategories)
	publicRoutes.Get("/categories/:id", handler.GetCategoryById)

	sellerRoutes := v1.Group("/seller", auth.AuthorizeSeller)
	sellerRoutes.Post("/categories", handler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", handler.EditCategory)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategory)

	sellerRoutes.Get("/products", handler.GetProducts)
	sellerRoutes.Get("products/:id", handler.GetProduct)
	sellerRoutes.Post("/products", handler.CreateProducts)
	sellerRoutes.Patch("/products/:id", handler.UpdateStock)
	sellerRoutes.Put("/products/:id", handler.EditProduct)
	sellerRoutes.Delete("/products/:id", handler.DeleteProduct)
}
