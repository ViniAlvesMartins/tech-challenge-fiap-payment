package sqs

import (
	"context"
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"log/slog"
)

type FailedProductionMessage struct {
	OrderId int    `json:"order_id"`
	Status  string `json:"status"`
}

type FailedProductHandler struct {
	payment contract.PaymentUseCase
	logger  *slog.Logger
}

func NewFailedProductionHandler(p contract.PaymentUseCase, l *slog.Logger) *FailedProductHandler {
	return &FailedProductHandler{payment: p, logger: l}
}

func (f *FailedProductHandler) Handle(ctx context.Context, b []byte) error {
	var message FailedProductionMessage

	f.logger.Info("Handling message...")

	if err := json.Unmarshal(b, &message); err != nil {
		return err
	}

	if message.Status != string(enum.OrderStatusCanceled) {
		return nil
	}

	if err := f.payment.CanceledPaymentNotification(ctx, message.OrderId); err != nil {
		return err
	}

	return nil
}
