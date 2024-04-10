package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	mock_databaseservice "main-server/database-service/mocks"
	mock_config "main-server/internal/config/mocks"
	"main-server/internal/handlers"
	"main-server/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleShorten(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	mockDbService := mock_databaseservice.NewMockDatabaseServiceInterface(mockCtrl)
	mockConfig := mock_config.NewMockConfigInterface(mockCtrl)

	mockConfig.EXPECT().Get("BASE_URL").Return("http://localhost:8080").AnyTimes()

	handlers := handlers.NewBaseHandler(logger, mockDbService, mockConfig)

	tests := map[string]struct {
		reqBody                  *models.RequestModel
		HandleShorten            *gomock.Call
		HandleShortenReturnError error
		HandleShortenReturnUrl   *models.ShortenResponseModel
		HandleShortenCallTimes   int
		ExpectedStatusCode       int
	}{
		"EmptyBody": {
			reqBody:                  nil,
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: nil,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   0,
			ExpectedStatusCode:       http.StatusBadRequest,
		},
		"EmptyUrl": {
			reqBody: &models.RequestModel{
				ExpiresAt: time.Now(),
			},
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: nil,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   0,
			ExpectedStatusCode:       http.StatusBadRequest,
		},
		"EmptyExpiresAt": {
			reqBody: &models.RequestModel{
				Url: "http://localhost:8080",
			},
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: nil,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   1,
			ExpectedStatusCode:       http.StatusOK,
		},
		"InvalidUrl": {
			reqBody: &models.RequestModel{
				Url: "/test",
			},
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: nil,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   0,
			ExpectedStatusCode:       http.StatusBadRequest,
		},
		"DatabaseServiceFail": {
			reqBody: &models.RequestModel{
				Url: "http://localhost:8080",
			},
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: assert.AnError,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   1,
			ExpectedStatusCode:       http.StatusInternalServerError,
		},
		"Success": {
			reqBody: &models.RequestModel{
				Url:       "http://localhost:8080",
				ExpiresAt: time.Now().AddDate(1, 0, 0),
			},
			HandleShorten:            mockDbService.EXPECT().HandleShorten(gomock.Any(), gomock.Any()),
			HandleShortenReturnError: nil,
			HandleShortenReturnUrl:   &models.ShortenResponseModel{},
			HandleShortenCallTimes:   1,
			ExpectedStatusCode:       http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.HandleShorten.Return(test.HandleShortenReturnUrl, test.HandleShortenReturnError).Times(test.HandleShortenCallTimes)

			body, err := json.Marshal(test.reqBody)

			if err != nil {
				t.Error("Error marshalling body")
			}

			req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			handlers.HandleShorten(resp, req)
			assert.Equal(t, test.ExpectedStatusCode, resp.Code, resp.Result().Status)
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	mockDbService := mock_databaseservice.NewMockDatabaseServiceInterface(mockCtrl)
	mockConfig := mock_config.NewMockConfigInterface(mockCtrl)

	handlers := handlers.NewBaseHandler(logger, mockDbService, mockConfig)

	// handleRedirectDbService :=

	tests := map[string]struct {
		reqUrl                    string
		HandleRedirect            *gomock.Call
		HandleRedirectReturnError error
		HandleRedirectReturnUrl   *models.RedirectResponseModel
		HandleRedirectCallTimes   int
		ExpectedStatusCode        int
		muxVars                   map[string]string
	}{
		"EmptyUrl": {
			reqUrl:                    "/",
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: nil,
			HandleRedirectReturnUrl:   &models.RedirectResponseModel{},
			HandleRedirectCallTimes:   0,
			ExpectedStatusCode:        http.StatusBadRequest,
			muxVars: map[string]string{
				"url": "",
			},
		},
		"InvalidUrl": {
			reqUrl:                    "/absdn",
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: errors.New(http.StatusText(http.StatusNotFound)),
			HandleRedirectReturnUrl:   &models.RedirectResponseModel{},
			HandleRedirectCallTimes:   1,
			ExpectedStatusCode:        http.StatusNotFound,
			muxVars: map[string]string{
				"url": "absdn",
			},
		},
		"DatabaseServiceFail": {
			reqUrl:                    "/absdn",
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: assert.AnError,
			HandleRedirectReturnUrl:   &models.RedirectResponseModel{},
			HandleRedirectCallTimes:   1,
			ExpectedStatusCode:        http.StatusInternalServerError,
			muxVars: map[string]string{
				"url": "absdn",
			},
		},
		"Success": {
			reqUrl:                    "/adksjlkda",
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: nil,
			HandleRedirectReturnUrl: &models.RedirectResponseModel{
				Url: "https://google.com",
			},
			HandleRedirectCallTimes: 1,
			ExpectedStatusCode:      http.StatusMovedPermanently,
			muxVars: map[string]string{
				"url": "adksjlkda",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.HandleRedirect.Return(test.HandleRedirectReturnUrl, test.HandleRedirectReturnError).Times(test.HandleRedirectCallTimes)

			req := httptest.NewRequest("GET", test.reqUrl, nil)
			resp := httptest.NewRecorder()
			req = mux.SetURLVars(req, test.muxVars)
			handlers.HandleRedirect(resp, req)
			assert.Equal(t, test.ExpectedStatusCode, resp.Code, resp.Result().Status)
		})
	}
}
