package cache

import (
	"context"
	"fmt"

	"github.com/golang/groupcache"
	"github.com/sandeep-jaiswar/dreamteam/pkg/logger"
	"go.uber.org/zap"
)

var groupCache *groupcache.Group

// InitializeGroupCache initializes the groupcache instance
func InitializeGroupCache(name string, cacheSize int64, loader groupcache.GetterFunc) {
	// Initialize logger
	logInstance, err := logger.GetLoggerInstance()
	if err != nil {
		logInstance.DPanic("Error loading config: %v", zap.Error(err))
	}
	groupCache = groupcache.NewGroup(name, cacheSize, loader)
	logInstance.Info("Groupcache initialized with name: %s and size: %d", zap.String("name", name), zap.String("cacheSize", fmt.Sprintf("%d", cacheSize)))
}

// GetFromGroupCache retrieves a value from groupcache
func GetFromGroupCache(key string) ([]byte, error) {
	var data []byte
	err := groupCache.Get(context.TODO(), key, groupcache.AllocatingByteSliceSink(&data))
	return data, err
}
