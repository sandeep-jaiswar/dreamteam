package profiling

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/sandeep-jaiswar/dreamteam/pkg/config"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/zap"
)

func StartProfilingServer() {
	go func() {
		// Initialize logger
		logInstance, err := logger.GetLoggerInstance()
		port:= config.AppConfig.Profiling.Port
		if err != nil {
			logInstance.DPanic("Error loading config: %v", zap.Error(err))
		}
		logInstance.Info("Starting pprof server on port %d", zap.String("port", strconv.Itoa(port)))
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
			logInstance.Fatal("Failed to start pprof server", zap.Error(err))
		}
	}()
}
