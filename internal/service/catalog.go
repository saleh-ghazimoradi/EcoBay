package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/helper"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

type CatalogService interface {
	CreateCategory(ctx context.Context, input *dto.Category) error
	EditCategory(ctx context.Context, id uint, input *dto.UpdateCategory) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryById(ctx context.Context, id uint) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]*domain.Category, error)
}

type catalogService struct {
	catalogRepository repository.CatalogRepository
	authentication    helper.Authentication
}

func (c *catalogService) CreateCategory(ctx context.Context, input *dto.Category) error {
	if err := c.catalogRepository.CreateCategory(ctx, &domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	}); err != nil {
		slg.Logger.Error("error creating category", "error", err.Error())
		return errors.New("error creating category: " + err.Error())
	}
	return nil
}

func (c *catalogService) GetCategoryById(ctx context.Context, id uint) (*domain.Category, error) {
	category, err := c.catalogRepository.FindCategoryById(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (c *catalogService) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := c.catalogRepository.FindCategories(ctx)
	if err != nil {
		return nil, errors.New("categories not found")
	}
	return categories, nil
}

func (c *catalogService) EditCategory(ctx context.Context, id uint, input *dto.UpdateCategory) (*domain.Category, error) {
	existingCategory, err := c.GetCategoryById(ctx, id)
	if err != nil {
		return nil, customErr.ErrNotFound
	}

	if input.Name != nil {
		existingCategory.Name = *input.Name
	}

	if input.ParentId != nil {
		existingCategory.ParentId = *input.ParentId
	}

	if input.ImageUrl != nil {
		existingCategory.ImageUrl = *input.ImageUrl
	}

	if input.DisplayOrder != nil {
		existingCategory.DisplayOrder = *input.DisplayOrder
	}

	category, err := c.catalogRepository.EditCategory(ctx, existingCategory)
	if err != nil {
		return nil, customErr.ErrUpdate
	}

	return category, nil
}

func (c *catalogService) DeleteCategory(ctx context.Context, id uint) error {
	if err := c.catalogRepository.DeleteCategory(ctx, id); err != nil {
		slg.Logger.Error("failed to delete category", "id", id, "error", err)
		return customErr.ErrDelete
	}
	return nil
}

func NewCatalogService(catalogRepository repository.CatalogRepository, authentication helper.Authentication) CatalogService {
	return &catalogService{
		catalogRepository: catalogRepository,
		authentication:    authentication,
	}
}
