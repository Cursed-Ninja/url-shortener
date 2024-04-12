// Make sure redis server is running
package cache_test

import (
	"cache-server/internal/cache"
	mock_config "cache-server/internal/config/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewCache(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	config := mock_config.NewMockConfigInterface(mockCtrl)

	t.Run("Success", func(t *testing.T) {
		config.EXPECT().Get("REDIS_ADDR").Return("localhost:6379")
		config.EXPECT().Get("REDIS_PASSWORD").Return("")
		config.EXPECT().Get("REDIS_DB").Return("0")

		_, err := cache.NewCache(config, logger)

		assert.Nil(t, err, "Error connecting to Redis")
	})
}

func TestSetValue(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	config := mock_config.NewMockConfigInterface(mockCtrl)

	t.Run("Success", func(t *testing.T) {
		config.EXPECT().Get("REDIS_ADDR").Return("localhost:6379")
		config.EXPECT().Get("REDIS_PASSWORD").Return("")
		config.EXPECT().Get("REDIS_DB").Return("0")

		cache, _ := cache.NewCache(config, logger)

		err := cache.SetValue("key", "value", "requestId", 10)

		assert.Nil(t, err, "Error setting value in cache")

		err = cache.Flush()

		assert.Nil(t, err, "Error flushing cache")
	})
}

func TestGetValue(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	config := mock_config.NewMockConfigInterface(mockCtrl)

	t.Run("Key does not exist", func(t *testing.T) {
		config.EXPECT().Get("REDIS_ADDR").Return("localhost:6379")
		config.EXPECT().Get("REDIS_PASSWORD").Return("")
		config.EXPECT().Get("REDIS_DB").Return("0")

		cache, _ := cache.NewCache(config, logger)

		_, err := cache.GetValue("key", "requestId")

		assert.NotNil(t, err, "Error getting value from cache")
	})

	t.Run("Key exists", func(t *testing.T) {
		config.EXPECT().Get("REDIS_ADDR").Return("localhost:6379")
		config.EXPECT().Get("REDIS_PASSWORD").Return("")
		config.EXPECT().Get("REDIS_DB").Return("0")

		cache, _ := cache.NewCache(config, logger)

		err := cache.SetValue("key", "value", "requestId", time.Second*10)

		assert.Nil(t, err, "Error setting value in cache")

		val, err := cache.GetValue("key", "requestId")

		assert.Nil(t, err, "Error getting value from cache")
		assert.Equal(t, "value", val)

		err = cache.Flush()

		assert.Nil(t, err, "Error flushing cache")
	})
}
