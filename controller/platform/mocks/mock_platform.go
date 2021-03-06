// Code generated by MockGen. DO NOT EDIT.
// Source: opensearch-cli/controller/platform (interfaces: Controller)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	platform "opensearch-cli/entity/platform"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockController is a mock of Controller interface
type MockController struct {
	ctrl     *gomock.Controller
	recorder *MockControllerMockRecorder
}

// MockControllerMockRecorder is the mock recorder for MockController
type MockControllerMockRecorder struct {
	mock *MockController
}

// NewMockController creates a new mock instance
func NewMockController(ctrl *gomock.Controller) *MockController {
	mock := &MockController{ctrl: ctrl}
	mock.recorder = &MockControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockController) EXPECT() *MockControllerMockRecorder {
	return m.recorder
}

// Curl mocks base method
func (m *MockController) Curl(arg0 context.Context, arg1 platform.CurlCommandRequest) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Curl", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Curl indicates an expected call of Curl
func (mr *MockControllerMockRecorder) Curl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Curl", reflect.TypeOf((*MockController)(nil).Curl), arg0, arg1)
}

// GetDistinctValues mocks base method
func (m *MockController) GetDistinctValues(arg0 context.Context, arg1, arg2 string) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDistinctValues", arg0, arg1, arg2)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDistinctValues indicates an expected call of GetDistinctValues
func (mr *MockControllerMockRecorder) GetDistinctValues(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDistinctValues", reflect.TypeOf((*MockController)(nil).GetDistinctValues), arg0, arg1, arg2)
}
