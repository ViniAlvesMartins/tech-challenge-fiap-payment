package entity

import (
	"time"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type Order struct {
	ID          int              `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientId    *int             `json:"client_id"`
	StatusOrder enum.StatusOrder `json:"status_order"`
	Amount      float32          `json:"amount"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
}

func (o *Order) SetAmount(amount float32) {
	o.Amount = amount
}
