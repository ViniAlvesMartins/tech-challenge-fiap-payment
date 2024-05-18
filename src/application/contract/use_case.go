//go:generate mockgen -destination=mock/use_case.go -source=use_case.go -package=mock
package contract

import (
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type OrderUseCase interface {
	GetById(id int) (*entity.Order, error)
}

type PaymentUseCase interface {
	Create(payment *entity.Payment) error
	CreateQRCode(order *entity.Order) (*response_payment_service.CreateQRCode, error)
	GetLastPaymentStatus(paymentId int) (enum.PaymentStatus, error)
	PaymentNotification(paymentId int) error
}
