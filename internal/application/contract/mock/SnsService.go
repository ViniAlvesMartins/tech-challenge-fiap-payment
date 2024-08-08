// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	entity "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"

	mock "github.com/stretchr/testify/mock"
)

// SnsService is an autogenerated mock type for the SnsService type
type SnsService struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: ctx, message
func (_m *SnsService) SendMessage(ctx context.Context, message entity.PaymentMessage) error {
	ret := _m.Called(ctx, message)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.PaymentMessage) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSnsService creates a new instance of SnsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSnsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *SnsService {
	mock := &SnsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}