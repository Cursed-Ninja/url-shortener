package main

import (
	databaseservice "main-server/database-service"
	"main-server/internal/config"
	"main-server/internal/handlers"
	"main-server/internal/logging"
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

	databaseservice := databaseservice.NewDatabaseService(config, logger)

	handlers := handlers.NewBaseHandler(logger, databaseservice, config)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handlers.HandleShorten).Methods(http.MethodPost)
	r.HandleFunc("/{url}", handlers.HandleRedirect).Methods(http.MethodGet)

	http.Handle("/", r)
	logger.Error(http.ListenAndServe(":8080", nil))
}
