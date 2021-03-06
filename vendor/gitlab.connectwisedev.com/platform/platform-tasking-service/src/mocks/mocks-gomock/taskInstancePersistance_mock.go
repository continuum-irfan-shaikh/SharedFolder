// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/models (interfaces: TaskInstancePersistence)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gocql "github.com/gocql/gocql"
	gomock "github.com/golang/mock/gomock"
	models "gitlab.connectwisedev.com/platform/platform-tasking-service/src/models"
	reflect "reflect"
	time "time"
)

// MockTaskInstancePersistence is a mock of TaskInstancePersistence interface
type MockTaskInstancePersistence struct {
	ctrl     *gomock.Controller
	recorder *MockTaskInstancePersistenceMockRecorder
}

// MockTaskInstancePersistenceMockRecorder is the mock recorder for MockTaskInstancePersistence
type MockTaskInstancePersistenceMockRecorder struct {
	mock *MockTaskInstancePersistence
}

// NewMockTaskInstancePersistence creates a new mock instance
func NewMockTaskInstancePersistence(ctrl *gomock.Controller) *MockTaskInstancePersistence {
	mock := &MockTaskInstancePersistence{ctrl: ctrl}
	mock.recorder = &MockTaskInstancePersistenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTaskInstancePersistence) EXPECT() *MockTaskInstancePersistenceMockRecorder {
	return m.recorder
}

// DeleteBatch mocks base method
func (m *MockTaskInstancePersistence) DeleteBatch(arg0 context.Context, arg1 []models.TaskInstance) error {
	ret := m.ctrl.Call(m, "DeleteBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch
func (mr *MockTaskInstancePersistenceMockRecorder) DeleteBatch(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockTaskInstancePersistence)(nil).DeleteBatch), arg0, arg1)
}

// GetByIDs mocks base method
func (m *MockTaskInstancePersistence) GetByIDs(arg0 context.Context, arg1 ...gocql.UUID) ([]models.TaskInstance, error) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByIDs", varargs...)
	ret0, _ := ret[0].([]models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDs indicates an expected call of GetByIDs
func (mr *MockTaskInstancePersistenceMockRecorder) GetByIDs(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetByIDs), varargs...)
}

// GetByStartedAtAfter mocks base method
func (m *MockTaskInstancePersistence) GetByStartedAtAfter(arg0 context.Context, arg1 string, arg2, arg3 time.Time) ([]models.TaskInstance, error) {
	ret := m.ctrl.Call(m, "GetByStartedAtAfter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByStartedAtAfter indicates an expected call of GetByStartedAtAfter
func (mr *MockTaskInstancePersistenceMockRecorder) GetByStartedAtAfter(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByStartedAtAfter", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetByStartedAtAfter), arg0, arg1, arg2, arg3)
}

// GetByTaskID mocks base method
func (m *MockTaskInstancePersistence) GetByTaskID(arg0 context.Context, arg1 gocql.UUID) ([]models.TaskInstance, error) {
	ret := m.ctrl.Call(m, "GetByTaskID", arg0, arg1)
	ret0, _ := ret[0].([]models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTaskID indicates an expected call of GetByTaskID
func (mr *MockTaskInstancePersistenceMockRecorder) GetByTaskID(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTaskID", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetByTaskID), arg0, arg1)
}

// GetInstancesCountByTaskID mocks base method
func (m *MockTaskInstancePersistence) GetInstancesCountByTaskID(arg0 context.Context, arg1 gocql.UUID) (int, error) {
	ret := m.ctrl.Call(m, "GetInstancesCountByTaskID", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstancesCountByTaskID indicates an expected call of GetInstancesCountByTaskID
func (mr *MockTaskInstancePersistenceMockRecorder) GetInstancesCountByTaskID(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstancesCountByTaskID", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetInstancesCountByTaskID), arg0, arg1)
}

// GetMinimalInstanceByID mocks base method
func (m *MockTaskInstancePersistence) GetMinimalInstanceByID(arg0 context.Context, arg1 gocql.UUID) (models.TaskInstance, error) {
	ret := m.ctrl.Call(m, "GetMinimalInstanceByID", arg0, arg1)
	ret0, _ := ret[0].(models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinimalInstanceByID indicates an expected call of GetMinimalInstanceByID
func (mr *MockTaskInstancePersistenceMockRecorder) GetMinimalInstanceByID(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinimalInstanceByID", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetMinimalInstanceByID), arg0, arg1)
}

// GetNearestInstanceAfter mocks base method
func (m *MockTaskInstancePersistence) GetNearestInstanceAfter(arg0 gocql.UUID, arg1 time.Time) (models.TaskInstance, error) {
	ret := m.ctrl.Call(m, "GetNearestInstanceAfter", arg0, arg1)
	ret0, _ := ret[0].(models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNearestInstanceAfter indicates an expected call of GetNearestInstanceAfter
func (mr *MockTaskInstancePersistenceMockRecorder) GetNearestInstanceAfter(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNearestInstanceAfter", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetNearestInstanceAfter), arg0, arg1)
}

// GetTopInstancesByTaskID mocks base method
func (m *MockTaskInstancePersistence) GetTopInstancesByTaskID(arg0 context.Context, arg1 gocql.UUID) ([]models.TaskInstance, error) {
	ret := m.ctrl.Call(m, "GetTopInstancesByTaskID", arg0, arg1)
	ret0, _ := ret[0].([]models.TaskInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopInstancesByTaskID indicates an expected call of GetTopInstancesByTaskID
func (mr *MockTaskInstancePersistenceMockRecorder) GetTopInstancesByTaskID(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopInstancesByTaskID", reflect.TypeOf((*MockTaskInstancePersistence)(nil).GetTopInstancesByTaskID), arg0, arg1)
}

// Insert mocks base method
func (m *MockTaskInstancePersistence) Insert(arg0 context.Context, arg1 models.TaskInstance) error {
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockTaskInstancePersistenceMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockTaskInstancePersistence)(nil).Insert), arg0, arg1)
}

// UpdateStatuses mocks base method
func (m *MockTaskInstancePersistence) UpdateStatuses(arg0 context.Context, arg1 models.TaskInstance) error {
	ret := m.ctrl.Call(m, "UpdateStatuses", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatuses indicates an expected call of UpdateStatuses
func (mr *MockTaskInstancePersistenceMockRecorder) UpdateStatuses(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatuses", reflect.TypeOf((*MockTaskInstancePersistence)(nil).UpdateStatuses), arg0, arg1)
}
