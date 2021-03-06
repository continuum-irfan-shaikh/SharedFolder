// Code generated by MockGen. DO NOT EDIT.
// Source: ./permission.go

// Package permission is a generated GoMock package.
package permission

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockuserUC is a mock of userUC interface
type MockuserUC struct {
	ctrl     *gomock.Controller
	recorder *MockuserUCMockRecorder
}

// MockuserUCMockRecorder is the mock recorder for MockuserUC
type MockuserUCMockRecorder struct {
	mock *MockuserUC
}

// NewMockuserUC creates a new mock instance
func NewMockuserUC(ctrl *gomock.Controller) *MockuserUC {
	mock := &MockuserUC{ctrl: ctrl}
	mock.recorder = &MockuserUCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockuserUC) EXPECT() *MockuserUCMockRecorder {
	return m.recorder
}

// Sites mocks base method
func (m *MockuserUC) Sites(ctx context.Context) ([]string, error) {
	ret := m.ctrl.Call(m, "Sites", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sites indicates an expected call of Sites
func (mr *MockuserUCMockRecorder) Sites(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sites", reflect.TypeOf((*MockuserUC)(nil).Sites), ctx)
}

// Endpoints mocks base method
func (m *MockuserUC) Endpoints(ctx context.Context, sites []string) ([]string, error) {
	ret := m.ctrl.Call(m, "Endpoints", ctx, sites)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Endpoints indicates an expected call of Endpoints
func (mr *MockuserUCMockRecorder) Endpoints(ctx, sites interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Endpoints", reflect.TypeOf((*MockuserUC)(nil).Endpoints), ctx, sites)
}
