package logging

import (
	"log"
	"os"

	"go.uber.org/zap"
)

func NewLogger() (*zap.SugaredLogger, error) {
	env := os.Getenv("APP_ENV")

	newLogger, err := zap.NewProduction()

	if env == "DEVELOPMENT" {
		newLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
		return nil, err
	}

	sugaredLogger := newLogger.Sugar()

	return sugaredLogger, nil
}
