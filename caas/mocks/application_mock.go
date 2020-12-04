// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/caas (interfaces: Application)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	caas "github.com/juju/juju/caas"
	watcher "github.com/juju/juju/core/watcher"
	reflect "reflect"
)

// MockApplication is a mock of Application interface
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockApplication) Delete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete")
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockApplicationMockRecorder) Delete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockApplication)(nil).Delete))
}

// Ensure mocks base method
func (m *MockApplication) Ensure(arg0 caas.ApplicationConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ensure", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ensure indicates an expected call of Ensure
func (mr *MockApplicationMockRecorder) Ensure(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ensure", reflect.TypeOf((*MockApplication)(nil).Ensure), arg0)
}

// Exists mocks base method
func (m *MockApplication) Exists() (caas.DeploymentState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists")
	ret0, _ := ret[0].(caas.DeploymentState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists
func (mr *MockApplicationMockRecorder) Exists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockApplication)(nil).Exists))
}

// State mocks base method
func (m *MockApplication) State() (caas.ApplicationState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(caas.ApplicationState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// State indicates an expected call of State
func (mr *MockApplicationMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockApplication)(nil).State))
}

// Units mocks base method
func (m *MockApplication) Units() ([]caas.Unit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Units")
	ret0, _ := ret[0].([]caas.Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units
func (mr *MockApplicationMockRecorder) Units() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockApplication)(nil).Units))
}

// UpdatePorts mocks base method
func (m *MockApplication) UpdatePorts(arg0 []caas.ServicePort, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePorts", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePorts indicates an expected call of UpdatePorts
func (mr *MockApplicationMockRecorder) UpdatePorts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePorts", reflect.TypeOf((*MockApplication)(nil).UpdatePorts), arg0, arg1)
}

// UpdateService mocks base method
func (m *MockApplication) UpdateService(arg0 caas.ServiceParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateService", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateService indicates an expected call of UpdateService
func (mr *MockApplicationMockRecorder) UpdateService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateService", reflect.TypeOf((*MockApplication)(nil).UpdateService), arg0)
}

// Watch mocks base method
func (m *MockApplication) Watch() (watcher.NotifyWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch")
	ret0, _ := ret[0].(watcher.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockApplicationMockRecorder) Watch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockApplication)(nil).Watch))
}

// WatchReplicas mocks base method
func (m *MockApplication) WatchReplicas() (watcher.NotifyWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchReplicas")
	ret0, _ := ret[0].(watcher.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchReplicas indicates an expected call of WatchReplicas
func (mr *MockApplicationMockRecorder) WatchReplicas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchReplicas", reflect.TypeOf((*MockApplication)(nil).WatchReplicas))
}
