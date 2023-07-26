// Code generated by MockGen. DO NOT EDIT.
// Source: justanother.org/justanotherbotkit/transport (interfaces: Transport)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	transport "justanother.org/justanotherbotkit/transport"
)

// MockTransport is a mock of Transport interface.
type MockTransport struct {
	ctrl     *gomock.Controller
	recorder *MockTransportMockRecorder
}

// MockTransportMockRecorder is the mock recorder for MockTransport.
type MockTransportMockRecorder struct {
	mock *MockTransport
}

// NewMockTransport creates a new mock instance.
func NewMockTransport(ctrl *gomock.Controller) *MockTransport {
	mock := &MockTransport{ctrl: ctrl}
	mock.recorder = &MockTransportMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransport) EXPECT() *MockTransportMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockTransport) SendMessage(arg0 context.Context, arg1 string, arg2 ...transport.MsgOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockTransportMockRecorder) SendMessage(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockTransport)(nil).SendMessage), varargs...)
}

// Start mocks base method.
func (m *MockTransport) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockTransportMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTransport)(nil).Start), arg0)
}
