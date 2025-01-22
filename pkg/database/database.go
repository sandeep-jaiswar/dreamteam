package database

import (
	"fmt"

	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}

func Connect() (*gorm.DB) {
    	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		logInstance.DPanic("Error loading config: %v", zap.Error(err))
	}
    dbConfig := config.AppConfig.Database

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)
    DB, err := ConnectDatabase(dsn)
    if err != nil {
		logInstance.DPanic("Failed to connect to database: %v", zap.Error(err))
	}

	logInstance.Info("Database connected successfully!")

    return DB
}
