//go:generate mockgen -destination=mock/repository.go -source=repository.go -package=mock
package contract

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
)

type PaymentRepository interface {
	Create(payment entity.Payment) (*entity.Payment, error)
	GetLastPaymentStatus(orderId int) (*entity.Payment, error)
}
