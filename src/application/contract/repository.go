package contract

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"
)

type PaymentRepository interface {
	Create(payment entity.Payment) (*entity.Payment, error)
	GetLastPaymentStatus(orderId int) (*entity.Payment, error)
	UpdateStatus(orderId int, status enum.PaymentStatus) (bool, error)
}
