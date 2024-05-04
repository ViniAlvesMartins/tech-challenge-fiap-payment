package use_case

import (
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
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
	order, err := o.orderService.GetById(id)

	if err != nil {
		return nil, err
	}

	return &order, nil
}
