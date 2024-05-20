package use_case

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log/slog"
	"strconv"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type PaymentUseCase struct {
	repository             contract.PaymentRepository
	externalPaymentService contract.ExternalPaymentService
	snsService             contract.SnsService
	logger                 *slog.Logger
}

func NewPaymentUseCase(r contract.PaymentRepository, e contract.ExternalPaymentService, s contract.SnsService, logger *slog.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		repository:             r,
		externalPaymentService: e,
		snsService:             s,
		logger:                 logger,
	}
}

func (p *PaymentUseCase) Create(ctx context.Context, payment *entity.Payment) error {
	if _, err := p.repository.Create(ctx, *payment); err != nil {
		return err
	}

	return nil
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

func (p *PaymentUseCase) CreateQRCode(ctx context.Context, order *entity.Order) (*responsepaymentservice.CreateQRCode, error) {
	lastPaymentStatus, err := p.GetLastPaymentStatus(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if lastPaymentStatus == enum.CONFIRMED {
		p.logger.Error(fmt.Sprintf("Last payment status: %s", lastPaymentStatus))
		return nil, nil
	}

	payment := &entity.Payment{
		OrderID:      strconv.Itoa(order.ID),
		Type:         enum.QRCODE,
		CurrentState: enum.PENDING,
		Amount:       order.Amount,
	}

	p.Create(ctx, payment)

	qrCode, _ := p.externalPaymentService.CreateQRCode(*payment)

	return &qrCode, nil
}

func (p *PaymentUseCase) PaymentNotification(ctx context.Context, paymentId int) error {
	if err := p.repository.UpdateStatus(ctx, paymentId, enum.CONFIRMED); err != nil {
		return err
	}

	if err := p.snsService.SendMessage(paymentId, enum.CONFIRMED); err != nil {
		return err
	}

	return nil
}
