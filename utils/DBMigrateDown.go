package utils

import (
	"github.com/saleh-ghazimoradi/EcoBay/internal/domain"
	"gorm.io/gorm"
)

func DBMigrateDrop(db *gorm.DB) error {
	err := db.Migrator().DropTable(
		&domain.User{},
		&domain.Address{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
	)
	if err != nil {
		return err
	}
	return nil
}
