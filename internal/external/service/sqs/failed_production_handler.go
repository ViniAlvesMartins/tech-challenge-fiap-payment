package sqs

import (
	"encoding/json"
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

func (f *FailedProductHandler) Handle(b []byte) error {
	var message FailedProductionMessage
	if err := json.Unmarshal(b, &message); err != nil {
		return err
	}

	return nil
}
