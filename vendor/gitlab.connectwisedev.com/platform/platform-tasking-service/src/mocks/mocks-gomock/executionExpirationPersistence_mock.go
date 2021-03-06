// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/models (interfaces: ExecutionExpirationPersistence)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	gocql "github.com/gocql/gocql"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockExecutionExpirationPersistence is a mock of ExecutionExpirationPersistence interface
type MockExecutionExpirationPersistence struct {
	ctrl     *gomock.Controller
	recorder *MockExecutionExpirationPersistenceMockRecorder
}

// MockExecutionExpirationPersistenceMockRecorder is the mock recorder for MockExecutionExpirationPersistence
type MockExecutionExpirationPersistenceMockRecorder struct {
	mock *MockExecutionExpirationPersistence
}

// NewMockExecutionExpirationPersistence creates a new mock instance
func NewMockExecutionExpirationPersistence(ctrl *gomock.Controller) *MockExecutionExpirationPersistence {
	mock := &MockExecutionExpirationPersistence{ctrl: ctrl}
	mock.recorder = &MockExecutionExpirationPersistenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutionExpirationPersistence) EXPECT() *MockExecutionExpirationPersistenceMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockExecutionExpirationPersistence) Delete(arg0 models.ExecutionExpiration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockExecutionExpirationPersistenceMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockExecutionExpirationPersistence)(nil).Delete), arg0)
}

// GetByExpirationTime mocks base method
func (m *MockExecutionExpirationPersistence) GetByExpirationTime(arg0 context.Context, arg1 time.Time) ([]models.ExecutionExpiration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByExpirationTime", arg0, arg1)
	ret0, _ := ret[0].([]models.ExecutionExpiration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByExpirationTime indicates an expected call of GetByExpirationTime
func (mr *MockExecutionExpirationPersistenceMockRecorder) GetByExpirationTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByExpirationTime", reflect.TypeOf((*MockExecutionExpirationPersistence)(nil).GetByExpirationTime), arg0, arg1)
}

// GetByTaskInstanceIDs mocks base method
func (m *MockExecutionExpirationPersistence) GetByTaskInstanceIDs(arg0 string, arg1 []gocql.UUID) ([]models.ExecutionExpiration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTaskInstanceIDs", arg0, arg1)
	ret0, _ := ret[0].([]models.ExecutionExpiration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTaskInstanceIDs indicates an expected call of GetByTaskInstanceIDs
func (mr *MockExecutionExpirationPersistenceMockRecorder) GetByTaskInstanceIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTaskInstanceIDs", reflect.TypeOf((*MockExecutionExpirationPersistence)(nil).GetByTaskInstanceIDs), arg0, arg1)
}

// InsertExecutionExpiration mocks base method
func (m *MockExecutionExpirationPersistence) InsertExecutionExpiration(arg0 context.Context, arg1 models.ExecutionExpiration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertExecutionExpiration", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertExecutionExpiration indicates an expected call of InsertExecutionExpiration
func (mr *MockExecutionExpirationPersistenceMockRecorder) InsertExecutionExpiration(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertExecutionExpiration", reflect.TypeOf((*MockExecutionExpirationPersistence)(nil).InsertExecutionExpiration), arg0, arg1)
}
