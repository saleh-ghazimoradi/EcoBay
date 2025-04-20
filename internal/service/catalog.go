package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/internal/dto"
	"github.com/saleh-ghazimoradi/EcoBay/internal/repository"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
)

type CatalogService interface {
	CreateCategory(ctx context.Context, input *dto.Category) error
	EditCategory(ctx context.Context, id uint, input *dto.UpdateCategory) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	GetCategoryById(ctx context.Context, id uint) (*domain.Category, error)
	GetCategories(ctx context.Context) ([]*domain.Category, error)
	CreateProduct(ctx context.Context, input *dto.Product, user *domain.User) error
	EditProduct(ctx context.Context, id uint, input *dto.UpdateProduct, user *domain.User) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint, user *domain.User) error
	GetProductById(ctx context.Context, id uint) (*domain.Product, error)
	GetProducts(ctx context.Context) ([]*domain.Product, error)
	GetSellerProducts(ctx context.Context, id uint) ([]*domain.Product, error)
	UpdateStock(ctx context.Context, product *domain.Product) (*domain.Product, error)
}

type catalogService struct {
	catalogRepository repository.CatalogRepository
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

func (c *catalogService) CreateProduct(ctx context.Context, input *dto.Product, user *domain.User) error {
	if err := c.catalogRepository.CreateProduct(ctx, &domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		ImageUrl:    input.ImageUrl,
		UserId:      user.ID,
		Stock:       uint(input.Stock),
	}); err != nil {
		slg.Logger.Error("failed to create a product", "error", err.Error())
		return customErr.ErrsCreate
	}
	return nil
}

func (c *catalogService) EditProduct(ctx context.Context, id uint, input *dto.UpdateProduct, user *domain.User) (*domain.Product, error) {
	existProduct, err := c.GetProductById(ctx, id)
	if err != nil {
		slg.Logger.Error("failed to fetch a product", "error", err.Error())
		return nil, err
	}

	if existProduct.UserId != user.ID {
		return nil, errors.New("you don't have manage rights of this products")
	}

	if input.Name != nil {
		existProduct.Name = *input.Name
	}

	if input.Description != nil {
		existProduct.Description = *input.Description
	}

	if input.Price != nil {
		existProduct.Price = *input.Price
	}

	if input.CategoryId != nil {
		existProduct.CategoryId = *input.CategoryId
	}

	updatedProduct, err := c.catalogRepository.EditProduct(ctx, existProduct)
	if err != nil {
		slg.Logger.Error("failed to update a product", "error", err.Error())
		return nil, err
	}

	return updatedProduct, nil
}

func (c *catalogService) DeleteProduct(ctx context.Context, id uint, user *domain.User) error {
	existingProduct, err := c.GetProductById(ctx, id)
	if err != nil {
		slg.Logger.Error("failed to fetch a product", "error", err.Error())
		return customErr.ErrNotFound
	}

	if existingProduct.UserId != user.ID {
		return errors.New("you don't have manage rights of this products")
	}

	if err = c.catalogRepository.DeleteProduct(ctx, existingProduct); err != nil {
		slg.Logger.Error("failed to delete a product", "error", err.Error())
		return err
	}

	return nil
}

func (c *catalogService) GetProductById(ctx context.Context, id uint) (*domain.Product, error) {
	product, err := c.catalogRepository.FindProductById(ctx, id)
	if err != nil {
		slg.Logger.Error("failed to fetch a product", "error", err.Error())
		return nil, err
	}

	return product, nil
}

func (c *catalogService) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	products, err := c.catalogRepository.FindProducts(ctx)
	if err != nil {
		return nil, errors.New("products not found")
	}

	return products, nil
}

func (c *catalogService) GetSellerProducts(ctx context.Context, id uint) ([]*domain.Product, error) {
	products, err := c.catalogRepository.FindSellerProducts(ctx, id)
	if err != nil {
		return nil, errors.New("products not found")
	}
	return products, nil
}

func (c *catalogService) UpdateStock(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	eProduct, err := c.GetProductById(ctx, product.ID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if eProduct.UserId != product.UserId {
		return nil, errors.New("you do not have manage rights of this products")
	}

	eProduct.Stock = product.Stock

	editProduct, err := c.catalogRepository.EditProduct(ctx, eProduct)
	if err != nil {
		slg.Logger.Error("failed to update a product", "error", err.Error())
		return nil, err
	}

	return editProduct, nil
}

func NewCatalogService(catalogRepository repository.CatalogRepository) CatalogService {
	return &catalogService{
		catalogRepository: catalogRepository,
	}
}
