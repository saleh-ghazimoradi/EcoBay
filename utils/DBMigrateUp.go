package utils

import (
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
	)
}
