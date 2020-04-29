// Code generated by MockGen. DO NOT EDIT.
// Source: channelsolution.go

// Package status is a generated GoMock package.
package status

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockResponder is a mock of Responder interface
type MockResponder struct {
	ctrl     *gomock.Controller
	recorder *MockResponderMockRecorder
}

// MockResponderMockRecorder is the mock recorder for MockResponder
type MockResponderMockRecorder struct {
	mock *MockResponder
}

// NewMockResponder creates a new mock instance
func NewMockResponder(ctrl *gomock.Controller) *MockResponder {
	mock := &MockResponder{ctrl: ctrl}
	mock.recorder = &MockResponderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResponder) EXPECT() *MockResponderMockRecorder {
	return m.recorder
}

// Reply mocks base method
func (m *MockResponder) Reply(fr *friendResponse) error {
	ret := m.ctrl.Call(m, "Reply", fr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reply indicates an expected call of Reply
func (mr *MockResponderMockRecorder) Reply(fr interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reply", reflect.TypeOf((*MockResponder)(nil).Reply), fr)
}

// IsStateless mocks base method
func (m *MockResponder) IsStateless() bool {
	ret := m.ctrl.Call(m, "IsStateless")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsStateless indicates an expected call of IsStateless
func (mr *MockResponderMockRecorder) IsStateless() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsStateless", reflect.TypeOf((*MockResponder)(nil).IsStateless))
}