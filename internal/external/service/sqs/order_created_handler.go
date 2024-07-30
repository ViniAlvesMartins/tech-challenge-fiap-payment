package sqs

import (
	"context"
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"log/slog"
)

type (
	OrderCreatedMessage struct {
		ID          int              `json:"id" gorm:"primaryKey;autoIncrement"`
		OrderStatus enum.OrderStatus `json:"order_status"`
		Amount      float32          `json:"amount"`
	}

	OrderCreatedHandler struct {
		payment contract.PaymentUseCase
		logger  *slog.Logger
	}
)

func NewOrderCreatedHandler(p contract.PaymentUseCase, l *slog.Logger) *OrderCreatedHandler {
	return &OrderCreatedHandler{payment: p, logger: l}
}

func (f *OrderCreatedHandler) Handle(ctx context.Context, b []byte) error {
	var message OrderCreatedMessage

	f.logger.Info("Handling message...")

	if err := json.Unmarshal(b, &message); err != nil {
		return err
	}

	if message.OrderStatus != enum.OrderStatusAwaitingPayment {
		return nil
	}

	order := &entity.Order{
		ID:     message.ID,
		Amount: message.Amount,
	}

	_, err := f.payment.CreateQRCode(ctx, order)

	return err
}
