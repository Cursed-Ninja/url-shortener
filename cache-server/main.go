package main

import (
	databaseservice "cache-server/external/database-service"
	"cache-server/internal/cache"
	"cache-server/internal/config"
	"cache-server/internal/handlers"
	"cache-server/internal/logging"
	"cache-server/internal/middlewares"
	"net/http"

	kafka "github.com/cursed-ninja/go-kafka-producer"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	initialLogger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	logger := initialLogger.Sugar()

	config, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatalw("Could not load config", zap.Error(err))
	}

	producer := kafka.NewKafkaProducer([]string{config.Get("KAFKA_SERVICE_BASE_URL")}, "cache-server", true)

	logger, err = logging.NewLogger(producer)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	dbService := databaseservice.NewDatabaseService(config, logger)

	cacheService, err := cache.NewCache(config, logger)
	if err != nil {
		logger.Fatalw("Could not create cache service", zap.Error(err))
	}
	defer cacheService.Close()

	handler := handlers.NewHandler(cacheService, logger, config, dbService)

	r := mux.NewRouter()
	r.HandleFunc("/redirect", handler.HandleRedirect).Methods(http.MethodPost)

	http.Handle("/", middlewares.LoggingMiddleware(r))
	logger.Error(http.ListenAndServe(":8082", nil))
}
