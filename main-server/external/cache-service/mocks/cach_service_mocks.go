// Code generated by MockGen. DO NOT EDIT.
// Source: D:\Programming\Go\url-shortner\main-server\external\cache-service\cache_service.go

// Package mock_cacheservice is a generated GoMock package.
package mock_cacheservice

import (
	io "io"
	models "main-server/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCacheServiceInterface is a mock of CacheServiceInterface interface.
type MockCacheServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCacheServiceInterfaceMockRecorder
}

// MockCacheServiceInterfaceMockRecorder is the mock recorder for MockCacheServiceInterface.
type MockCacheServiceInterfaceMockRecorder struct {
	mock *MockCacheServiceInterface
}

// NewMockCacheServiceInterface creates a new mock instance.
func NewMockCacheServiceInterface(ctrl *gomock.Controller) *MockCacheServiceInterface {
	mock := &MockCacheServiceInterface{ctrl: ctrl}
	mock.recorder = &MockCacheServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheServiceInterface) EXPECT() *MockCacheServiceInterfaceMockRecorder {
	return m.recorder
}

// HandleRedirect mocks base method.
func (m *MockCacheServiceInterface) HandleRedirect(body io.Reader, requestId string) (*models.RedirectResponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleRedirect", body, requestId)
	ret0, _ := ret[0].(*models.RedirectResponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleRedirect indicates an expected call of HandleRedirect.
func (mr *MockCacheServiceInterfaceMockRecorder) HandleRedirect(body, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleRedirect", reflect.TypeOf((*MockCacheServiceInterface)(nil).HandleRedirect), body, requestId)
}