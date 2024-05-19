package use_case

import (
	"context"
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract/mock"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestPaymentUseCase_Create(t *testing.T) {
	t.Run("create payment successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		ctx := context.Background()

		payment := entity.Payment{
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		result := payment
		result.PaymentID = "65cf595b-19b9-431b-9a81-9818dec845ff"

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().Create(ctx, payment).Times(1).Return(&result, nil)
		sns := mock.NewMockSnsService(ctrl)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.Create(ctx, &payment)

		assert.Nil(t, err)
	})

	t.Run("error creating payment", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error connecting to database")

		payment := entity.Payment{
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().Create(ctx, payment).Times(1).Return(nil, expectedErr)
		sns := mock.NewMockSnsService(ctrl)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.Create(ctx, &payment)

		assert.Error(t, expectedErr, err)
	})
}

func TestPaymentUseCase_GetLastPaymentStatus(t *testing.T) {
	for _, tt := range lastPaymentStatusPayments() {
		t.Run(fmt.Sprintf("get last payment status successfully:%s", tt.Type), func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

			sns := mock.NewMockSnsService(ctrl)
			repo := mock.NewMockPaymentRepository(ctrl)
			repo.EXPECT().GetLastPaymentStatus(ctx, 1).Times(1).Return(tt.Payment, nil)

			paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
			status, err := paymentUseCase.GetLastPaymentStatus(ctx, 1)

			assert.Equal(t, tt.Expected, status)
			assert.Nil(t, err)
		})
	}

	t.Run("error getting last payment status", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error connecting to database")

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       132.45,
		}

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		sns := mock.NewMockSnsService(ctrl)
		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().GetLastPaymentStatus(ctx, 1).Times(1).Return(payment, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		status, err := paymentUseCase.GetLastPaymentStatus(ctx, 1)

		assert.Equal(t, payment.CurrentState, status)
		assert.Error(t, expectedErr, err)
	})
}

type paymentTest struct {
	Type     string
	Payment  *entity.Payment
	Expected enum.PaymentStatus
}

func TestPaymentUseCase_CreateQRCode(t *testing.T) {
	t.Run("create qr code successfully", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.PREPARING,
			Amount:      123.45,
		}

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      order.ID,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       order.Amount,
		}

		qrCode := responsepaymentservice.CreateQRCode{QrData: "qr data"}

		sns := mock.NewMockSnsService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		getLastPaymentStatus := repo.EXPECT().GetLastPaymentStatus(ctx, 1).Times(1).Return(payment, nil)
		createPayment := repo.EXPECT().Create(ctx, gomock.Any()).Times(1).Return(payment, nil).After(getLastPaymentStatus)

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)
		externalPaymentMock.EXPECT().CreateQRCode(gomock.Any()).Times(1).Return(qrCode, nil).After(createPayment)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Equal(t, code.QrData, qrCode.QrData)
		assert.Nil(t, err)
	})

	t.Run("last status confirmed", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.PREPARING,
			Amount:      123.45,
		}

		payment := &entity.Payment{
			PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID:      order.ID,
			Type:         enum.QRCODE,
			CurrentState: enum.CONFIRMED,
			Amount:       order.Amount,
		}

		sns := mock.NewMockSnsService(ctrl)
		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().GetLastPaymentStatus(ctx, 1).Times(1).Return(payment, nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Nil(t, code)
		assert.Nil(t, err)
	})

	t.Run("get last payment status error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error getting status")

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.PREPARING,
			Amount:      123.45,
		}

		sns := mock.NewMockSnsService(ctrl)
		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().GetLastPaymentStatus(ctx, 1).Times(1).Return(nil, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		code, err := paymentUseCase.CreateQRCode(ctx, order)

		assert.Nil(t, code)
		assert.Error(t, expectedErr, err)
	})
}

func TestPaymentUseCase_PaymentNotification(t *testing.T) {
	t.Run("send notification successfully", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)
		repo := mock.NewMockPaymentRepository(ctrl)
		updateStatus := repo.EXPECT().UpdateStatus(ctx, 1, enum.CONFIRMED).Times(1).Return(nil)

		sns := mock.NewMockSnsService(ctrl)
		sns.EXPECT().SendMessage(1, enum.CONFIRMED).Times(1).Return(nil).After(updateStatus)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.PaymentNotification(ctx, 1)

		assert.Nil(t, err)
	})

	t.Run("error updating status", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error updating status")

		sns := mock.NewMockSnsService(ctrl)
		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().UpdateStatus(ctx, 1, enum.CONFIRMED).Times(1).Return(expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.PaymentNotification(ctx, 1)

		assert.Error(t, expectedErr, err)
	})

	t.Run("error sending message", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error sending message")

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)
		repo := mock.NewMockPaymentRepository(ctrl)
		updateStatus := repo.EXPECT().UpdateStatus(ctx, 1, enum.CONFIRMED).Times(1).Return(nil)

		sns := mock.NewMockSnsService(ctrl)
		sns.EXPECT().SendMessage(1, enum.CONFIRMED).Times(1).Return(expectedErr).After(updateStatus)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, sns, logger)
		err := paymentUseCase.PaymentNotification(ctx, 1)

		assert.Error(t, expectedErr, err)
	})
}

func lastPaymentStatusPayments() []paymentTest {
	return []paymentTest{
		{
			Payment: &entity.Payment{
				PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f0",
				OrderID:      1,
				Type:         enum.QRCODE,
				CurrentState: enum.PENDING,
				Amount:       132.45,
			},
			Expected: enum.PENDING,
			Type:     fmt.Sprintf("[%s status]", string(enum.PENDING)),
		},
		{
			Payment: &entity.Payment{
				PaymentID:    "65cf595b-19b9-431b-9a81-9818dec845f1",
				OrderID:      1,
				Type:         enum.QRCODE,
				CurrentState: "",
				Amount:       132.45,
			},
			Expected: enum.PENDING,
			Type:     "[empty status]",
		},
	}
}
