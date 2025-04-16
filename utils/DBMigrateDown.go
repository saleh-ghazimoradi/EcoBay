package utils

import "gorm.io/gorm"

func DBMigrateDrop(db *gorm.DB) error {
	err := db.Migrator().DropTable()
	if err != nil {
		return err
	}
	return nil
}
