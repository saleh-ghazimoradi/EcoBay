package utils

import (
	"github.com/saleh-ghazimoradi/EcoBay/internal/service/service_models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(&service_models.User{})
}
