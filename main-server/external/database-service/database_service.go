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

	d.logger.Infow("Sending request to database service", zap.String("Request Id", requestId), zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Errorw("Error creating request at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Errorw("Error sending request to database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Errorw("Request failed at database service", zap.String("Request Id", requestId), zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New("request failed at database service")
	}

	d.logger.Infow("Request successful", zap.String("Request Id", requestId), zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Errorw("Empty response body from database service", zap.String("Request Id", requestId))
		return nil, errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Errorw("Error reading response body at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	d.logger.Infow("Successfully read response body at database service", zap.String("Request Id", requestId), zap.Any("response", httpBody))

	unmarsheledBody := &models.ShortenResponseModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		d.logger.Errorw("Error unmarshalling response body at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	d.logger.Infow("Successfully unmarshalled response body at database service", zap.String("Request Id", requestId), zap.Any("request", unmarsheledBody))

	return unmarsheledBody, nil
}

func (d *databaseService) HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error) {
	reqUrl := d.config.Get("DATABASE_SERVICE_BASE_URL") + "/redirect"

	d.logger.Infow("Sending request to database service", zap.String("Request Id", requestId), zap.String("url", reqUrl))

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)

	if err != nil {
		d.logger.Errorw("Error creating request at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-request-id", requestId)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		d.logger.Errorw("Error sending request to database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		d.logger.Errorw("Request failed at database service", zap.String("Request Id", requestId), zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New(http.StatusText(http.StatusNotFound))
	}

	if resp.StatusCode != http.StatusOK {
		d.logger.Errorw("Request failed at database service", zap.String("Request Id", requestId), zap.Int("status", resp.StatusCode), zap.String("status", resp.Status))
		return nil, errors.New("request failed at database service")
	}

	d.logger.Infow("Request successful", zap.String("Request Id", requestId), zap.String("status", resp.Status))

	if resp.Body == nil {
		d.logger.Errorw("Empty response body from database service", zap.String("Request Id", requestId))
		return nil, errors.New("empty response body")
	}

	httpBody, err := io.ReadAll(resp.Body)

	if err != nil {
		d.logger.Errorw("Error reading response body at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	d.logger.Infow("Successfully read response body at database service", zap.String("Request Id", requestId), zap.Any("response", httpBody))

	unmarsheledBody := &models.RedirectResponseModel{}

	err = json.Unmarshal(httpBody, unmarsheledBody)

	if err != nil {
		d.logger.Errorw("Error unmarshalling response body at database service", zap.String("Request Id", requestId), zap.Error(err))
		return nil, err
	}

	d.logger.Infow("Successfully unmarshalled response body at database service", zap.String("Request Id", requestId), zap.Any("request", unmarsheledBody))

	return unmarsheledBody, nil
}
