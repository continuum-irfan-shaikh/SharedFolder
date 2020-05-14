// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/plugin/wmi (interfaces: Wrapper)

package wmiMock

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Wrapper interface
type MockWrapper struct {
	ctrl     *gomock.Controller
	recorder *_MockWrapperRecorder
}

// Recorder for MockWrapper (not exported)
type _MockWrapperRecorder struct {
	mock *MockWrapper
}

func NewMockWrapper(ctrl *gomock.Controller) *MockWrapper {
	mock := &MockWrapper{ctrl: ctrl}
	mock.recorder = &_MockWrapperRecorder{mock}
	return mock
}

func (_m *MockWrapper) EXPECT() *_MockWrapperRecorder {
	return _m.recorder
}

func (_m *MockWrapper) CreateQuery(_param0 interface{}, _param1 string) string {
	ret := _m.ctrl.Call(_m, "CreateQuery", _param0, _param1)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockWrapperRecorder) CreateQuery(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateQuery", arg0, arg1)
}

func (_m *MockWrapper) Query(_param0 string, _param1 interface{}, _param2 ...interface{}) error {
	_s := []interface{}{_param0, _param1}
	for _, _x := range _param2 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Query", _s...)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockWrapperRecorder) Query(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Query", _s...)
}
