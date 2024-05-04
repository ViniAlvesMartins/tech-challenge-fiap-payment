package output

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"time"
)

type OrderDto struct {
	ID          int          `json:"id"`
	ClientID    *int         `json:"client_id"`
	StatusOrder string       `json:"status_order"`
	Amount      float32      `json:"amount"`
	CreatedAt   time.Time    `json:"created_at"`
	Products    []ProductDto `json:"products"`
}

func OrderFromEntity(order entity.Order) OrderDto {
	var orderDto = OrderDto{
		ID:          order.ID,
		ClientID:    order.ClientId,
		StatusOrder: string(order.StatusOrder),
		Amount:      order.Amount,
		CreatedAt:   order.CreatedAt,
	}

	for _, p := range order.Products {
		orderDto.Products = append(orderDto.Products, ProductFromEntity(*p))
	}

	return orderDto
}

func OrderListFromEntity(orders []entity.Order) []OrderDto {
	var ordersDto []OrderDto
	for _, o := range orders {
		ordersDto = append(ordersDto, OrderFromEntity(o))
	}

	return ordersDto
}
