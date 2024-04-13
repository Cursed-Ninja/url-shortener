package cacheservice

import (
	"encoding/json"
	"errors"
	"io"
	"main-server/internal/config"
	"main-server/internal/models"
	"net/http"

	"go.uber.org/zap"
)

type CacheServiceInterface interface {
	HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error)
}

type cacheService struct {
	config config.ConfigInterface
	logger *zap.SugaredLogger
}

func NewCacheService(config config.ConfigInterface, logger *zap.SugaredLogger) *cacheService {
	return &cacheService{
		config: config,
		logger: logger,
	}
}

func (c *cacheService) HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error) {
	reqUrl := c.config.Get("CACHE_SERVICE_BASE_URL") + "/redirect"

	c.logger.Info(zap.String("requestId", requestId), "Sending request to cache service", zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		c.logger.Error(zap.String("requestId", requestId), "Error creating request at cache service", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.logger.Error(zap.String("requestId", requestId), "Error sending request to cache service", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusNotFound {
		c.logger.Error(zap.String("requestId", requestId), "Request failed at cache service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	c.logger.Info(zap.String("requestId", requestId), "Request successful", zap.String("status", resp.Status))

	if resp.Body == nil {
		c.logger.Error(zap.String("Request Id", requestId), "Empty response body from cache service")
		return nil, errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		c.logger.Error(zap.String("Request Id", requestId), "Error reading response body at cache service", zap.Error(err))
		return nil, err
	}

	c.logger.Info(zap.String("Request Id", requestId), "Successfully read response body at cache service", zap.Any("response", httpBody))

	unmarsheledBody := &models.RedirectResponseModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		c.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling response body at database service", zap.Error(err))
		return nil, err
	}

	c.logger.Info(zap.String("Request Id", requestId), "Successfully unmarshalled response body at database service", zap.Any("request", unmarsheledBody))

	return unmarsheledBody, nil
}
