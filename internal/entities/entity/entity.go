package entity

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"time"
)

type (
	PaymentMessage struct {
		OrderId int                `json:"order_id"`
		Status  enum.PaymentStatus `json:"status"`
	}

	QRCodePayment struct {
		QRCode  string `json:"qr_code"`
		OrderID int    `json:"-"`
	}

	Payment struct {
		PaymentID    string             `json:"payment_id" dynamodbav:"payment_id"`
		OrderID      int                `json:"order_id" dynamodbav:"order_id"`
		Type         enum.PaymentType   `json:"type" dynamodbav:"type"`
		CurrentState enum.PaymentStatus `json:"status" dynamodbav:"current_state"`
		Amount       float32            `json:"amount" dynamodbav:"amount"`
		CreatedAt    *time.Time         `json:"created_at,omitempty" dynamodbav:"created_at"`
		UpdatedAt    *time.Time         `json:"updated_at,omitempty" dynamodbav:"updated_at"`
		DeletedAt    *time.Time         `json:"deleted_at,omitempty" dynamodbav:"deleted_at"`
	}

	Order struct {
		ID          int              `json:"id" gorm:"primaryKey;autoIncrement"`
		ClientId    *int             `json:"client_id"`
		OrderStatus enum.OrderStatus `json:"order_status"`
		Amount      float32          `json:"amount"`
		CreatedAt   time.Time        `json:"created_at,omitempty"`
	}

	Client struct {
		ID    int    `json:"id"`
		Cpf   int    `json:"cpf"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)
