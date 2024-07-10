package use_case

import (
	"context"
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log/slog"
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
	var notFoundErr *types.ResourceNotFoundException

	payment, err := p.repository.GetLastPaymentStatus(ctx, paymentId)
	if err != nil {
		if errors.As(err, &notFoundErr) {
			return enum.PENDING, nil
		}

		return enum.PENDING, err
	}

	if payment != nil && payment.CurrentState == "" {
		return enum.PENDING, nil
	}

	return payment.CurrentState, nil
}

func (p *PaymentUseCase) CreateQRCode(ctx context.Context, order *entity.Order) (*entity.QRCodePayment, error) {
	lastPaymentStatus, err := p.GetLastPaymentStatus(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if lastPaymentStatus == enum.CONFIRMED {
		p.logger.Error(fmt.Sprintf("Last payment status: %s", lastPaymentStatus))
		return nil, nil
	}

	payment := entity.Payment{
		OrderID:      order.ID,
		Type:         enum.QRCODE,
		CurrentState: enum.PENDING,
		Amount:       order.Amount,
	}

	p.create(ctx, &payment)

	qrCode, err := p.qrCode.Process(payment)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (p *PaymentUseCase) PaymentNotification(ctx context.Context, paymentId int) error {
	if err := p.repository.UpdateStatus(ctx, paymentId, enum.CONFIRMED); err != nil {
		return err
	}

	payment := entity.PaymentMessage{
		OrderId: paymentId,
		Status:  enum.CONFIRMED,
	}

	if err := p.snsService.SendMessage(ctx, payment); err != nil {
		return err
	}

	return nil
}

//func (p *PaymentUseCase) CancelPayment(ctx context.Context, orderId int) error {
//	err := p.repository.UpdateStatus(ctx, paymentId, enum.CONFIRMED); err != nil {
//		return err
//	}
//
//
//
//	if err := p.snsService.SendMessage(ctx, payment); err != nil {
//		return err
//	}
//
//	return nil
//}
