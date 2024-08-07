// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	entity "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	mock "github.com/stretchr/testify/mock"
)

// PaymentInterface is an autogenerated mock type for the PaymentInterface type
type PaymentInterface[T interface{}] struct {
	mock.Mock
}

// Process provides a mock function with given fields: p
func (_m *PaymentInterface[T]) Process(p entity.Payment) (*T, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.Payment) (*T, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func(entity.Payment) *T); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(entity.Payment) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPaymentInterface creates a new instance of PaymentInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentInterface[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentInterface[T] {
	mock := &PaymentInterface[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
