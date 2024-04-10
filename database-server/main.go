package main

import (
	"net/http"
	"url-shortner-database/internal/config"
	"url-shortner-database/internal/database"
	"url-shortner-database/internal/handlers"
	"url-shortner-database/internal/logging"

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
	r.HandleFunc("/{url}", handlers.HandleRedirect).Methods(http.MethodGet)

	http.Handle("/", r)
	logger.Error(http.ListenAndServe(":8081", nil))
}
