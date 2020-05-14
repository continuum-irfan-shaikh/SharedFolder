// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/pluginUtils (interfaces: IOReaderFactory,IOWriterFactory)

package mock

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
)

// Mock of IOReaderFactory interface
type MockIOReaderFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockIOReaderFactoryRecorder
}

// Recorder for MockIOReaderFactory (not exported)
type _MockIOReaderFactoryRecorder struct {
	mock *MockIOReaderFactory
}

func NewMockIOReaderFactory(ctrl *gomock.Controller) *MockIOReaderFactory {
	mock := &MockIOReaderFactory{ctrl: ctrl}
	mock.recorder = &_MockIOReaderFactoryRecorder{mock}
	return mock
}

func (_m *MockIOReaderFactory) EXPECT() *_MockIOReaderFactoryRecorder {
	return _m.recorder
}

func (_m *MockIOReaderFactory) GetReader() io.Reader {
	ret := _m.ctrl.Call(_m, "GetReader")
	ret0, _ := ret[0].(io.Reader)
	return ret0
}

func (_mr *_MockIOReaderFactoryRecorder) GetReader() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetReader")
}

// Mock of IOWriterFactory interface
type MockIOWriterFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockIOWriterFactoryRecorder
}

// Recorder for MockIOWriterFactory (not exported)
type _MockIOWriterFactoryRecorder struct {
	mock *MockIOWriterFactory
}

func NewMockIOWriterFactory(ctrl *gomock.Controller) *MockIOWriterFactory {
	mock := &MockIOWriterFactory{ctrl: ctrl}
	mock.recorder = &_MockIOWriterFactoryRecorder{mock}
	return mock
}

func (_m *MockIOWriterFactory) EXPECT() *_MockIOWriterFactoryRecorder {
	return _m.recorder
}

func (_m *MockIOWriterFactory) GetWriter() io.Writer {
	ret := _m.ctrl.Call(_m, "GetWriter")
	ret0, _ := ret[0].(io.Writer)
	return ret0
}

func (_mr *_MockIOWriterFactoryRecorder) GetWriter() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetWriter")
}
