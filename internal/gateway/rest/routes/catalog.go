package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/gateway/rest/handlers"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"gorm.io/gorm"
)

func CatalogRoutes(app *fiber.App, db *gorm.DB) {
	catalogRepository := repository.NewCatalogRepository(db)
	catalogService := service.NewCatalogService(catalogRepository)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	// Listing Products and Categories
	app.Get("/products", catalogHandler.GetProducts)
	app.Get("/products/:id", catalogHandler.GetProduct)
	app.Get("/categories", catalogHandler.GetCategories)
	app.Get("/categories/:id", catalogHandler.GetCategoryById)

	sellerRoutes := app.Group("/seller", authService.AuthorizeSeller)
	// Categories
	sellerRoutes.Post("/categories", catalogHandler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", catalogHandler.EditCategory)
	sellerRoutes.Delete("categories/:id", catalogHandler.DeleteCategory)

	// Products
	sellerRoutes.Post("/products", catalogHandler.CreateProducts)
	sellerRoutes.Get("/products", catalogHandler.GetProducts)
	sellerRoutes.Get("/products/:id", catalogHandler.GetProduct)
	sellerRoutes.Put("/products/:id", catalogHandler.EditProduct)
	sellerRoutes.Patch("/products/:id", catalogHandler.UpdateStock)
	sellerRoutes.Delete("/products/:id", catalogHandler.DeleteProduct)

}
