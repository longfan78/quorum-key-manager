// Code generated by MockGen. DO NOT EDIT.
// Source: secrets.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	entities "github.com/longfan78/quorum-key-manager/src/stores/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSecretStore is a mock of SecretStore interface
type MockSecretStore struct {
	ctrl     *gomock.Controller
	recorder *MockSecretStoreMockRecorder
}

// MockSecretStoreMockRecorder is the mock recorder for MockSecretStore
type MockSecretStoreMockRecorder struct {
	mock *MockSecretStore
}

// NewMockSecretStore creates a new mock instance
func NewMockSecretStore(ctrl *gomock.Controller) *MockSecretStore {
	mock := &MockSecretStore{ctrl: ctrl}
	mock.recorder = &MockSecretStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecretStore) EXPECT() *MockSecretStoreMockRecorder {
	return m.recorder
}

// Set mocks base method
func (m *MockSecretStore) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, id, value, attr)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set
func (mr *MockSecretStoreMockRecorder) Set(ctx, id, value, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockSecretStore)(nil).Set), ctx, id, value, attr)
}

// Get mocks base method
func (m *MockSecretStore) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id, version)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockSecretStoreMockRecorder) Get(ctx, id, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSecretStore)(nil).Get), ctx, id, version)
}

// List mocks base method
func (m *MockSecretStore) List(ctx context.Context, limit, offset uint64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockSecretStoreMockRecorder) List(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSecretStore)(nil).List), ctx, limit, offset)
}

// Delete mocks base method
func (m *MockSecretStore) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSecretStoreMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSecretStore)(nil).Delete), ctx, id)
}

// GetDeleted mocks base method
func (m *MockSecretStore) GetDeleted(ctx context.Context, id string) (*entities.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, id)
	ret0, _ := ret[0].(*entities.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockSecretStoreMockRecorder) GetDeleted(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockSecretStore)(nil).GetDeleted), ctx, id)
}

// ListDeleted mocks base method
func (m *MockSecretStore) ListDeleted(ctx context.Context, limit, offset uint64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeleted", ctx, limit, offset)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeleted indicates an expected call of ListDeleted
func (mr *MockSecretStoreMockRecorder) ListDeleted(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeleted", reflect.TypeOf((*MockSecretStore)(nil).ListDeleted), ctx, limit, offset)
}

// Restore mocks base method
func (m *MockSecretStore) Restore(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockSecretStoreMockRecorder) Restore(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockSecretStore)(nil).Restore), ctx, id)
}

// Destroy mocks base method
func (m *MockSecretStore) Destroy(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockSecretStoreMockRecorder) Destroy(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockSecretStore)(nil).Destroy), ctx, id)
}
