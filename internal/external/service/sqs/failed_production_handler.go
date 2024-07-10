package sqs

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
)

type FailedProductionMessage struct {
	OrderId int `json:"order_id"`
}

type FailedProductHandler struct {
	payment contract.PaymentUseCase
}

func NewFailedProductionHandler(p contract.PaymentUseCase) *FailedProductHandler {
	return &FailedProductHandler{payment: p}
}

func (f *FailedProductHandler) Handle(message FailedProductionMessage) error {

	return nil
}
