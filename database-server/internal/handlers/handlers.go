package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortner-database/internal/database"
	"url-shortner-database/internal/models"
	"url-shortner-database/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type baseHandler struct {
	dbConnection database.DBInterface
	logger       *zap.SugaredLogger
}

func NewBaseHandler(logger *zap.SugaredLogger, dbConnection database.DBInterface) *baseHandler {
	return &baseHandler{
		dbConnection: dbConnection,
		logger:       logger,
	}
}

func (h *baseHandler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("X-request-id")
	h.logger.Info(zap.String("Request Id", requestId), "Handling shorten request", zap.Any("request", r.Body))

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

	unmarsheledBody := &models.ShortenRequestModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
		return
	}

	if unmarsheledBody.Url == "" {
		h.logger.Error(zap.String("Request Id", requestId), "Empty URL in request body")
		http.Error(w, "Empty URL in request body", http.StatusBadRequest)
		return
	}

	url := models.URL{
		ShortUrlPath: utils.KeyGenerationService(unmarsheledBody.Url + requestId),
		OriginalUrl:  unmarsheledBody.Url,
		ExpiresAt:    utils.GetExpirationTime(unmarsheledBody.ExpiresAt),
	}

	for _, err := h.dbConnection.FindOne(bson.D{{Key: "shorturlpath", Value: url.ShortUrlPath}}); err == nil; {
		url.ShortUrlPath = utils.KeyGenerationService(unmarsheledBody.Url + requestId)
	}

	h.logger.Info(zap.String("Request Id", requestId), "Generated shortened URL", zap.Any("url", url))

	err = h.dbConnection.InsertOne(url)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error inserting document", zap.Error(err))
		http.Error(w, "Error inserting document", http.StatusInternalServerError)
		return
	}

	shortenedUrl := models.ShortenResponseModel{
		ShortUrlPath: url.ShortUrlPath,
	}

	jsonBody, err := json.Marshal(shortenedUrl)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

	h.logger.Info(zap.String("Request Id", requestId), "Successfully shortened URL", zap.Any("response", shortenedUrl))
}

func (h *baseHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
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

	url, err := h.dbConnection.FindOne(bson.D{{Key: "shorturlpath", Value: unmarsheledBody.ShortUrlPath}})

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Document not found", zap.Error(err))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	response := models.RedirectResponseModel{
		Url: url.OriginalUrl,
	}

	h.logger.Info(zap.String("Request Id", requestId), "Found document", zap.Any("document", url))

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error marshalling JSON", zap.Error(err))
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	h.logger.Info(zap.String("Request Id", requestId), "Successfully redirected URL", zap.Any("response", jsonResponse))
}
