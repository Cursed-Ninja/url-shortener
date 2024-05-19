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

	d.logger.Infow("Sending request to database service", zap.String("Request Id", requestId), zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Errorw("Error creating request at database service", zap.String("Request Id", requestId), zap.Error(err))
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Errorw("Error sending request to database service", zap.String("Request Id", requestId), zap.Error(err))
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		d.logger.Errorw("Request failed at database service", zap.String("Request Id", requestId), zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return "", errors.New(http.StatusText(http.StatusNotFound))
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Errorw("Request failed at database service", zap.String("Request Id", requestId), zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return "", errors.New("request failed at database service")
	}

	d.logger.Infow("Request successful", zap.String("Request Id", requestId), zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Errorw("Empty response body from database service", zap.String("Request Id", requestId))
		return "", errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Errorw("Error reading response body at database service", zap.String("Request Id", requestId), zap.Error(err))
		return "", err
	}

	d.logger.Infow("Successfully read response body at database service", zap.String("Request Id", requestId), zap.Any("response", string(httpBody)))

	return string(httpBody), nil
}
