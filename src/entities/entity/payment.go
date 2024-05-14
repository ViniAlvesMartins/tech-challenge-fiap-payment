package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type PaymentInterface interface {
	GETJSONValue() string
}

type Payment struct {
	ID        string             `json:"id"`
	OrderID   int                `json:"orderId"`
	Type      enum.PaymentType   `json:"type"`
	Status    enum.PaymentStatus `json:"status"`
	Amount    float32            `json:"amount"`
	CreatedAt *time.Time         `json:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty"`
}

func (p *Payment) GetJSONValue() (string, error) {
	b, err := json.Marshal(p)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(b), nil
}
