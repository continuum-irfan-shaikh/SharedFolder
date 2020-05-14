// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/exec (interfaces: Command)

package mock

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Command interface
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *_MockCommandRecorder
}

// Recorder for MockCommand (not exported)
type _MockCommandRecorder struct {
	mock *MockCommand
}

func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &_MockCommandRecorder{mock}
	return mock
}

func (_m *MockCommand) EXPECT() *_MockCommandRecorder {
	return _m.recorder
}

func (_m *MockCommand) Run(_param0 string, _param1 ...string) error {
	_s := []interface{}{_param0}
	for _, _x := range _param1 {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Run", _s...)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockCommandRecorder) Run(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0}, arg1...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Run", _s...)
}
