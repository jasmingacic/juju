// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/provider/openstack (interfaces: SSLHostnameConfig,Networking)

// Package openstack is a generated GoMock package.
package openstack

import (
	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	instance "github.com/juju/juju/core/instance"
	network "github.com/juju/juju/core/network"
	neutron "gopkg.in/goose.v2/neutron"
	nova "gopkg.in/goose.v2/nova"
	reflect "reflect"
)

// MockSSLHostnameConfig is a mock of SSLHostnameConfig interface
type MockSSLHostnameConfig struct {
	ctrl     *gomock.Controller
	recorder *MockSSLHostnameConfigMockRecorder
}

// MockSSLHostnameConfigMockRecorder is the mock recorder for MockSSLHostnameConfig
type MockSSLHostnameConfigMockRecorder struct {
	mock *MockSSLHostnameConfig
}

// NewMockSSLHostnameConfig creates a new mock instance
func NewMockSSLHostnameConfig(ctrl *gomock.Controller) *MockSSLHostnameConfig {
	mock := &MockSSLHostnameConfig{ctrl: ctrl}
	mock.recorder = &MockSSLHostnameConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSSLHostnameConfig) EXPECT() *MockSSLHostnameConfigMockRecorder {
	return m.recorder
}

// SSLHostnameVerification mocks base method
func (m *MockSSLHostnameConfig) SSLHostnameVerification() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SSLHostnameVerification")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SSLHostnameVerification indicates an expected call of SSLHostnameVerification
func (mr *MockSSLHostnameConfigMockRecorder) SSLHostnameVerification() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SSLHostnameVerification", reflect.TypeOf((*MockSSLHostnameConfig)(nil).SSLHostnameVerification))
}

// MockNetworking is a mock of Networking interface
type MockNetworking struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingMockRecorder
}

// MockNetworkingMockRecorder is the mock recorder for MockNetworking
type MockNetworkingMockRecorder struct {
	mock *MockNetworking
}

// NewMockNetworking creates a new mock instance
func NewMockNetworking(ctrl *gomock.Controller) *MockNetworking {
	mock := &MockNetworking{ctrl: ctrl}
	mock.recorder = &MockNetworkingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworking) EXPECT() *MockNetworkingMockRecorder {
	return m.recorder
}

// AllocatePublicIP mocks base method
func (m *MockNetworking) AllocatePublicIP(arg0 instance.Id) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocatePublicIP", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocatePublicIP indicates an expected call of AllocatePublicIP
func (mr *MockNetworkingMockRecorder) AllocatePublicIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocatePublicIP", reflect.TypeOf((*MockNetworking)(nil).AllocatePublicIP), arg0)
}

// CreatePort mocks base method
func (m *MockNetworking) CreatePort(arg0, arg1 string, arg2 network.Id) (*neutron.PortV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePort", arg0, arg1, arg2)
	ret0, _ := ret[0].(*neutron.PortV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePort indicates an expected call of CreatePort
func (mr *MockNetworkingMockRecorder) CreatePort(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePort", reflect.TypeOf((*MockNetworking)(nil).CreatePort), arg0, arg1, arg2)
}

// DefaultNetworks mocks base method
func (m *MockNetworking) DefaultNetworks() ([]nova.ServerNetworks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultNetworks")
	ret0, _ := ret[0].([]nova.ServerNetworks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DefaultNetworks indicates an expected call of DefaultNetworks
func (mr *MockNetworkingMockRecorder) DefaultNetworks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultNetworks", reflect.TypeOf((*MockNetworking)(nil).DefaultNetworks))
}

// DeletePortByID mocks base method
func (m *MockNetworking) DeletePortByID(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePortByID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePortByID indicates an expected call of DeletePortByID
func (mr *MockNetworkingMockRecorder) DeletePortByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePortByID", reflect.TypeOf((*MockNetworking)(nil).DeletePortByID), arg0)
}

// FindNetworks mocks base method
func (m *MockNetworking) FindNetworks(arg0 bool) (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindNetworks", arg0)
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindNetworks indicates an expected call of FindNetworks
func (mr *MockNetworkingMockRecorder) FindNetworks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindNetworks", reflect.TypeOf((*MockNetworking)(nil).FindNetworks), arg0)
}

// NetworkInterfaces mocks base method
func (m *MockNetworking) NetworkInterfaces(arg0 []instance.Id) ([]network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkInterfaces", arg0)
	ret0, _ := ret[0].([]network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkInterfaces indicates an expected call of NetworkInterfaces
func (mr *MockNetworkingMockRecorder) NetworkInterfaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkInterfaces", reflect.TypeOf((*MockNetworking)(nil).NetworkInterfaces), arg0)
}

// ResolveNetwork mocks base method
func (m *MockNetworking) ResolveNetwork(arg0 string, arg1 bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveNetwork", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveNetwork indicates an expected call of ResolveNetwork
func (mr *MockNetworkingMockRecorder) ResolveNetwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveNetwork", reflect.TypeOf((*MockNetworking)(nil).ResolveNetwork), arg0, arg1)
}

// Subnets mocks base method
func (m *MockNetworking) Subnets(arg0 instance.Id, arg1 []network.Id) ([]network.SubnetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subnets", arg0, arg1)
	ret0, _ := ret[0].([]network.SubnetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subnets indicates an expected call of Subnets
func (mr *MockNetworkingMockRecorder) Subnets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subnets", reflect.TypeOf((*MockNetworking)(nil).Subnets), arg0, arg1)
}
