package main

import (
	databaseservice "cache-server/external/database-service"
	"cache-server/internal/cache"
	"cache-server/internal/config"
	"cache-server/internal/handlers"
	"cache-server/internal/logging"
	"cache-server/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatal("Could not load config", zap.Error(err))
	}

	dbService := databaseservice.NewDatabaseService(config, logger)

	cacheService, err := cache.NewCache(config, logger)
	if err != nil {
		logger.Fatal("Could not create cache service", zap.Error(err))
	}
	defer cacheService.Close()

	handler := handlers.NewHandler(cacheService, logger, config, dbService)

	r := mux.NewRouter()
	r.HandleFunc("/redirect", handler.HandleRedirect).Methods(http.MethodPost)

	http.Handle("/", middlewares.LoggingMiddleware(r))
	logger.Error(http.ListenAndServe(":8082", nil))
}
