package sqs

import (
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
)

type FailedProductionMessage struct {
	OrderId int    `json:"order_id"`
	Status  string `json:"status"`
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

	fmt.Println(message)

	if message.Status != string(enum.CANCELED) {
		return nil
	}

	return nil
}
