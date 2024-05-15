package use_case

import (
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract/mock"
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

		payment := entity.Payment{
			OrderID: 1,
			Type:    enum.QRCODE,
			Status:  enum.PENDING,
			Amount:  123.45,
		}

		result := payment
		result.ID = "65cf595b-19b9-431b-9a81-9818dec845ff"

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().Create(payment).Times(1).Return(&result, nil)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, logger)
		err := paymentUseCase.Create(&payment)

		assert.Nil(t, err)
	})

	t.Run("error creating payment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error connecting to database")

		payment := entity.Payment{
			OrderID: 1,
			Type:    enum.QRCODE,
			Status:  enum.PENDING,
			Amount:  123.45,
		}

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().Create(payment).Times(1).Return(nil, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, logger)
		err := paymentUseCase.Create(&payment)

		assert.Error(t, expectedErr, err)
	})
}

func TestPaymentUseCase_GetLastPaymentStatus(t *testing.T) {
	for _, tt := range lastPaymentStatusPayments() {
		t.Run(fmt.Sprintf("get last payment status successfully:%s", tt.Type), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

			repo := mock.NewMockPaymentRepository(ctrl)
			repo.EXPECT().GetLastPaymentStatus(1).Times(1).Return(tt.Payment, nil)

			paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, logger)
			status, err := paymentUseCase.GetLastPaymentStatus(1)

			assert.Equal(t, tt.Expected, status)
			assert.Nil(t, err)
		})
	}

	t.Run("error getting last payment status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("error connecting to database")

		payment := &entity.Payment{
			ID:      "65cf595b-19b9-431b-9a81-9818dec845f0",
			OrderID: 1,
			Type:    enum.QRCODE,
			Status:  enum.PENDING,
			Amount:  132.45,
		}

		externalPaymentMock := mock.NewMockExternalPaymentService(ctrl)

		repo := mock.NewMockPaymentRepository(ctrl)
		repo.EXPECT().GetLastPaymentStatus(1).Times(1).Return(payment, expectedErr)

		paymentUseCase := NewPaymentUseCase(repo, externalPaymentMock, logger)
		status, err := paymentUseCase.GetLastPaymentStatus(1)

		assert.Equal(t, payment.Status, status)
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
				ID:      "65cf595b-19b9-431b-9a81-9818dec845f0",
				OrderID: 1,
				Type:    enum.QRCODE,
				Status:  enum.PENDING,
				Amount:  132.45,
			},
			Expected: enum.PENDING,
			Type:     fmt.Sprintf("[%s status]", string(enum.PENDING)),
		},
		{
			Payment: &entity.Payment{
				ID:      "65cf595b-19b9-431b-9a81-9818dec845f1",
				OrderID: 1,
				Type:    enum.QRCODE,
				Status:  "",
				Amount:  132.45,
			},
			Expected: enum.PENDING,
			Type:     "[empty status]",
		},
	}
}
