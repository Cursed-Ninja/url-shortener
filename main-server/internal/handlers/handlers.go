package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	cacheservice "main-server/external/cache-service"
	databaseservice "main-server/external/database-service"
	"main-server/internal/config"
	"main-server/internal/models"
	"net/http"

	UrlVerifier "github.com/davidmytton/url-verifier"
	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type HandlerInterface interface {
	HandleShorten(w http.ResponseWriter, r *http.Request)
	HandleRedirect(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	logger          *zap.SugaredLogger
	databaseservice databaseservice.DatabaseServiceInterface
	config          config.ConfigInterface
	cacheservice    cacheservice.CacheServiceInterface
}

func NewBaseHandler(logger *zap.SugaredLogger, databaseservice databaseservice.DatabaseServiceInterface, config config.ConfigInterface, cacheservice cacheservice.CacheServiceInterface) *handler {
	return &handler{
		logger:          logger,
		databaseservice: databaseservice,
		config:          config,
		cacheservice:    cacheservice,
	}
}

func (h *handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("X-request-id")

	h.logger.Info(zap.String("Request Id", requestId), "Call to HandleShorten", zap.Any("Request", r.Body))

	if r.Body == nil {
		h.logger.Error(zap.String("Request Id", requestId), "Empty request body")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	httpBody, err := io.ReadAll(r.Body)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error reading request body", zap.Error(err))
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Successfully read request body", zap.Any("request", httpBody))

	unmarsheledBody := &models.RequestModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil || unmarsheledBody.Url == "" {
		h.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling request body", zap.Error(err))
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Successfully unmarshalled request body", zap.Any("request", unmarsheledBody))

	urlVerifier := UrlVerifier.NewVerifier()
	result, err := urlVerifier.Verify(unmarsheledBody.Url)

	if err != nil || result == nil || !result.IsURL {
		h.logger.Error(zap.String("Request Id", requestId), "Invalid URL", zap.Error(err))
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "URL is valid", zap.Any("url", unmarsheledBody.Url))

	shortenRequestModel := &models.ShortenRequestModel{
		Url:       unmarsheledBody.Url,
		ExpiresAt: unmarsheledBody.ExpiresAt,
	}

	h.logger.Info(zap.String("Request Id", requestId), "Shorten Request Model", zap.Any("model", shortenRequestModel))

	shortenRequestModelJson, err := json.Marshal(shortenRequestModel)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error marshalling shorten request model", zap.Error(err))
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Successfully marshalled shorten request model", zap.Any("model", shortenRequestModel))

	shortenResponseModel, err := h.databaseservice.HandleShorten(bytes.NewBuffer(shortenRequestModelJson), requestId)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error processing shorten request", zap.Error(err))
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	responseModel := &models.ResponseModel{
		Url: h.config.Get("BASE_URL") + "/" + shortenResponseModel.ShortUrlPath,
	}

	h.logger.Info(zap.String("Request Id", requestId), "Response Model", zap.Any("model", responseModel))

	jsonBody, err := json.Marshal(responseModel)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error marshalling shorten response model", zap.Error(err))
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

	h.logger.Info(zap.String("Request Id", requestId), "Successfully handled shorten request")
}

func (h *handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get("X-request-id")

	vars := mux.Vars(r)

	if url, ok := vars["url"]; !ok || url == "" {
		h.logger.Info(zap.String("Request Id", requestId), "URL variable not found")
		http.Error(w, "URL variable not found", http.StatusBadRequest)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Handling redirect request", zap.String("request", vars["url"]))

	redirectRequestModel := &models.RedirectRequestModel{
		ShortUrlPath: vars["url"],
	}

	h.logger.Info(zap.String("Request Id", requestId), "Redirect Request Model", zap.Any("model", redirectRequestModel))

	redirectRequestModelJson, err := json.Marshal(redirectRequestModel)

	if err != nil {
		h.logger.Error(zap.String("Request Id", requestId), "Error marshalling redirect request model", zap.Error(err))
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}

	h.logger.Info(zap.String("Request Id", requestId), "Successfully marshalled redirect request model", zap.Any("model", redirectRequestModelJson))

	redirectResponseModel, err := h.cacheservice.HandleRedirect(bytes.NewBuffer(redirectRequestModelJson), requestId)

	if err != nil {
		if err.Error() == http.StatusText(http.StatusNotFound) {
			h.logger.Error(zap.String("Request Id", requestId), "URL not found", zap.Error(err))
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		redirectResponseModel, err = h.databaseservice.HandleRedirect(bytes.NewBuffer(redirectRequestModelJson), requestId)
		if err != nil {
			if err.Error() == http.StatusText(http.StatusNotFound) {
				h.logger.Error(zap.String("Request Id", requestId), "URL not found", zap.Error(err))
				http.Error(w, "URL not found", http.StatusNotFound)
				return
			}

			h.logger.Error(zap.String("Request Id", requestId), "Error processing redirect request", zap.Error(err))
			http.Error(w, "Something went wrong!", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, redirectResponseModel.Url, http.StatusMovedPermanently)
	h.logger.Info(zap.String("Request Id", requestId), "Successfully handled redirect request")
}
