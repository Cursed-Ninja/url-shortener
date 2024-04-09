package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"url-shortner-database/database"
	"url-shortner-database/models"
	"url-shortner-database/utils"

	"github.com/gorilla/mux"
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
	h.logger.Info("Handling shorten request", zap.Any("request", r.Body))

	httpBody, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Error("Error reading request body", zap.Error(err))
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	unmarsheledBody := &models.RequestModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		h.logger.Error("Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	shortenedUrl := models.RequestModel{
		Url: utils.KeyGenerationService(unmarsheledBody.Url),
	}

	jsonBody, err := json.Marshal(shortenedUrl)

	if err != nil {
		h.logger.Error("Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	url := models.URL{
		ShortenedUrl: shortenedUrl.Url,
		Url:          unmarsheledBody.Url,
		ExpiresAt:    utils.GetExpirationTime(),
	}

	err = h.dbConnection.InsertOne(url)

	if err != nil {
		h.logger.Error("Error inserting document", zap.Error(err))
		http.Error(w, "Error inserting document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

	h.logger.Info("Successfully shortened URL", zap.Any("response", jsonBody))
}

func (h *baseHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	if _, ok := vars["url"]; !ok {
		h.logger.Error("URL variable not found")
		http.Error(w, "URL variable not found", http.StatusBadRequest)
		return
	}

	h.logger.Info("Handling redirect request", zap.String("request", vars["url"]))

	url, err := h.dbConnection.FindOne(bson.D{{Key: "shortenedurl", Value: vars["url"]}})

	if err != nil {
		h.logger.Info("Document not found", zap.Error(err))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if url.Url == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	response := models.RequestModel{
		Url: url.Url,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		h.logger.Error("Error marshalling JSON", zap.Error(err))
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	h.logger.Info("Successfully redirected URL", zap.Any("response", jsonResponse))
}
