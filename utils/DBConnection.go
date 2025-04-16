package utils

import (
	"fmt"
	"github.com/saleh-ghazimoradi/EcoBay/config"
	"github.com/saleh-ghazimoradi/EcoBay/slg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func postURI() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.Database.DatabaseHost, config.AppConfig.Database.DatabasePort, config.AppConfig.Database.DatabaseUser, config.AppConfig.Database.DatabasePassword, config.AppConfig.Database.DatabaseName, config.AppConfig.Database.DatabaseSSLMode)
}

func DBConnection(DBMigrator func(db *gorm.DB) error) (*gorm.DB, error) {
	uri := postURI()
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		slg.Logger.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	postDB, err := db.DB()
	if err != nil {
		slg.Logger.Error("Unable to get sql.DB", "error", err)
		return nil, err
	}

	postDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpenConn)
	postDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdleConn)
	postDB.SetConnMaxLifetime(config.AppConfig.Database.MaxLifetime)
	postDB.SetConnMaxIdleTime(config.AppConfig.Database.MaxIdleTime)

	slg.Logger.Info("Successfully connected to database")

	if err = DBMigrator(db); err != nil {
		slg.Logger.Error("Unable to migrate database", "error", err)
	}

	return db, nil
}
