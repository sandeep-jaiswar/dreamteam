package http

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideHTTPServer(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	return r
}

// StartHTTPServer registers lifecycle hooks to start and stop the HTTP server.
func StartHTTPServer(lifecycle fx.Lifecycle, r *gin.Engine, cfg *config.Config) {
	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		logInstance.DPanic("Error loading config: %v", zap.Error(err))
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				port := cfg.App.Port
				if port <= 0 {
					logInstance.DPanic("Invalid port number: ", zap.String("port", strconv.Itoa(port)))
				}
				address := fmt.Sprintf(":%d", port)
				fmt.Printf("Starting server on %s\n", address)
				if err := r.Run(address); err != nil {
					logInstance.DPanic("Failed to start server: %v", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return nil
		},
	})
}
