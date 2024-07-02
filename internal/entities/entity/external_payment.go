package entity

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/google/uuid"
)

type ExternalPayment struct {
	ID      *uuid.UUID       `json:"id"`
	OrderID int              `json:"-"`
	Type    enum.PaymentType `json:"type"`
	Amount  float32          `json:"amount"`
}
