package contract

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
)

type PaymentRepository interface {
	Create(payment entity.Payment) (*entity.Payment, error)
	GetLastPaymentStatus(orderId int) (*entity.Payment, error)
}
