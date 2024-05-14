package service

import (
	"fmt"
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/order_service"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type OrderService struct {
	client *resty.Client
	logger *slog.Logger
}

func NewOrderService(c *resty.Client, l *slog.Logger) *OrderService {
	return &OrderService{logger: l, client: c}
}

func (o *OrderService) GetById(id int) (*responseorderservice.GetByIdResp, error) {
	var order responseorderservice.GetByIdResp

	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(order).
		Get(fmt.Sprintf("/orders/%d", id))

	if err != nil {
		o.logger.Error("error making request", slog.String("error", err.Error()))
		return nil, err
	}

	if resp.IsError() {
		o.logger.Error("request response error", slog.Int("error_status_code", resp.StatusCode()))
		return nil, err
	}

	return &order, nil
}
