package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/service"
	"strconv"
)

type CatalogHandler struct {
	catalogService service.CatalogService
	auth           helper.Authentication
}

func (c *CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	payload := dto.Category{}

	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, err)
	}

	if err := c.catalogService.CreateCategory(ctx.Context(), &payload); err != nil {
		return serverErrorResponse(ctx, err)
	}
	return successResponse(ctx, fiber.StatusCreated, "category created successfully", nil)
}

func (c *CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)
	uintId := uint(id)

	category, err := c.catalogService.GetCategoryById(ctx.Context(), uintId)
	if err != nil {
		return notFoundResponse(ctx)
	}
	return successResponse(ctx, fiber.StatusOK, "category successfully retrieved", category)
}

func (c *CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	categories, err := c.catalogService.GetCategories(ctx.Context())
	if err != nil {
		return notFoundResponse(ctx)
	}
	return successResponse(ctx, fiber.StatusOK, "categories successfully retrieved", categories)
}

func (c *CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)
	uintId := uint(id)

	payload := dto.UpdateCategory{}

	if err := ctx.BodyParser(&payload); err != nil {
		return badRequestResponse(ctx, err)
	}

	updateCat, err := c.catalogService.EditCategory(ctx.Context(), uintId, &payload)
	if err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusOK, "category successfully edited", updateCat)
}

func (c *CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)
	uintId := uint(id)
	if err := c.catalogService.DeleteCategory(ctx.Context(), uintId); err != nil {
		return serverErrorResponse(ctx, err)
	}

	return successResponse(ctx, fiber.StatusNoContent, "category successfully deleted", nil)
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

func NewCatalogHandler(catalogService service.CatalogService, auth helper.Authentication) *CatalogHandler {
	return &CatalogHandler{
		catalogService: catalogService,
		auth:           auth,
	}
}
