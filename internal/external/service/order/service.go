package order

import (
	"errors"
	"fmt"
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/modules/response/order_service"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type Service struct {
	client *resty.Client
	logger *slog.Logger
}

func NewService(c *resty.Client, l *slog.Logger) *Service {
	return &Service{logger: l, client: c}
}

func (o *Service) GetById(id int) (*responseorderservice.GetByIdResp, error) {
	var order responseorderservice.GetByIdResp

	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&order).
		Get(fmt.Sprintf("/orders/%d", id))

	if err != nil {
		o.logger.Error("error making request", slog.String("error", err.Error()))
		return nil, err
	}

	if resp.IsError() {
		o.logger.Error("request response error", slog.Int("error_status_code", resp.StatusCode()))
		return nil, errors.New("request response error")
	}

	return &order, nil
}
