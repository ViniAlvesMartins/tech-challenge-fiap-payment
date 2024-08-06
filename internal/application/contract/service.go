//go:generate mockery --all --output mock --outpkg mock
package contract

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/order"
)

type OrderService interface {
	GetById(id int) (*order.GetByIdResp, error)
}

type SnsService interface {
	SendMessage(ctx context.Context, message entity.PaymentMessage) error
}
