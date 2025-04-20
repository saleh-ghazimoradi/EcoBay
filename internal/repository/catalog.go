package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/EcoBay/internal/customErr"
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	FindCategories(ctx context.Context) ([]*domain.Category, error)
	FindCategoryById(ctx context.Context, id uint) (*domain.Category, error)
	EditCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uint) error
	//CreateProduct(ctx context.Context, product *domain.Product) error
	//FindProducts(ctx context.Context) ([]*domain.Product, error)
	//FindProductById(ctx context.Context, id uint) (*domain.Product, error)
	//FindSellerProducts(ctx context.Context, id uint) ([]*domain.Product, error)
	//EditProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	//DeleteProduct(ctx context.Context, product *domain.Product) error
}

type catalogRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (c *catalogRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	if err := c.dbWrite.WithContext(ctx).Create(&category).Error; err != nil {
		slg.Logger.Error("create category", "error", err)
		return customErr.ErrsCreate
	}
	return nil
}

func (c *catalogRepository) FindCategories(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	if err := c.dbRead.WithContext(ctx).Find(&categories).Error; err != nil {
		slg.Logger.Error("find categories", "error", err)
		return nil, customErr.ErrNotFound
	}

	return categories, nil
}

func (c *catalogRepository) FindCategoryById(ctx context.Context, id uint) (*domain.Category, error) {
	var category *domain.Category
	if err := c.dbRead.WithContext(ctx).First(&category, id).Error; err != nil {
		slg.Logger.Error("find category", "error", err)
		return nil, customErr.ErrNotFound
	}
	return category, nil
}

func (c *catalogRepository) EditCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := c.dbWrite.WithContext(ctx).Save(&category).Error; err != nil {
		slg.Logger.Error("edit category", "error", err)
		return nil, customErr.ErrUpdate
	}
	return category, nil
}

func (c *catalogRepository) DeleteCategory(ctx context.Context, id uint) error {
	if err := c.dbWrite.WithContext(ctx).Delete(&domain.Category{}, id).Error; err != nil {
		return customErr.ErrDelete
	}
	return nil
}

//func (c *catalogRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c catalogRepository) FindProducts(ctx context.Context) ([]*domain.Product, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c catalogRepository) FindProductById(ctx context.Context, id uint) (*domain.Product, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c catalogRepository) FindSellerProducts(ctx context.Context, id uint) ([]*domain.Product, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c catalogRepository) EditProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c catalogRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
//	//TODO implement me
//	panic("implement me")
//}

func NewCatalogRepository(dbWrite, dbRead *gorm.DB) CatalogRepository {
	return &catalogRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
