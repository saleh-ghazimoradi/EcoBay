package utils

import (
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func dbURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.DbConfig.DbHost, config.AppConfig.DbConfig.DbPort, config.AppConfig.DbConfig.DbUser, config.AppConfig.DbConfig.DbPassword, config.AppConfig.DbConfig.DbName, config.AppConfig.DbConfig.DbSslMode)
}

func DBConnection(DBMigrator func(db *gorm.DB) error) (*gorm.DB, error) {
	uri := dbURI()
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	if err = DBMigrator(db); err != nil {
		log.Fatalf("Unable to migrate database: %v", err)
	}

	return db, nil
}
