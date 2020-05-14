// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.connectwisedev.com/platform/platform-tasking-service/src/integration (interfaces: AutomationEngine)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entities "gitlab.connectwisedev.com/platform/platform-tasking-service/src/entities"
	reflect "reflect"
)

// MockAutomationEngine is a mock of AutomationEngine interface
type MockAutomationEngine struct {
	ctrl     *gomock.Controller
	recorder *MockAutomationEngineMockRecorder
}

// MockAutomationEngineMockRecorder is the mock recorder for MockAutomationEngine
type MockAutomationEngineMockRecorder struct {
	mock *MockAutomationEngine
}

// NewMockAutomationEngine creates a new mock instance
func NewMockAutomationEngine(ctrl *gomock.Controller) *MockAutomationEngine {
	mock := &MockAutomationEngine{ctrl: ctrl}
	mock.recorder = &MockAutomationEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAutomationEngine) EXPECT() *MockAutomationEngineMockRecorder {
	return m.recorder
}

// RemovePolicy mocks base method
func (m *MockAutomationEngine) RemovePolicy(arg0 context.Context, arg1 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "RemovePolicy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePolicy indicates an expected call of RemovePolicy
func (mr *MockAutomationEngineMockRecorder) RemovePolicy(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePolicy", reflect.TypeOf((*MockAutomationEngine)(nil).RemovePolicy), arg0, arg1)
}

// UpdateRemotePolicies mocks base method
func (m *MockAutomationEngine) UpdateRemotePolicies(arg0 context.Context, arg1 []entities.TriggerDefinition) (string, error) {
	ret := m.ctrl.Call(m, "UpdateRemotePolicies", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRemotePolicies indicates an expected call of UpdateRemotePolicies
func (mr *MockAutomationEngineMockRecorder) UpdateRemotePolicies(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRemotePolicies", reflect.TypeOf((*MockAutomationEngine)(nil).UpdateRemotePolicies), arg0, arg1)
}
