//go:generate mockgen -destination=mock/use_case.go -source=use_case.go -package=mock
package contract

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
)

type PaymentInterface[T interface{}] interface {
	Process(p entity.Payment) (*T, error)
}

type OrderUseCase interface {
	GetById(id int) (*entity.Order, error)
}

type PaymentUseCase interface {
	CreateQRCode(ctx context.Context, order *entity.Order) (*entity.QRCodePayment, error)
	GetLastPaymentStatus(ctx context.Context, id int) (enum.PaymentStatus, error)
	ConfirmedPaymentNotification(ctx context.Context, id int) error
	CanceledPaymentNotification(ctx context.Context, id int) error
}
