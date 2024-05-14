package use_case

import (
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

	if len(orderServiceResponse.Data) == 0 {
		return nil, nil
	}

	r := orderServiceResponse.Data[0]

	return &entity.Order{
		ID:          r.ID,
		ClientId:    r.ClientId,
		StatusOrder: r.StatusOrder,
		Amount:      r.Amount,
		CreatedAt:   r.CreatedAt,
	}, nil
}
