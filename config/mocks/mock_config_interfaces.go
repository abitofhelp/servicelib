// Code generated by MockGen. DO NOT EDIT.
// Source: ../interfaces/config.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	interfaces "github.com/abitofhelp/servicelib/config/interfaces"
	gomock "github.com/golang/mock/gomock"
)

// MockAppConfig is a mock of AppConfig interface.
type MockAppConfig struct {
	ctrl     *gomock.Controller
	recorder *MockAppConfigMockRecorder
}

// MockAppConfigMockRecorder is the mock recorder for MockAppConfig.
type MockAppConfigMockRecorder struct {
	mock *MockAppConfig
}

// NewMockAppConfig creates a new mock instance.
func NewMockAppConfig(ctrl *gomock.Controller) *MockAppConfig {
	mock := &MockAppConfig{ctrl: ctrl}
	mock.recorder = &MockAppConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppConfig) EXPECT() *MockAppConfigMockRecorder {
	return m.recorder
}

// GetEnvironment mocks base method.
func (m *MockAppConfig) GetEnvironment() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnvironment")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetEnvironment indicates an expected call of GetEnvironment.
func (mr *MockAppConfigMockRecorder) GetEnvironment() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnvironment", reflect.TypeOf((*MockAppConfig)(nil).GetEnvironment))
}

// GetName mocks base method.
func (m *MockAppConfig) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName.
func (mr *MockAppConfigMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockAppConfig)(nil).GetName))
}

// GetVersion mocks base method.
func (m *MockAppConfig) GetVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockAppConfigMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockAppConfig)(nil).GetVersion))
}

// MockDatabaseConfig is a mock of DatabaseConfig interface.
type MockDatabaseConfig struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseConfigMockRecorder
}

// MockDatabaseConfigMockRecorder is the mock recorder for MockDatabaseConfig.
type MockDatabaseConfigMockRecorder struct {
	mock *MockDatabaseConfig
}

// NewMockDatabaseConfig creates a new mock instance.
func NewMockDatabaseConfig(ctrl *gomock.Controller) *MockDatabaseConfig {
	mock := &MockDatabaseConfig{ctrl: ctrl}
	mock.recorder = &MockDatabaseConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseConfig) EXPECT() *MockDatabaseConfigMockRecorder {
	return m.recorder
}

// GetCollectionName mocks base method.
func (m *MockDatabaseConfig) GetCollectionName(entityType string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionName", entityType)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCollectionName indicates an expected call of GetCollectionName.
func (mr *MockDatabaseConfigMockRecorder) GetCollectionName(entityType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionName", reflect.TypeOf((*MockDatabaseConfig)(nil).GetCollectionName), entityType)
}

// GetConnectionString mocks base method.
func (m *MockDatabaseConfig) GetConnectionString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionString")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetConnectionString indicates an expected call of GetConnectionString.
func (mr *MockDatabaseConfigMockRecorder) GetConnectionString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionString", reflect.TypeOf((*MockDatabaseConfig)(nil).GetConnectionString))
}

// GetDatabaseName mocks base method.
func (m *MockDatabaseConfig) GetDatabaseName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabaseName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDatabaseName indicates an expected call of GetDatabaseName.
func (mr *MockDatabaseConfigMockRecorder) GetDatabaseName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabaseName", reflect.TypeOf((*MockDatabaseConfig)(nil).GetDatabaseName))
}

// GetType mocks base method.
func (m *MockDatabaseConfig) GetType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockDatabaseConfigMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockDatabaseConfig)(nil).GetType))
}

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// GetApp mocks base method.
func (m *MockConfig) GetApp() interfaces.AppConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApp")
	ret0, _ := ret[0].(interfaces.AppConfig)
	return ret0
}

// GetApp indicates an expected call of GetApp.
func (mr *MockConfigMockRecorder) GetApp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApp", reflect.TypeOf((*MockConfig)(nil).GetApp))
}

// GetDatabase mocks base method.
func (m *MockConfig) GetDatabase() interfaces.DatabaseConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatabase")
	ret0, _ := ret[0].(interfaces.DatabaseConfig)
	return ret0
}

// GetDatabase indicates an expected call of GetDatabase.
func (mr *MockConfigMockRecorder) GetDatabase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabase", reflect.TypeOf((*MockConfig)(nil).GetDatabase))
}
