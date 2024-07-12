package order

import (
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/go-resty/resty/v2"
	"log/slog"
)

type Service struct {
	client *resty.Client
	logger *slog.Logger
}

type GetByIdResp struct {
	Error string        `json:"error"`
	Data  *entity.Order `json:"data"`
}

func NewService(c *resty.Client, l *slog.Logger) *Service {
	return &Service{logger: l, client: c}
}

func (o *Service) GetById(id int) (*GetByIdResp, error) {
	var order GetByIdResp

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
