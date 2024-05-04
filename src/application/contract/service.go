package contract

import (
	response_order_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/order_service"
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
)

type ExternalPaymentService interface {
	CreateQRCode(payment entity.Payment) (response_payment_service.CreateQRCode, error)
}

type OrderService interface {
	GetById(orderId int) (*response_order_service.GetByIdResp, error)
}