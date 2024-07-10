//go:generate mockgen -destination=mock/service.go -source=service.go -package=mock
package contract

import (
	"context"
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/modules/response/order_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
)

type OrderService interface {
	GetById(id int) (*responseorderservice.GetByIdResp, error)
}

type SnsService interface {
	SendMessage(ctx context.Context, message entity.PaymentMessage) error
}
