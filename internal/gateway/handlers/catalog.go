package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
)

type CatalogHandler struct {
	catalogService service.CatalogService
	auth           helper.Authentication
}

func (c *CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	user := c.auth.GetCurrentUser(ctx)

	return successResponse(ctx, fiber.StatusCreated, "", user)
}

func (c *CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	return nil
}

func (c *CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {}

func (c *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {}

func NewCatalogHandler(catalogService service.CatalogService, auth helper.Authentication) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
		auth:           auth,
	}
}
