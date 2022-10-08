// Code generated by MockGen. DO NOT EDIT.
// Source: ethereum.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	ethereum "github.com/longfan78/quorum-key-manager/pkg/ethereum"
	entities "github.com/longfan78/quorum-key-manager/src/stores/entities"
	types "github.com/longfan78/quorum/core/types"
	common "github.com/ethereum/go-ethereum/common"
	types0 "github.com/ethereum/go-ethereum/core/types"
	core "github.com/ethereum/go-ethereum/signer/core"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockEthStore is a mock of EthStore interface
type MockEthStore struct {
	ctrl     *gomock.Controller
	recorder *MockEthStoreMockRecorder
}

// MockEthStoreMockRecorder is the mock recorder for MockEthStore
type MockEthStoreMockRecorder struct {
	mock *MockEthStore
}

// NewMockEthStore creates a new mock instance
func NewMockEthStore(ctrl *gomock.Controller) *MockEthStore {
	mock := &MockEthStore{ctrl: ctrl}
	mock.recorder = &MockEthStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEthStore) EXPECT() *MockEthStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockEthStore) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, id, attr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockEthStoreMockRecorder) Create(ctx, id, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockEthStore)(nil).Create), ctx, id, attr)
}

// Import mocks base method
func (m *MockEthStore) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", ctx, id, privKey, attr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Import indicates an expected call of Import
func (mr *MockEthStoreMockRecorder) Import(ctx, id, privKey, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockEthStore)(nil).Import), ctx, id, privKey, attr)
}

// Get mocks base method
func (m *MockEthStore) Get(ctx context.Context, addr common.Address) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, addr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockEthStoreMockRecorder) Get(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockEthStore)(nil).Get), ctx, addr)
}

// List mocks base method
func (m *MockEthStore) List(ctx context.Context, limit, offset uint64) ([]common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset)
	ret0, _ := ret[0].([]common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockEthStoreMockRecorder) List(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockEthStore)(nil).List), ctx, limit, offset)
}

// Update mocks base method
func (m *MockEthStore) Update(ctx context.Context, addr common.Address, attr *entities.Attributes) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, addr, attr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockEthStoreMockRecorder) Update(ctx, addr, attr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockEthStore)(nil).Update), ctx, addr, attr)
}

// Delete mocks base method
func (m *MockEthStore) Delete(ctx context.Context, addr common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockEthStoreMockRecorder) Delete(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockEthStore)(nil).Delete), ctx, addr)
}

// GetDeleted mocks base method
func (m *MockEthStore) GetDeleted(ctx context.Context, addr common.Address) (*entities.ETHAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeleted", ctx, addr)
	ret0, _ := ret[0].(*entities.ETHAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeleted indicates an expected call of GetDeleted
func (mr *MockEthStoreMockRecorder) GetDeleted(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeleted", reflect.TypeOf((*MockEthStore)(nil).GetDeleted), ctx, addr)
}

// ListDeleted mocks base method
func (m *MockEthStore) ListDeleted(ctx context.Context, limit, offset uint64) ([]common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListDeleted", ctx, limit, offset)
	ret0, _ := ret[0].([]common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListDeleted indicates an expected call of ListDeleted
func (mr *MockEthStoreMockRecorder) ListDeleted(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListDeleted", reflect.TypeOf((*MockEthStore)(nil).ListDeleted), ctx, limit, offset)
}

// Restore mocks base method
func (m *MockEthStore) Restore(ctx context.Context, addr common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore
func (mr *MockEthStoreMockRecorder) Restore(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockEthStore)(nil).Restore), ctx, addr)
}

// Destroy mocks base method
func (m *MockEthStore) Destroy(ctx context.Context, addr common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", ctx, addr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockEthStoreMockRecorder) Destroy(ctx, addr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockEthStore)(nil).Destroy), ctx, addr)
}

// Sign mocks base method
func (m *MockEthStore) Sign(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sign", ctx, addr, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sign indicates an expected call of Sign
func (mr *MockEthStoreMockRecorder) Sign(ctx, addr, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sign", reflect.TypeOf((*MockEthStore)(nil).Sign), ctx, addr, data)
}

// SignMessage mocks base method
func (m *MockEthStore) SignMessage(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignMessage", ctx, addr, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignMessage indicates an expected call of SignMessage
func (mr *MockEthStoreMockRecorder) SignMessage(ctx, addr, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignMessage", reflect.TypeOf((*MockEthStore)(nil).SignMessage), ctx, addr, data)
}

// SignTypedData mocks base method
func (m *MockEthStore) SignTypedData(ctx context.Context, addr common.Address, typedData *core.TypedData) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTypedData", ctx, addr, typedData)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTypedData indicates an expected call of SignTypedData
func (mr *MockEthStoreMockRecorder) SignTypedData(ctx, addr, typedData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTypedData", reflect.TypeOf((*MockEthStore)(nil).SignTypedData), ctx, addr, typedData)
}

// SignTransaction mocks base method
func (m *MockEthStore) SignTransaction(ctx context.Context, addr common.Address, chainID *big.Int, tx *types0.Transaction) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTransaction", ctx, addr, chainID, tx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTransaction indicates an expected call of SignTransaction
func (mr *MockEthStoreMockRecorder) SignTransaction(ctx, addr, chainID, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTransaction", reflect.TypeOf((*MockEthStore)(nil).SignTransaction), ctx, addr, chainID, tx)
}

// SignEEA mocks base method
func (m *MockEthStore) SignEEA(ctx context.Context, addr common.Address, chainID *big.Int, tx *types0.Transaction, args *ethereum.PrivateArgs) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignEEA", ctx, addr, chainID, tx, args)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignEEA indicates an expected call of SignEEA
func (mr *MockEthStoreMockRecorder) SignEEA(ctx, addr, chainID, tx, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignEEA", reflect.TypeOf((*MockEthStore)(nil).SignEEA), ctx, addr, chainID, tx, args)
}

// SignPrivate mocks base method
func (m *MockEthStore) SignPrivate(ctx context.Context, addr common.Address, tx *types.Transaction) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignPrivate", ctx, addr, tx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignPrivate indicates an expected call of SignPrivate
func (mr *MockEthStoreMockRecorder) SignPrivate(ctx, addr, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignPrivate", reflect.TypeOf((*MockEthStore)(nil).SignPrivate), ctx, addr, tx)
}

// Encrypt mocks base method
func (m *MockEthStore) Encrypt(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", ctx, addr, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt
func (mr *MockEthStoreMockRecorder) Encrypt(ctx, addr, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*MockEthStore)(nil).Encrypt), ctx, addr, data)
}

// Decrypt mocks base method
func (m *MockEthStore) Decrypt(ctx context.Context, addr common.Address, data []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrypt", ctx, addr, data)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrypt indicates an expected call of Decrypt
func (mr *MockEthStoreMockRecorder) Decrypt(ctx, addr, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrypt", reflect.TypeOf((*MockEthStore)(nil).Decrypt), ctx, addr, data)
}
