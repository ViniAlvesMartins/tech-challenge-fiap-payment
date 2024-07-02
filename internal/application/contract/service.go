//go:generate mockgen -destination=mock/service.go -source=service.go -package=mock
package contract

import (
	"context"
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/modules/response/order_service"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/sns_producer"
)

type ExternalPaymentService interface {
	CreateQRCode(payment entity.Payment) (responsepaymentservice.CreateQRCode, error)
}

type OrderService interface {
	GetById(id int) (*responseorderservice.GetByIdResp, error)
}

type SnsService interface {
	SendMessage(ctx context.Context, message sns_producer.Message) error
}
