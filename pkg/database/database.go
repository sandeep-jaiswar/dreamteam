package database

import (
	"fmt"

	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProvideDatabase(cfg *config.Config) *gorm.DB {
	dbConfig := cfg.Database

	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		logInstance.DPanic("Error loading config: %v", zap.Error(err))
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logInstance.DPanic("Failed to connect to database: %v", zap.Error(err))
	}

	logInstance.Info("Database connected successfully!")
	return db
}
