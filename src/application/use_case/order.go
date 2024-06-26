package use_case

import (
	"errors"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"log/slog"
)

type OrderUseCase struct {
	orderService contract.OrderService
	logger       *slog.Logger
}

func NewOrderUseCase(orderService contract.OrderService, logger *slog.Logger) *OrderUseCase {
	return &OrderUseCase{
		orderService: orderService,
		logger:       logger,
	}
}

func (o *OrderUseCase) GetById(id int) (*entity.Order, error) {
	orderServiceResponse, err := o.orderService.GetById(id)
	if err != nil {
		return nil, err
	}

	if orderServiceResponse.Error != "" {
		return nil, errors.New(orderServiceResponse.Error)
	}

	if orderServiceResponse.Data == nil {
		return nil, nil
	}

	return &entity.Order{
		ID:          orderServiceResponse.Data.ID,
		ClientId:    orderServiceResponse.Data.ClientId,
		StatusOrder: orderServiceResponse.Data.StatusOrder,
		Amount:      orderServiceResponse.Data.Amount,
		CreatedAt:   orderServiceResponse.Data.CreatedAt,
	}, nil
}
