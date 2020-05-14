// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/persistency/memcached (interfaces: Cache)

// Package memcached is a generated GoMock package.
package memcached

import (
	memcache "github.com/bradfitz/gomemcache/memcache"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCache is a mock of Cache interface
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockCache) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockCacheMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCache)(nil).Delete), arg0)
}

// Get mocks base method
func (m *MockCache) Get(arg0 string) (*memcache.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*memcache.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockCacheMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), arg0)
}

// Set mocks base method
func (m *MockCache) Set(arg0 *memcache.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set
func (mr *MockCacheMockRecorder) Set(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), arg0)
}
