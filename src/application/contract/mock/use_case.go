// Code generated by MockGen. DO NOT EDIT.
// Source: use_case.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	entity "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	enum "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderUseCase is a mock of OrderUseCase interface.
type MockOrderUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUseCaseMockRecorder
}

// MockOrderUseCaseMockRecorder is the mock recorder for MockOrderUseCase.
type MockOrderUseCaseMockRecorder struct {
	mock *MockOrderUseCase
}

// NewMockOrderUseCase creates a new mock instance.
func NewMockOrderUseCase(ctrl *gomock.Controller) *MockOrderUseCase {
	mock := &MockOrderUseCase{ctrl: ctrl}
	mock.recorder = &MockOrderUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUseCase) EXPECT() *MockOrderUseCaseMockRecorder {
	return m.recorder
}

// GetById mocks base method.
func (m *MockOrderUseCase) GetById(id int) (*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockOrderUseCaseMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockOrderUseCase)(nil).GetById), id)
}

// MockPaymentUseCase is a mock of PaymentUseCase interface.
type MockPaymentUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentUseCaseMockRecorder
}

// MockPaymentUseCaseMockRecorder is the mock recorder for MockPaymentUseCase.
type MockPaymentUseCaseMockRecorder struct {
	mock *MockPaymentUseCase
}

// NewMockPaymentUseCase creates a new mock instance.
func NewMockPaymentUseCase(ctrl *gomock.Controller) *MockPaymentUseCase {
	mock := &MockPaymentUseCase{ctrl: ctrl}
	mock.recorder = &MockPaymentUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentUseCase) EXPECT() *MockPaymentUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPaymentUseCase) Create(ctx context.Context, payment *entity.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, payment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPaymentUseCaseMockRecorder) Create(ctx, payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPaymentUseCase)(nil).Create), ctx, payment)
}

// CreateQRCode mocks base method.
func (m *MockPaymentUseCase) CreateQRCode(ctx context.Context, order *entity.Order) (*response_payment_service.CreateQRCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQRCode", ctx, order)
	ret0, _ := ret[0].(*response_payment_service.CreateQRCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQRCode indicates an expected call of CreateQRCode.
func (mr *MockPaymentUseCaseMockRecorder) CreateQRCode(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQRCode", reflect.TypeOf((*MockPaymentUseCase)(nil).CreateQRCode), ctx, order)
}

// GetLastPaymentStatus mocks base method.
func (m *MockPaymentUseCase) GetLastPaymentStatus(ctx context.Context, paymentId int) (enum.PaymentStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastPaymentStatus", ctx, paymentId)
	ret0, _ := ret[0].(enum.PaymentStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastPaymentStatus indicates an expected call of GetLastPaymentStatus.
func (mr *MockPaymentUseCaseMockRecorder) GetLastPaymentStatus(ctx, paymentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastPaymentStatus", reflect.TypeOf((*MockPaymentUseCase)(nil).GetLastPaymentStatus), ctx, paymentId)
}

// PaymentNotification mocks base method.
func (m *MockPaymentUseCase) PaymentNotification(ctx context.Context, paymentId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaymentNotification", ctx, paymentId)
	ret0, _ := ret[0].(error)
	return ret0
}

// PaymentNotification indicates an expected call of PaymentNotification.
func (mr *MockPaymentUseCaseMockRecorder) PaymentNotification(ctx, paymentId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaymentNotification", reflect.TypeOf((*MockPaymentUseCase)(nil).PaymentNotification), ctx, paymentId)
}
