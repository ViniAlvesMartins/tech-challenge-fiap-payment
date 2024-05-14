//go:generate mockgen -destination=mock/service.go -source=service.go -package=mock
package contract

import (
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/order_service"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
)

type ExternalPaymentService interface {
	CreateQRCode(payment entity.Payment) (responsepaymentservice.CreateQRCode, error)
}

type OrderService interface {
	GetById(id int) (*responseorderservice.GetByIdResp, error)
}
