// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/config (interfaces: ConfigurationService)

package mock

import (
	config "gitlab.connectwisedev.com/platform/platform-common-lib/src/config"
	gomock "github.com/golang/mock/gomock"
)

// Mock of ConfigurationService interface
type MockConfigurationService struct {
	ctrl     *gomock.Controller
	recorder *_MockConfigurationServiceRecorder
}

// Recorder for MockConfigurationService (not exported)
type _MockConfigurationServiceRecorder struct {
	mock *MockConfigurationService
}

func NewMockConfigurationService(ctrl *gomock.Controller) *MockConfigurationService {
	mock := &MockConfigurationService{ctrl: ctrl}
	mock.recorder = &_MockConfigurationServiceRecorder{mock}
	return mock
}

func (_m *MockConfigurationService) EXPECT() *_MockConfigurationServiceRecorder {
	return _m.recorder
}

func (_m *MockConfigurationService) Update(_param0 config.Configuration) ([]config.UpdatedConfig, error) {
	ret := _m.ctrl.Call(_m, "Update", _param0)
	ret0, _ := ret[0].([]config.UpdatedConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockConfigurationServiceRecorder) Update(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0)
}
