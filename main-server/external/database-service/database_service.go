package databaseservice

import (
	"encoding/json"
	"errors"
	"io"
	"main-server/internal/config"
	"main-server/internal/models"
	"net/http"

	"go.uber.org/zap"
)

type DatabaseServiceInterface interface {
	HandleShorten(body io.Reader, requestId string) (*models.ShortenResponseModel, error)
	HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error)
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

func (d *databaseService) HandleShorten(body io.Reader, requestId string) (*models.ShortenResponseModel, error) {
	reqUrl := d.config.Get("DATABASE_SERVICE_BASE_URL") + "/shorten"

	d.logger.Info(zap.String("requestId", requestId), "Sending request to database service", zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error creating request at database service", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error sending request to database service", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Error(zap.String("requestId", requestId), "Request failed at database service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New("request failed at database service")
	}

	d.logger.Info(zap.String("requestId", requestId), "Request successful", zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Error(zap.String("Request Id", requestId), "Empty response body from database service")
		return nil, errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Error(zap.String("Request Id", requestId), "Error reading response body at database service", zap.Error(err))
		return nil, err
	}

	d.logger.Info(zap.String("Request Id", requestId), "Successfully read response body at database service", zap.Any("response", httpBody))

	unmarsheledBody := &models.ShortenResponseModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		d.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling response body at database service", zap.Error(err))
		return nil, err
	}

	d.logger.Info(zap.String("Request Id", requestId), "Successfully unmarshalled response body at database service", zap.Any("request", unmarsheledBody))

	return unmarsheledBody, nil
}

func (d *databaseService) HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error) {
	reqUrl := d.config.Get("DATABASE_SERVICE_BASE_URL") + "/redirect"

	d.logger.Info(zap.String("requestId", requestId), "Sending request to database service", zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error creating request at database service", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Error(zap.String("requestId", requestId), "Error sending request to database service", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		d.logger.Error(zap.String("requestId", requestId), "Request failed at database service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Error(zap.String("requestId", requestId), "Request failed at database service", zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New("request failed at database service")
	}

	d.logger.Info(zap.String("requestId", requestId), "Request successful", zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Error(zap.String("Request Id", requestId), "Empty response body from database service")
		return nil, errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Error(zap.String("Request Id", requestId), "Error reading response body at database service", zap.Error(err))
		return nil, err
	}

	d.logger.Info(zap.String("Request Id", requestId), "Successfully read response body at database service", zap.Any("response", httpBody))

	unmarsheledBody := &models.RedirectResponseModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		d.logger.Error(zap.String("Request Id", requestId), "Error unmarshalling response body at database service", zap.Error(err))
		return nil, err
	}

	d.logger.Info(zap.String("Request Id", requestId), "Successfully unmarshalled response body at database service", zap.Any("request", unmarsheledBody))

	return unmarsheledBody, nil
}
