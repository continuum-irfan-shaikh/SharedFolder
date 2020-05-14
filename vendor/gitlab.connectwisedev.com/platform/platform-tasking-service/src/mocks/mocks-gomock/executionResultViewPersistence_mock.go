// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/models (interfaces: ExecutionResultViewPersistence)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	gocql "github.com/gocql/gocql"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExecutionResultViewPersistence is a mock of ExecutionResultViewPersistence interface
type MockExecutionResultViewPersistence struct {
	ctrl     *gomock.Controller
	recorder *MockExecutionResultViewPersistenceMockRecorder
}

// MockExecutionResultViewPersistenceMockRecorder is the mock recorder for MockExecutionResultViewPersistence
type MockExecutionResultViewPersistenceMockRecorder struct {
	mock *MockExecutionResultViewPersistence
}

// NewMockExecutionResultViewPersistence creates a new mock instance
func NewMockExecutionResultViewPersistence(ctrl *gomock.Controller) *MockExecutionResultViewPersistence {
	mock := &MockExecutionResultViewPersistence{ctrl: ctrl}
	mock.recorder = &MockExecutionResultViewPersistenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutionResultViewPersistence) EXPECT() *MockExecutionResultViewPersistenceMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockExecutionResultViewPersistence) Get(arg0 context.Context, arg1 string, arg2 gocql.UUID, arg3 int, arg4 bool) ([]*models.ExecutionResultView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*models.ExecutionResultView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockExecutionResultViewPersistenceMockRecorder) Get(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockExecutionResultViewPersistence)(nil).Get), arg0, arg1, arg2, arg3, arg4)
}

// History mocks base method
func (m *MockExecutionResultViewPersistence) History(arg0 context.Context, arg1 string, arg2, arg3 gocql.UUID, arg4 int, arg5 bool) ([]*models.ExecutionResultView, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "History", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]*models.ExecutionResultView)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// History indicates an expected call of History
func (mr *MockExecutionResultViewPersistenceMockRecorder) History(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "History", reflect.TypeOf((*MockExecutionResultViewPersistence)(nil).History), arg0, arg1, arg2, arg3, arg4, arg5)
}
