// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/env (interfaces: FactoryEnv,Env)

package mock

import (
	env "gitlab.connectwisedev.com/platform/platform-common-lib/src/env"
	gomock "github.com/golang/mock/gomock"
	io "io"
)

// Mock of FactoryEnv interface
type MockFactoryEnv struct {
	ctrl     *gomock.Controller
	recorder *_MockFactoryEnvRecorder
}

// Recorder for MockFactoryEnv (not exported)
type _MockFactoryEnvRecorder struct {
	mock *MockFactoryEnv
}

func NewMockFactoryEnv(ctrl *gomock.Controller) *MockFactoryEnv {
	mock := &MockFactoryEnv{ctrl: ctrl}
	mock.recorder = &_MockFactoryEnvRecorder{mock}
	return mock
}

func (_m *MockFactoryEnv) EXPECT() *_MockFactoryEnvRecorder {
	return _m.recorder
}

func (_m *MockFactoryEnv) GetEnv() env.Env {
	ret := _m.ctrl.Call(_m, "GetEnv")
	ret0, _ := ret[0].(env.Env)
	return ret0
}

func (_mr *_MockFactoryEnvRecorder) GetEnv() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEnv")
}

// Mock of Env interface
type MockEnv struct {
	ctrl     *gomock.Controller
	recorder *_MockEnvRecorder
}

// Recorder for MockEnv (not exported)
type _MockEnvRecorder struct {
	mock *MockEnv
}

func NewMockEnv(ctrl *gomock.Controller) *MockEnv {
	mock := &MockEnv{ctrl: ctrl}
	mock.recorder = &_MockEnvRecorder{mock}
	return mock
}

func (_m *MockEnv) EXPECT() *_MockEnvRecorder {
	return _m.recorder
}

func (_m *MockEnv) ExecuteBash(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "ExecuteBash", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvRecorder) ExecuteBash(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ExecuteBash", arg0)
}

func (_m *MockEnv) GetCommandReader(_param0 string, _param1 ...string) (io.ReadCloser, error) {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetCommandReader", _s...)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvRecorder) GetCommandReader(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCommandReader", _s...)
}

func (_m *MockEnv) GetDirectoryFileCount(_param0 string, _param1 ...[]string) (io.ReadCloser, error) {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetDirectoryFileCount", _s...)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvRecorder) GetDirectoryFileCount(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDirectoryFileCount", _s...)
}

func (_m *MockEnv) GetExeDir() (string, error) {
	ret := _m.ctrl.Call(_m, "GetExeDir")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvRecorder) GetExeDir() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetExeDir")
}

func (_m *MockEnv) GetFileReader(_param0 string) (io.ReadCloser, error) {
	ret := _m.ctrl.Call(_m, "GetFileReader", _param0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEnvRecorder) GetFileReader(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetFileReader", arg0)
}
