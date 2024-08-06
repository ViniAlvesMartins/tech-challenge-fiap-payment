package use_case

import (
	"errors"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	orderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/order"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestOrderUseCase_GetById(t *testing.T) {
	t.Run("get order by id successfully", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		clientID := 1

		order := &entity.Order{
			ID:          1,
			ClientId:    &clientID,
			OrderStatus: "PENDING",
			Amount:      123.4,
			CreatedAt:   time.Now(),
		}

		orderServiceResponse := &orderservice.GetByIdResp{
			Error: "",
			Data:  order,
		}

		orderService := new(mock.OrderService)
		orderService.On("GetById", 1).Return(orderServiceResponse, nil).Once()

		orderUseCase := NewOrderUseCase(orderService, logger)
		orders, err := orderUseCase.GetById(1)

		assert.Nil(t, err)
		assert.Equal(t, orders, order)
	})

	t.Run("error getting order by id", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedError := errors.New("error getting order by id")

		orderServiceResponse := &orderservice.GetByIdResp{
			Error: expectedError.Error(),
			Data:  nil,
		}

		orderService := new(mock.OrderService)
		orderService.On("GetById", 1).Return(orderServiceResponse, nil).Once()

		orderUseCase := NewOrderUseCase(orderService, logger)
		orders, err := orderUseCase.GetById(1)

		assert.Errorf(t, expectedError, err.Error())
		assert.Nil(t, orders)
	})
}
