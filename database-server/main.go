package main

import (
	"net/http"
	"url-shortner-database/internal/config"
	"url-shortner-database/internal/database"
	"url-shortner-database/internal/handlers"
	"url-shortner-database/internal/logging"
	"url-shortner-database/internal/middlewares"

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

	producer := kafka.NewKafkaProducer([]string{config.Get("KAFKA_SERVICE_BASE_URL")}, "database-server", true)

	logger, err = logging.NewLogger(producer)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	MONGO_URI := config.Get("MONGO_URI")
	DB_NAME := config.Get("DB_NAME")
	COLLECTION_NAME := config.Get("COLLECTION_NAME")

	mongoClient, err := database.NewDbConnection(logger, MONGO_URI, DB_NAME, COLLECTION_NAME)
	if err != nil {
		logger.Panic("Could not connect to database", zap.Error(err))
	}
	defer mongoClient.Disconnect()

	handlers := handlers.NewBaseHandler(logger, mongoClient)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handlers.HandleShorten).Methods(http.MethodPost)
	r.HandleFunc("/redirect", handlers.HandleRedirect).Methods(http.MethodPost)

	http.Handle("/", middlewares.LoggingMiddleware(r))
	logger.Error(http.ListenAndServe(":8081", nil))
}
