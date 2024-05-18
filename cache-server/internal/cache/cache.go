package cache

import (
	"cache-server/internal/config"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheInterface interface {
	GetValue(key, requestId string) (string, error)
	SetValue(key, value, requestId string, expiryTime time.Duration) error
}

type cache struct {
	client *redis.Client
	logger *zap.SugaredLogger
}

func NewCache(config config.ConfigInterface, logger *zap.SugaredLogger) (*cache, error) {
	REDIS_ADDR := config.Get("REDIS_ADDR")
	REDIS_PASSWORD := config.Get("REDIS_PASSWORD")
	REDIS_DB, err := strconv.Atoi(config.Get("REDIS_DB"))

	if err != nil {
		logger.Errorw("Error converting RedisDb to int", zap.Error(err))
		return nil, err
	}

	cache := &cache{
		client: redis.NewClient(&redis.Options{
			Addr:     REDIS_ADDR,
			Password: REDIS_PASSWORD,
			DB:       REDIS_DB,
		}),
		logger: logger,
	}

	_, err = cache.client.Ping(context.Background()).Result()

	if err != nil {
		logger.Errorw("Error connecting to Redis", zap.Error(err))
		defer cache.client.Close()
		return nil, err
	}

	logger.Infow("Successfully connected to Redis")
	return cache, nil
}

func (cache *cache) GetValue(key, requestId string) (string, error) {
	cache.logger.Infow("Retrieve from cache", zap.String("Request Id", requestId), zap.String("key", key))

	val, err := cache.client.Get(context.Background(), key).Result()

	if err != nil {
		if err == redis.Nil {
			cache.logger.Infow("Key not found in cache", zap.String("Request Id", requestId), zap.String("key", key))
		} else {
			cache.logger.Errorw("Error retrieving key", zap.String("Request Id", requestId), zap.Error(err))
		}

		return "", err
	}

	cache.logger.Infow("Successfully retrieved value", zap.String("Request Id", requestId), zap.String("key", key), zap.String("value", val))

	return val, nil
}

func (cache *cache) SetValue(key, value, requestId string, expiryTime time.Duration) error {
	cache.logger.Infow("Set value in cache", zap.String("Request Id", requestId), zap.String("key", key), zap.String("value", value), zap.Any("expiryTime", expiryTime))

	err := cache.client.Set(context.Background(), key, value, expiryTime).Err()

	if err != nil {
		cache.logger.Errorw("Error setting value", zap.String("Request Id", requestId), zap.Error(err))
	}

	return err
}

func (cache *cache) Flush() error {
	return cache.client.FlushAll(context.Background()).Err()
}

func (cache *cache) Close() error {
	return cache.client.Close()
}
