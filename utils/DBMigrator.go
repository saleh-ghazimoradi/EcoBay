package utils

import "gorm.io/gorm"

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate()
}
