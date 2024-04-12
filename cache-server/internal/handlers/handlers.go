package handlers

import (
	"bytes"
	databaseservice "cache-server/external/database-service"
	"cache-server/internal/cache"
	"cache-server/internal/config"
	"cache-server/internal/models"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HandlerInterface interface {
	HandleRedirect(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	cache     cache.CacheInterface
	logger    *zap.SugaredLogger
	config    config.ConfigInterface
	dbService databaseservice.DatabaseServiceInterface
}

func NewHandler(cache cache.CacheInterface, logger *zap.SugaredLogger, config config.ConfigInterface, dbService databaseservice.DatabaseServiceInterface) *handler {
	return &handler{
		cache:     cache,
		logger:    logger,
		config:    config,
		dbService: dbService,
	}
}

func (h *handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("X-request-id")

	h.logger.Info(zap.String("Request Id", requestId), "Handling redirect request", zap.Any("request", r.Body))

	if r.Body == nil {
		h.logger.Error(zap.String("Request Id", requestId), "Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	httpBody, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error reading request body", zap.Error(err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Successfully read request body", zap.Any("request", httpBody))

	unmarsheledBody := &models.RedirectRequestModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		return
	}

	if unmarsheledBody.ShortUrlPath == "" {
		h.logger.Error(zap.String("Request Id", requestId), "Empty URL in request body")
		http.Error(w, "Empty URL in request body", http.StatusBadRequest)
		return
	}

	val, err := h.cache.GetValue(string(httpBody), requestId)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error retrieving value from cache", zap.Error(err))
		val, err = h.dbService.HandleRedirect(bytes.NewBuffer(httpBody), requestId)

		if err != nil {
			if err.Error() == http.StatusText(http.StatusNotFound) {
				h.logger.Error(zap.String("Request Id", requestId), "URL not found", zap.Error(err))
				http.Error(w, "URL not found", http.StatusNotFound)
				err = h.cache.SetValue(string(httpBody), "", requestId, 10*time.Second)

				if err != nil {
					h.logger.Error(zap.String("Request Id", requestId), "Error setting value in cache", zap.Error(err))
				}

				return
			}

			h.logger.Error(zap.String("Request Id", requestId), "Error processing redirect request", zap.Error(err))
			http.Error(w, "Error processing redirect request", http.StatusInternalServerError)
			return
		}

		err = h.cache.SetValue(string(httpBody), val, requestId, 3*time.Minute)

		if err != nil {
			h.logger.Error(zap.String("Request Id", requestId), "Error setting value in cache", zap.Error(err))
		}

		h.logger.Info(zap.String("Request Id", requestId), "Successfully processed redirect request", zap.String("value", val))
	} else {
		h.logger.Info(zap.String("Request Id", requestId), "Successfully retrieved value from cache", zap.String("value", val))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))

	h.logger.Info(zap.String("Request Id", requestId), "Successfully responded to redirect request")

}
