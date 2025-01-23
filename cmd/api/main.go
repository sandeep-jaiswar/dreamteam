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
	"github.com/sandeep-jaiswar/dreamteam/pkg/middleware"
	// "github.com/sandeep-jaiswar/dreamteam/pkg/profiling"
	"github.com/sandeep-jaiswar/dreamteam/pkg/rbac"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	config.ProvideConfig()

	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer func() {
		if err := logInstance.Sync(); err != nil {
			panic("Failed to sync logger: %v\n" + err.Error())
		}
	}()

	app := fx.New(
		// Provide components
		fx.Provide(
			config.ProvideConfig,
			database.ProvideDatabase,
			http.ProvideHTTPServer,
		),
		// Register lifecycle hooks
		fx.Invoke(
			middleware.SecurityHeaders,
			rbac.InitializeRBAC,
			http.StartHTTPServer,
			//profiling.StartProfilingServer,
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
