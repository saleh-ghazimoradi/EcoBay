package repository

import "gorm.io/gorm"

type CatalogRepository interface {
}

type catalogRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func NewCatalogRepository(dbWrite, dbRead *gorm.DB) CatalogRepository {
	return &catalogRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
