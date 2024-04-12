package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_databaseservice "cache-server/external/database-service/mocks"
	mock_cache "cache-server/internal/cache/mocks"
	mock_config "cache-server/internal/config/mocks"
	"cache-server/internal/handlers"
	"cache-server/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleRedirect(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)

	mockCache := mock_cache.NewMockCacheInterface(mockCtrl)
	mockConfig := mock_config.NewMockConfigInterface(mockCtrl)
	mockDbService := mock_databaseservice.NewMockDatabaseServiceInterface(mockCtrl)

	handler := handlers.NewHandler(mockCache, logger.Sugar(), mockConfig, mockDbService)

	tests := map[string]struct {
		requestBody               *models.RedirectRequestModel
		GetValue                  *gomock.Call
		GetValueReturnError       error
		GetValueReturnVal         string
		GetValueCallTimes         int
		HandleRedirect            *gomock.Call
		HandleRedirectReturnError error
		HandleRedirectReturnVal   string
		HandleRedirectCallTimes   int
		SetValue                  *gomock.Call
		SetValueReturnError       error
		SetValueCallTimes         int
		ExpectedStatusCode        int
	}{
		"EmptyBody": {
			requestBody:               nil,
			GetValue:                  mockCache.EXPECT().GetValue(gomock.Any(), gomock.Any()),
			GetValueReturnError:       nil,
			GetValueReturnVal:         "",
			GetValueCallTimes:         0,
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: nil,
			HandleRedirectReturnVal:   "",
			HandleRedirectCallTimes:   0,
			SetValue:                  mockCache.EXPECT().SetValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			SetValueReturnError:       nil,
			SetValueCallTimes:         0,
			ExpectedStatusCode:        http.StatusBadRequest,
		},
		"Cache Miss With DB Error": {
			requestBody:               &models.RedirectRequestModel{ShortUrlPath: "shortUrl"},
			GetValue:                  mockCache.EXPECT().GetValue(gomock.Any(), gomock.Any()),
			GetValueReturnError:       assert.AnError,
			GetValueReturnVal:         "",
			GetValueCallTimes:         1,
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: assert.AnError,
			HandleRedirectReturnVal:   "",
			HandleRedirectCallTimes:   1,
			SetValue:                  mockCache.EXPECT().SetValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			SetValueReturnError:       nil,
			SetValueCallTimes:         0,
			ExpectedStatusCode:        http.StatusInternalServerError,
		},
		"Cache Miss With DB Not Found": {
			requestBody:               &models.RedirectRequestModel{ShortUrlPath: "shortUrl"},
			GetValue:                  mockCache.EXPECT().GetValue(gomock.Any(), gomock.Any()),
			GetValueReturnError:       assert.AnError,
			GetValueReturnVal:         "",
			GetValueCallTimes:         1,
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: errors.New(http.StatusText(http.StatusNotFound)),
			HandleRedirectReturnVal:   "",
			HandleRedirectCallTimes:   1,
			SetValue:                  mockCache.EXPECT().SetValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			SetValueReturnError:       nil,
			SetValueCallTimes:         1,
			ExpectedStatusCode:        http.StatusNotFound,
		},
		"Cache Miss With DB Success": {
			requestBody:               &models.RedirectRequestModel{ShortUrlPath: "shortUrl"},
			GetValue:                  mockCache.EXPECT().GetValue(gomock.Any(), gomock.Any()),
			GetValueReturnError:       assert.AnError,
			GetValueReturnVal:         "",
			GetValueCallTimes:         1,
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: nil,
			HandleRedirectReturnVal:   "",
			HandleRedirectCallTimes:   1,
			SetValue:                  mockCache.EXPECT().SetValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			SetValueReturnError:       nil,
			SetValueCallTimes:         1,
			ExpectedStatusCode:        http.StatusOK,
		},
		"Cache Hit": {
			requestBody:               &models.RedirectRequestModel{ShortUrlPath: "shortUrl"},
			GetValue:                  mockCache.EXPECT().GetValue(gomock.Any(), gomock.Any()),
			GetValueReturnError:       nil,
			GetValueReturnVal:         "",
			GetValueCallTimes:         1,
			HandleRedirect:            mockDbService.EXPECT().HandleRedirect(gomock.Any(), gomock.Any()),
			HandleRedirectReturnError: nil,
			HandleRedirectReturnVal:   "",
			HandleRedirectCallTimes:   0,
			SetValue:                  mockCache.EXPECT().SetValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			SetValueReturnError:       nil,
			SetValueCallTimes:         0,
			ExpectedStatusCode:        http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.requestBody)
			test.GetValue.Return(test.GetValueReturnVal, test.GetValueReturnError).Times(test.GetValueCallTimes)
			test.HandleRedirect.Return(test.GetValueReturnVal, test.HandleRedirectReturnError).Times(test.HandleRedirectCallTimes)
			test.SetValue.Return(test.SetValueReturnError).Times(test.SetValueCallTimes)

			assert.Nil(t, err)

			req := httptest.NewRequest("POST", "/redirect", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			handler.HandleRedirect(resp, req)

			assert.Equal(t, test.ExpectedStatusCode, resp.Code)
		})
	}
}
