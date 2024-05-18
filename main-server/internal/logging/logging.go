package logging

import (
	"io"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(kafkaProducer io.Writer) (*zap.SugaredLogger, error) {
	env := os.Getenv("APP_ENV")

	newLogger, err := zap.NewDevelopment()

	if env == "production" {
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		var kafkaPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel
		})

		newLogger = zap.New(zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(kafkaProducer)), kafkaPriority))
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
		return nil, err
	}

	sugaredLogger := newLogger.Sugar()

	return sugaredLogger, nil
}
