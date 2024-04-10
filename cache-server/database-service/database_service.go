package databaseservice

import (
	"cache-server/internal/config"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type DatabaseServiceInterface interface {
	HandleRedirect(body io.Reader, requestId string) (string, error)
}

type databaseService struct {
	config config.ConfigInterface
	logger *zap.SugaredLogger
}

func NewDatabaseService(config config.ConfigInterface, logger *zap.SugaredLogger) *databaseService {
	return &databaseService{
		config: config,
		logger: logger,
	}
}

func (d *databaseService) HandleRedirect(body io.Reader, requestId string) (string, error) {
	reqUrl := d.config.Get("DATABASE_SERVICE_BASE_URL") + "/redirect"

	d.logger.Info(zap.String("requestId", requestId), "Sending request to database service", zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error creating request at database service", zap.Error(err))
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error sending request to database service", zap.Error(err))
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		d.logger.Error(zap.String("requestId", requestId), "Request failed at database service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return "", errors.New(http.StatusText(http.StatusNotFound))
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Error(zap.String("requestId", requestId), "Request failed at database service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return "", errors.New("request failed at database service")
	}

	d.logger.Info(zap.String("requestId", requestId), "Request successful", zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Error(zap.String("Request Id", requestId), "Empty response body from database service")
		return "", errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Error(zap.String("Request Id", requestId), "Error reading response body at database service", zap.Error(err))
		return "", err
	}
	

	d.logger.Info(zap.String("Request Id", requestId), "Successfully read response body at database service", zap.Any("response", string(httpBody)))

	return string(httpBody), nil
}
