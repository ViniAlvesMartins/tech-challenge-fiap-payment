package use_case

import (
	"errors"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	orderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/order"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestOrderUseCase_GetById(t *testing.T) {
	t.Run("get order by id successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		clientID := 1

		order := &entity.Order{
			ID:          1,
			ClientId:    &clientID,
			StatusOrder: "PENDING",
			Amount:      123.4,
			CreatedAt:   time.Now(),
		}

		orderServiceResponse := &orderservice.GetByIdResp{
			Error: "",
			Data:  order,
		}

		orderService := mock.NewMockOrderService(ctrl)
		orderService.EXPECT().GetById(1).Return(orderServiceResponse, nil).Times(1)

		orderUseCase := NewOrderUseCase(orderService, logger)
		orders, err := orderUseCase.GetById(1)

		assert.Nil(t, err)
		assert.Equal(t, orders, order)
	})

	t.Run("error getting order by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedError := errors.New("error getting order by id")

		orderServiceResponse := &orderservice.GetByIdResp{
			Error: expectedError.Error(),
			Data:  nil,
		}

		orderService := mock.NewMockOrderService(ctrl)
		orderService.EXPECT().GetById(1).Return(orderServiceResponse, nil).Times(1)

		orderUseCase := NewOrderUseCase(orderService, logger)
		orders, err := orderUseCase.GetById(1)

		assert.Errorf(t, expectedError, err.Error())
		assert.Nil(t, orders)
	})
}
