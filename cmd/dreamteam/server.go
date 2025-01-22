package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/database"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"github.com/sandeep-jaiswar/dreamteam/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logInstance.Sync()

	// Connect to the database
	database.Connect()

	// Initialize Gin
	r := gin.Default()

	// Define the server port
	port := config.AppConfig.App.Port

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Middleware setup for logging and security headers
	r.Use(middleware.SecurityHeaders())

	// Graceful shutdown handling
	serverErrChan := make(chan error)

	go func() {
		logInstance.Info("Starting Gin server on port " + strconv.Itoa(port))
		serverErrChan <- r.Run(":" + strconv.Itoa(port)) // Start the Gin server
	}()

	// Capture shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		logInstance.Info("Received shutdown signal")
	case err := <-serverErrChan:
		if err != nil {
			logInstance.Fatal("Server encountered an error", zap.Error(err))
		}
	}

	// Gracefully stop the Gin server
	logInstance.Info("Shutting down gracefully...")
}
