package use_case

import (
	"context"
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"log/slog"
)

var (
	ErrConfirmingPayment = errors.New("error confirming payment")
	ErrCancelingPayment  = errors.New("error canceling payment")
)

type PaymentUseCase struct {
	repository contract.PaymentRepository
	qrCode     contract.PaymentInterface[entity.QRCodePayment]
	snsService contract.SnsService
	logger     *slog.Logger
}

func NewPaymentUseCase(r contract.PaymentRepository, q contract.PaymentInterface[entity.QRCodePayment], s contract.SnsService, logger *slog.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		repository: r,
		qrCode:     q,
		snsService: s,
		logger:     logger,
	}
}

func (p *PaymentUseCase) create(ctx context.Context, payment *entity.Payment) error {
	return p.repository.Create(ctx, *payment)
}

func (p *PaymentUseCase) GetLastPaymentStatus(ctx context.Context, paymentId int) (enum.PaymentStatus, error) {
	payment, err := p.repository.GetLastPaymentStatus(ctx, paymentId)
	if err != nil {
		return enum.PaymentStatusPending, err
	}

	if payment == nil {
		return enum.PaymentStatusPending, nil
	}

	return payment.CurrentState, nil
}

func (p *PaymentUseCase) CreateQRCode(ctx context.Context, order *entity.Order) (*entity.QRCodePayment, error) {
	lastPaymentStatus, err := p.GetLastPaymentStatus(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if lastPaymentStatus == enum.PaymentStatusConfirmed {
		p.logger.Error(fmt.Sprintf("Last payment status: %s", lastPaymentStatus))
		return nil, nil
	}

	payment := entity.Payment{
		OrderID:      order.ID,
		Type:         enum.PaymentTypeQRCode,
		CurrentState: enum.PaymentStatusPending,
		Amount:       order.Amount,
	}

	p.create(ctx, &payment)

	qrCode, err := p.qrCode.Process(payment)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (p *PaymentUseCase) ConfirmedPaymentNotification(ctx context.Context, id int) error {
	return p.processPayment(ctx, id, enum.PaymentStatusConfirmed, ErrConfirmingPayment)
}

func (p *PaymentUseCase) CanceledPaymentNotification(ctx context.Context, id int) error {
	return p.processPayment(ctx, id, enum.PaymentStatusCanceled, ErrCancelingPayment)
}

func (p *PaymentUseCase) processPayment(ctx context.Context, id int, status enum.PaymentStatus, operationErr error) error {
	lastStatus, err := p.GetLastPaymentStatus(ctx, id)
	if err != nil {
		return errors.Join(operationErr, err)
	}

	if err = p.repository.UpdateStatus(ctx, id, status); err != nil {
		return errors.Join(operationErr, err)
	}

	payment := entity.PaymentMessage{
		OrderId: id,
		Status:  status,
	}

	if err = p.snsService.SendMessage(ctx, payment); err != nil {
		// rollback payment status
		if err = p.repository.UpdateStatus(ctx, id, lastStatus); err != nil {
			return errors.Join(operationErr, err)
		}

		return errors.Join(operationErr, err)
	}

	return nil
}
