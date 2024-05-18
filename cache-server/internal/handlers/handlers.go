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

	h.logger.Infow("Handling redirect request", zap.String("Request Id", requestId), zap.Any("request", r.Body))

	if r.Body == nil {
		h.logger.Errorw("Empty request body", zap.String("Request Id", requestId))
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	httpBody, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Errorw("Error reading request body", zap.String("Request Id", requestId), zap.Error(err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("Successfully read request body", zap.String("Request Id", requestId), zap.Any("request", httpBody))

	unmarsheledBody := &models.RedirectRequestModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		h.logger.Errorw("Error unmarshalling request body", zap.String("Request Id", requestId), zap.Error(err))
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		return
	}

	if unmarsheledBody.ShortUrlPath == "" {
		h.logger.Errorw("Empty URL in request body", zap.String("Request Id", requestId))
		http.Error(w, "Empty URL in request body", http.StatusBadRequest)
		return
	}

	val, err := h.cache.GetValue(string(httpBody), requestId)

	if err != nil {
		h.logger.Errorw("Error retrieving value from cache", zap.String("Request Id", requestId), zap.Error(err))
		val, err = h.dbService.HandleRedirect(bytes.NewBuffer(httpBody), requestId)

		if err != nil {
			if err.Error() == http.StatusText(http.StatusNotFound) {
				h.logger.Errorw("URL not found", zap.String("Request Id", requestId), zap.Error(err))
				http.Error(w, "URL not found", http.StatusNotFound)
				err = h.cache.SetValue(string(httpBody), "", requestId, 10*time.Second)

				if err != nil {
					h.logger.Errorw("Error setting value in cache", zap.String("Request Id", requestId), zap.Error(err))
				}

				return
			}

			h.logger.Errorw("Error processing redirect request", zap.String("Request Id", requestId), zap.Error(err))
			http.Error(w, "Error processing redirect request", http.StatusInternalServerError)
			return
		}

		err = h.cache.SetValue(string(httpBody), val, requestId, 3*time.Minute)

		if err != nil {
			h.logger.Errorw("Error setting value in cache", zap.String("Request Id", requestId), zap.Error(err))
		}

		h.logger.Infow("Successfully processed redirect request", zap.String("Request Id", requestId), zap.String("value", val))
	} else {
		h.logger.Infow("Successfully retrieved value from cache", zap.String("Request Id", requestId), zap.String("value", val))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))

	h.logger.Infow("Successfully responded to redirect request", zap.String("Request Id", requestId))

}
