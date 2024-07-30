package use_case

import (
	"context"
	"errors"
	"fmt"
	contractmock "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"testing"
)

func TestPaymentUseCase_GetLastPaymentStatus(t *testing.T) {
	for _, tt := range lastPaymentStatusPayments() {
		t.Run(fmt.Sprintf("get last payment status successfully:%s", tt.Type), func(t *testing.T) {
			ctx := context.Background()
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			externalPaymentMock := new(contractmock.PaymentInterface[entity.QRCodePayment])
			sns := contractmock.NewSnsService(t)
			repo := contractmock.NewPaymentRepository(t)
			repo.On("GetLastPaymentStatus", ctx, 1).Return(tt.Payment, nil).Once()

			paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
			status, err := paymentUseCase.GetLastPaymentStatus(ctx, 1)

			assert.Equal(t, tt.Expected, status)
			assert.Nil(t, err)
		})
	}

	t.Run("error getting last payment status", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error connecting to database")

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      1,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusPending,
			Amount:       132.45,
		}

		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)
		sns := contractmock.NewSnsService(t)
		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, 1).Once().Return(payment, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		status, err := paymentUseCase.GetLastPaymentStatus(ctx, 1)

		assert.Equal(t, payment.CurrentState, status)
		assert.Error(t, expectedErr, err)
	})
}

func TestPaymentUseCase_CreateQRCode(t *testing.T) {
	t.Run("create qr code successfully", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusPreparing,
			Amount:      123.45,
		}

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      order.ID,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusPending,
			Amount:       order.Amount,
		}

		qrCode := &entity.QRCodePayment{QRCode: "qr data"}

		sns := contractmock.NewSnsService(t)
		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, 1).Once().Return(payment, nil)
		repo.On("Create", ctx, mock.AnythingOfType("entity.Payment")).Once().Return(nil)

		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)
		externalPaymentMock.On("Process", mock.AnythingOfType("entity.Payment")).Once().Return(qrCode, nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Equal(t, code.QRCode, qrCode.QRCode)
		assert.Nil(t, err)
	})

	t.Run("last status confirmed", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusPreparing,
			Amount:      123.45,
		}

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      order.ID,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusConfirmed,
			Amount:       order.Amount,
		}

		sns := contractmock.NewSnsService(t)
		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)

		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, 1).Once().Return(payment, nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Nil(t, code)
		assert.Nil(t, err)
	})

	t.Run("get last payment status error", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error getting status")

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusPreparing,
			Amount:      123.45,
		}

		sns := contractmock.NewSnsService(t)
		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)

		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, 1).Once().Return(nil, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Nil(t, code)
		assert.Error(t, expectedErr, err)
	})
}

func TestPaymentUseCase_PaymentNotification(t *testing.T) {
	t.Run("send notification successfully", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		payment := entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      1,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusPending,
			Amount:       123.45,
		}

		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)
		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, payment.OrderID).Once().Return(&payment, nil)

		repo.On("UpdateStatus", ctx, payment.OrderID, enum.PaymentStatusConfirmed).Once().Return(nil)

		sns := contractmock.NewSnsService(t)
		sns.On("SendMessage", ctx, entity.PaymentMessage{
			OrderId: 1,
			Status:  enum.PaymentStatusConfirmed,
		}).Once().Return(nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.ConfirmedPaymentNotification(ctx, payment.OrderID)

		assert.Nil(t, err)
	})

	t.Run("error updating status", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error updating status")

		payment := entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      1,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusPending,
			Amount:       123.45,
		}

		sns := contractmock.NewSnsService(t)
		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)

		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, payment.OrderID).Once().Return(&payment, nil)
		repo.On("UpdateStatus", ctx, payment.OrderID, enum.PaymentStatusConfirmed).Once().Return(expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.ConfirmedPaymentNotification(ctx, payment.OrderID)

		assert.Error(t, expectedErr, err)
	})

	t.Run("error sending message", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error sending message")

		payment := entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      1,
			Type:         enum.PaymentTypeQRCode,
			CurrentState: enum.PaymentStatusPending,
			Amount:       123.45,
		}

		externalPaymentMock := contractmock.NewPaymentInterface[entity.QRCodePayment](t)
		repo := contractmock.NewPaymentRepository(t)
		repo.On("GetLastPaymentStatus", ctx, payment.OrderID).Once().Return(&payment, nil)
		repo.On("UpdateStatus", ctx, 1, enum.PaymentStatusConfirmed).Once().Return(nil)

		sns := contractmock.NewSnsService(t)
		sns.On("SendMessage", ctx, entity.PaymentMessage{
			OrderId: 1,
			Status:  enum.PaymentStatusConfirmed,
		}).Once().Return(expectedErr)

		repo.On("UpdateStatus", ctx, 1, enum.PaymentStatusPending).Once().Return(nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.ConfirmedPaymentNotification(ctx, 1)

		assert.Error(t, expectedErr, err)
	})
}

type paymentTest struct {
	Type     string
	Payment  *entity.Payment
	Expected enum.PaymentStatus
}

func lastPaymentStatusPayments() []paymentTest {
	return []paymentTest{
		{
			Payment: &entity.Payment{
				PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
				OrderID:      1,
				Type:         enum.PaymentTypeQRCode,
				CurrentState: enum.PaymentStatusPending,
				Amount:       132.45,
			},
			Expected: enum.PaymentStatusPending,
			Type:     fmt.Sprintf("[%s status]", string(enum.PaymentStatusPending)),
		},
		{
			Payment:  nil,
			Expected: enum.PaymentStatusPending,
			Type:     "[empty status]",
		},
	}
}
