// Automatically generated by MockGen. DO NOT EDIT!
// Source: gitlab.connectwisedev.com/platform/platform-common-lib/src/kafka (interfaces: ProducerFactory,ConsumerFactory,ProducerService,ConsumerService,Limiter)

package mock

import (
	kafka "gitlab.connectwisedev.com/platform/platform-common-lib/src/kafka"
	encode "gitlab.connectwisedev.com/platform/platform-common-lib/src/kafka/encode"
	gomock "github.com/golang/mock/gomock"
)

// MockProducerFactory is a mock of ProducerFactory interface
type MockProducerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockProducerFactoryMockRecorder
}

// MockProducerFactoryMockRecorder is the mock recorder for MockProducerFactory
type MockProducerFactoryMockRecorder struct {
	mock *MockProducerFactory
}

// NewMockProducerFactory creates a new mock instance
func NewMockProducerFactory(ctrl *gomock.Controller) *MockProducerFactory {
	mock := &MockProducerFactory{ctrl: ctrl}
	mock.recorder = &MockProducerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockProducerFactory) EXPECT() *MockProducerFactoryMockRecorder {
	return _m.recorder
}

// GetConfluentProducerService mocks base method
func (_m *MockProducerFactory) GetConfluentProducerService(_param0 kafka.ProducerConfig) (kafka.ProducerService, error) {
	ret := _m.ctrl.Call(_m, "GetConfluentProducerService", _param0)
	ret0, _ := ret[0].(kafka.ProducerService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfluentProducerService indicates an expected call of GetConfluentProducerService
func (_mr *MockProducerFactoryMockRecorder) GetConfluentProducerService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfluentProducerService", arg0)
}

// GetProducerService mocks base method
func (_m *MockProducerFactory) GetProducerService(_param0 kafka.ProducerConfig) (kafka.ProducerService, error) {
	ret := _m.ctrl.Call(_m, "GetProducerService", _param0)
	ret0, _ := ret[0].(kafka.ProducerService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducerService indicates an expected call of GetProducerService
func (_mr *MockProducerFactoryMockRecorder) GetProducerService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetProducerService", arg0)
}

// MockConsumerFactory is a mock of ConsumerFactory interface
type MockConsumerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerFactoryMockRecorder
}

// MockConsumerFactoryMockRecorder is the mock recorder for MockConsumerFactory
type MockConsumerFactoryMockRecorder struct {
	mock *MockConsumerFactory
}

// NewMockConsumerFactory creates a new mock instance
func NewMockConsumerFactory(ctrl *gomock.Controller) *MockConsumerFactory {
	mock := &MockConsumerFactory{ctrl: ctrl}
	mock.recorder = &MockConsumerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockConsumerFactory) EXPECT() *MockConsumerFactoryMockRecorder {
	return _m.recorder
}

// GetConsumerService mocks base method
func (_m *MockConsumerFactory) GetConsumerService(_param0 kafka.ConsumerConfig) (kafka.ConsumerService, error) {
	ret := _m.ctrl.Call(_m, "GetConsumerService", _param0)
	ret0, _ := ret[0].(kafka.ConsumerService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConsumerService indicates an expected call of GetConsumerService
func (_mr *MockConsumerFactoryMockRecorder) GetConsumerService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConsumerService", arg0)
}

// MockProducerService is a mock of ProducerService interface
type MockProducerService struct {
	ctrl     *gomock.Controller
	recorder *MockProducerServiceMockRecorder
}

// MockProducerServiceMockRecorder is the mock recorder for MockProducerService
type MockProducerServiceMockRecorder struct {
	mock *MockProducerService
}

// NewMockProducerService creates a new mock instance
func NewMockProducerService(ctrl *gomock.Controller) *MockProducerService {
	mock := &MockProducerService{ctrl: ctrl}
	mock.recorder = &MockProducerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockProducerService) EXPECT() *MockProducerServiceMockRecorder {
	return _m.recorder
}

// CloseConnection mocks base method
func (_m *MockProducerService) CloseConnection() error {
	ret := _m.ctrl.Call(_m, "CloseConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseConnection indicates an expected call of CloseConnection
func (_mr *MockProducerServiceMockRecorder) CloseConnection() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CloseConnection")
}

// Push mocks base method
func (_m *MockProducerService) Push(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "Push", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push
func (_mr *MockProducerServiceMockRecorder) Push(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Push", arg0, arg1)
}

// PushEncoder mocks base method
func (_m *MockProducerService) PushEncoder(_param0 string, _param1 encode.Encoder) error {
	ret := _m.ctrl.Call(_m, "PushEncoder", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushEncoder indicates an expected call of PushEncoder
func (_mr *MockProducerServiceMockRecorder) PushEncoder(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PushEncoder", arg0, arg1)
}

// MockConsumerService is a mock of ConsumerService interface
type MockConsumerService struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerServiceMockRecorder
}

// MockConsumerServiceMockRecorder is the mock recorder for MockConsumerService
type MockConsumerServiceMockRecorder struct {
	mock *MockConsumerService
}

// NewMockConsumerService creates a new mock instance
func NewMockConsumerService(ctrl *gomock.Controller) *MockConsumerService {
	mock := &MockConsumerService{ctrl: ctrl}
	mock.recorder = &MockConsumerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockConsumerService) EXPECT() *MockConsumerServiceMockRecorder {
	return _m.recorder
}

// CloseConnection mocks base method
func (_m *MockConsumerService) CloseConnection() error {
	ret := _m.ctrl.Call(_m, "CloseConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseConnection indicates an expected call of CloseConnection
func (_mr *MockConsumerServiceMockRecorder) CloseConnection() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CloseConnection")
}

// Connect mocks base method
func (_m *MockConsumerService) Connect(_param0 *kafka.ConsumerKafkaInOutParams) error {
	ret := _m.ctrl.Call(_m, "Connect", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (_mr *MockConsumerServiceMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Connect", arg0)
}

// MarkOffset mocks base method
func (_m *MockConsumerService) MarkOffset(_param0 string, _param1 int32, _param2 int64) {
	_m.ctrl.Call(_m, "MarkOffset", _param0, _param1, _param2)
}

// MarkOffset indicates an expected call of MarkOffset
func (_mr *MockConsumerServiceMockRecorder) MarkOffset(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MarkOffset", arg0, arg1, arg2)
}

// PullHandler mocks base method
func (_m *MockConsumerService) PullHandler(_param0 kafka.ConsumerHandler) error {
	ret := _m.ctrl.Call(_m, "PullHandler", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullHandler indicates an expected call of PullHandler
func (_mr *MockConsumerServiceMockRecorder) PullHandler(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PullHandler", arg0)
}

// PullHandlerWithLimiter mocks base method
func (_m *MockConsumerService) PullHandlerWithLimiter(_param0 kafka.ConsumerHandler, _param1 kafka.Limiter) error {
	ret := _m.ctrl.Call(_m, "PullHandlerWithLimiter", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PullHandlerWithLimiter indicates an expected call of PullHandlerWithLimiter
func (_mr *MockConsumerServiceMockRecorder) PullHandlerWithLimiter(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PullHandlerWithLimiter", arg0, arg1)
}

// MockLimiter is a mock of Limiter interface
type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
}

// MockLimiterMockRecorder is the mock recorder for MockLimiter
type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

// NewMockLimiter creates a new mock instance
func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return _m.recorder
}

// IsConsumingAllowed mocks base method
func (_m *MockLimiter) IsConsumingAllowed() bool {
	ret := _m.ctrl.Call(_m, "IsConsumingAllowed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConsumingAllowed indicates an expected call of IsConsumingAllowed
func (_mr *MockLimiterMockRecorder) IsConsumingAllowed() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsConsumingAllowed")
}

// Wait mocks base method
func (_m *MockLimiter) Wait() {
	_m.ctrl.Call(_m, "Wait")
}

// Wait indicates an expected call of Wait
func (_mr *MockLimiterMockRecorder) Wait() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Wait")
}