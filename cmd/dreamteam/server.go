package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/database"
	"github.com/sandeep-jaiswar/dreamteam/pkg/http"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/fx"
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

	app := fx.New(
		// Provide components
		fx.Provide(
			config.ProvideConfig,
			database.ProvideDatabase,
			http.ProvideHTTPServer,
		),
		// Register lifecycle hooks
		fx.Invoke(
			http.StartHTTPServer,
		),
	)

	// Graceful shutdown handling
	serverErrChan := make(chan error)

	go func() {
		app.Run()
		serverErrChan <- nil
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
