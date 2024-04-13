package main

import (
	cacheservice "main-server/external/cache-service"
	databaseservice "main-server/external/database-service"
	"main-server/internal/config"
	"main-server/internal/handlers"
	"main-server/internal/logging"
	"main-server/internal/middlewares"
	"net/http"

	kafka "github.com/cursed-ninja/go-kafka-producer"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	producer, err := kafka.NewKafkaProducer([]string{"host1:9092", "host2:9092"}, "topic")

	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(producer)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatal("Could not load config", zap.Error(err))
	}

	databaseservice := databaseservice.NewDatabaseService(config, logger)
	cacheservice := cacheservice.NewCacheService(config, logger)

	handlers := handlers.NewBaseHandler(logger, databaseservice, config, cacheservice)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handlers.HandleShorten).Methods(http.MethodPost)
	r.HandleFunc("/{url}", handlers.HandleRedirect).Methods(http.MethodGet)

	http.Handle("/", middlewares.LoggingMiddleware(r))
	logger.Error(http.ListenAndServe(":8080", nil))
}
