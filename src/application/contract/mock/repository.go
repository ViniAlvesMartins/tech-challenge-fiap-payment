// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	entity "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockPaymentRepository is a mock of PaymentRepository interface.
type MockPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryMockRecorder
}

// MockPaymentRepositoryMockRecorder is the mock recorder for MockPaymentRepository.
type MockPaymentRepositoryMockRecorder struct {
	mock *MockPaymentRepository
}

// NewMockPaymentRepository creates a new mock instance.
func NewMockPaymentRepository(ctrl *gomock.Controller) *MockPaymentRepository {
	mock := &MockPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentRepository) EXPECT() *MockPaymentRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPaymentRepository) Create(payment entity.Payment) (*entity.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", payment)
	ret0, _ := ret[0].(*entity.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPaymentRepositoryMockRecorder) Create(payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPaymentRepository)(nil).Create), payment)
}

// GetLastPaymentStatus mocks base method.
func (m *MockPaymentRepository) GetLastPaymentStatus(orderId int) (*entity.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastPaymentStatus", orderId)
	ret0, _ := ret[0].(*entity.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastPaymentStatus indicates an expected call of GetLastPaymentStatus.
func (mr *MockPaymentRepositoryMockRecorder) GetLastPaymentStatus(orderId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastPaymentStatus", reflect.TypeOf((*MockPaymentRepository)(nil).GetLastPaymentStatus), orderId)
}