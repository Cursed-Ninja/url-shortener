package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	mock_database "url-shortner-database/internal/database/mocks"
	"url-shortner-database/internal/handlers"
	"url-shortner-database/internal/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleShorten(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	mockObj := mock_database.NewMockDBInterface(mockCtrl)
	mockObj.EXPECT().FindOne(gomock.Any()).Return(models.URL{}, errors.New("Not Found")).AnyTimes()

	handler := handlers.NewBaseHandler(logger, mockObj)

	tests := map[string]struct {
		reqBody              *models.ShortenRequestModel
		InsertOne            *gomock.Call
		InsertOneReturnError error
		InsertOneCall        int
		ExpectedStatusCode   int
	}{
		"Empty Request Body": {
			reqBody:              nil,
			InsertOne:            mockObj.EXPECT().InsertOne(gomock.Any()),
			InsertOneReturnError: nil,
			InsertOneCall:        0,
			ExpectedStatusCode:   http.StatusBadRequest,
		},
		"Error InsertOne": {
			reqBody:              &models.ShortenRequestModel{Url: "http://www.google.com"},
			InsertOne:            mockObj.EXPECT().InsertOne(gomock.Any()),
			InsertOneCall:        1,
			InsertOneReturnError: assert.AnError,
			ExpectedStatusCode:   http.StatusInternalServerError,
		},
		"Success": {
			reqBody:              &models.ShortenRequestModel{Url: "http://www.google.com"},
			InsertOne:            mockObj.EXPECT().InsertOne(gomock.Any()),
			InsertOneCall:        1,
			InsertOneReturnError: nil,
			ExpectedStatusCode:   http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.InsertOne.Return(test.InsertOneReturnError).Times(test.InsertOneCall)

			body, err := json.Marshal(test.reqBody)

			if err != nil {
				t.Error("Error marshalling request body")
			}

			req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			handler.HandleShorten(resp, req)
			assert.Equal(t, test.ExpectedStatusCode, resp.Code, resp.Result().Status)
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	logger := zap.NewNop().Sugar()

	mockCtrl := gomock.NewController(t)
	mockObj := mock_database.NewMockDBInterface(mockCtrl)

	handler := handlers.NewBaseHandler(logger, mockObj)

	tests := map[string]struct {
		reqBody            *models.RedirectRequestModel
		FindOne            *gomock.Call
		FindOneReturnError error
		FindOneReturnUrl   models.URL
		FindOneCall        int
		ExpectedStatusCode int
	}{
		"Empty Request URL": {
			reqBody:            nil,
			FindOne:            mockObj.EXPECT().FindOne(gomock.Any()),
			FindOneReturnError: nil,
			FindOneReturnUrl:   models.URL{},
			FindOneCall:        0,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"Error Find One": {
			reqBody:            &models.RedirectRequestModel{ShortUrlPath: "test"},
			FindOne:            mockObj.EXPECT().FindOne(gomock.Any()),
			FindOneReturnError: assert.AnError,
			FindOneReturnUrl:   models.URL{},
			FindOneCall:        1,
			ExpectedStatusCode: http.StatusNotFound,
		},
		"Success": {
			reqBody:            &models.RedirectRequestModel{ShortUrlPath: "test"},
			FindOne:            mockObj.EXPECT().FindOne(gomock.Any()),
			FindOneReturnError: nil,
			FindOneReturnUrl:   models.URL{ShortUrlPath: "test", OriginalUrl: "http://www.google.com", ExpiresAt: time.Now().AddDate(0, 1, 0)},
			FindOneCall:        1,
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.FindOne.Return(test.FindOneReturnUrl, test.FindOneReturnError).Times(test.FindOneCall)

			body, err := json.Marshal(test.reqBody)

			if err != nil {
				t.Error("Error marshalling request body")
			}

			req := httptest.NewRequest("POST", "/redirect", bytes.NewBuffer(body))
			resp := httptest.NewRecorder()
			handler.HandleRedirect(resp, req)
			assert.Equal(t, test.ExpectedStatusCode, resp.Code, resp.Result().Status)
		})
	}
}
