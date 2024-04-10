package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ConfigInterface interface {
	Get(key string) string
}

type config struct {
	viper  *viper.Viper
	logger *zap.SugaredLogger
}

func NewConfig(logger *zap.SugaredLogger) (*config, error) {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &config{
		viper:  viper.GetViper(),
		logger: logger,
	}, nil
}

func (config *config) Get(key string) string {
	config.logger.Info("Getting config value", zap.String("key", key), zap.String("value", config.viper.GetString(key)))
	return config.viper.GetString(key)
}
