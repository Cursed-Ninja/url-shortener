// Code generated by MockGen. DO NOT EDIT.
// Source: D:\Programming\Go\url-shortner\cache-server\internal\handlers\handlers.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandlerInterface is a mock of HandlerInterface interface.
type MockHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerInterfaceMockRecorder
}

// MockHandlerInterfaceMockRecorder is the mock recorder for MockHandlerInterface.
type MockHandlerInterfaceMockRecorder struct {
	mock *MockHandlerInterface
}

// NewMockHandlerInterface creates a new mock instance.
func NewMockHandlerInterface(ctrl *gomock.Controller) *MockHandlerInterface {
	mock := &MockHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandlerInterface) EXPECT() *MockHandlerInterfaceMockRecorder {
	return m.recorder
}

// HandleRedirect mocks base method.
func (m *MockHandlerInterface) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleRedirect", w, r)
}

// HandleRedirect indicates an expected call of HandleRedirect.
func (mr *MockHandlerInterfaceMockRecorder) HandleRedirect(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleRedirect", reflect.TypeOf((*MockHandlerInterface)(nil).HandleRedirect), w, r)
}