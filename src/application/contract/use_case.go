//go:generate mockgen -destination=mock/use_case.go -source=use_case.go -package=mock
package contract

import (
	"context"
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type OrderUseCase interface {
	GetById(id int) (*entity.Order, error)
}

type PaymentUseCase interface {
	Create(ctx context.Context, payment *entity.Payment) error
	CreateQRCode(ctx context.Context, order *entity.Order) (*response_payment_service.CreateQRCode, error)
	GetLastPaymentStatus(ctx context.Context, paymentId int) (enum.PaymentStatus, error)
	PaymentNotification(ctx context.Context, paymentId int) error
}
