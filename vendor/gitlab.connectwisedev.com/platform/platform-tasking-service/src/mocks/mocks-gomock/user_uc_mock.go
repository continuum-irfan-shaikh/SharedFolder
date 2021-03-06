// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/services/tasks (interfaces: UserUC)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	reflect "reflect"
)

// MockUserUC is a mock of UserUC interface
type MockUserUC struct {
	ctrl     *gomock.Controller
	recorder *MockUserUCMockRecorder
}

// MockUserUCMockRecorder is the mock recorder for MockUserUC
type MockUserUCMockRecorder struct {
	mock *MockUserUC
}

// NewMockUserUC creates a new mock instance
func NewMockUserUC(ctrl *gomock.Controller) *MockUserUC {
	mock := &MockUserUC{ctrl: ctrl}
	mock.recorder = &MockUserUCMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserUC) EXPECT() *MockUserUCMockRecorder {
	return m.recorder
}

// EndpointsFromAsset mocks base method
func (m *MockUserUC) EndpointsFromAsset(arg0 context.Context, arg1 []string) ([]entities.Endpoints, error) {
	ret := m.ctrl.Call(m, "EndpointsFromAsset", arg0, arg1)
	ret0, _ := ret[0].([]entities.Endpoints)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EndpointsFromAsset indicates an expected call of EndpointsFromAsset
func (mr *MockUserUCMockRecorder) EndpointsFromAsset(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndpointsFromAsset", reflect.TypeOf((*MockUserUC)(nil).EndpointsFromAsset), arg0, arg1)
}

// SaveEndpoints mocks base method
func (m *MockUserUC) SaveEndpoints(arg0 context.Context, arg1 []entities.Endpoints) {
	m.ctrl.Call(m, "SaveEndpoints", arg0, arg1)
}

// SaveEndpoints indicates an expected call of SaveEndpoints
func (mr *MockUserUCMockRecorder) SaveEndpoints(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEndpoints", reflect.TypeOf((*MockUserUC)(nil).SaveEndpoints), arg0, arg1)
}
