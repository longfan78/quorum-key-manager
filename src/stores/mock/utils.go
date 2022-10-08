// Code generated by MockGen. DO NOT EDIT.
// Source: utils.go

// Package mock is a generated GoMock package.
package mock

import (
	"github.com/longfan78/quorum-key-manager/src/entities"
	common "github.com/ethereum/go-ethereum/common"
	core "github.com/ethereum/go-ethereum/signer/core"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUtils is a mock of Utils interface
type MockUtils struct {
	ctrl     *gomock.Controller
	recorder *MockUtilsMockRecorder
}

// MockUtilsMockRecorder is the mock recorder for MockUtils
type MockUtilsMockRecorder struct {
	mock *MockUtils
}

// NewMockUtils creates a new mock instance
func NewMockUtils(ctrl *gomock.Controller) *MockUtils {
	mock := &MockUtils{ctrl: ctrl}
	mock.recorder = &MockUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUtils) EXPECT() *MockUtilsMockRecorder {
	return m.recorder
}

// Verify mocks base method
func (m *MockUtils) Verify(pubKey, data, sig []byte, algo *entities.Algorithm) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", pubKey, data, sig, algo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify
func (mr *MockUtilsMockRecorder) Verify(pubKey, data, sig, algo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockUtils)(nil).Verify), pubKey, data, sig, algo)
}

// ECRecover mocks base method
func (m *MockUtils) ECRecover(data, sig []byte) (common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ECRecover", data, sig)
	ret0, _ := ret[0].(common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ECRecover indicates an expected call of ECRecover
func (mr *MockUtilsMockRecorder) ECRecover(data, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ECRecover", reflect.TypeOf((*MockUtils)(nil).ECRecover), data, sig)
}

// VerifyMessage mocks base method
func (m *MockUtils) VerifyMessage(addr common.Address, data, sig []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyMessage", addr, data, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyMessage indicates an expected call of VerifyMessage
func (mr *MockUtilsMockRecorder) VerifyMessage(addr, data, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyMessage", reflect.TypeOf((*MockUtils)(nil).VerifyMessage), addr, data, sig)
}

// VerifyTypedData mocks base method
func (m *MockUtils) VerifyTypedData(addr common.Address, typedData *core.TypedData, sig []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTypedData", addr, typedData, sig)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTypedData indicates an expected call of VerifyTypedData
func (mr *MockUtilsMockRecorder) VerifyTypedData(addr, typedData, sig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTypedData", reflect.TypeOf((*MockUtils)(nil).VerifyTypedData), addr, typedData, sig)
}
