package use_case

import (
	"fmt"
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type PaymentUseCase struct {
	repository             contract.PaymentRepository
	externalPaymentService contract.ExternalPaymentService
	snsService             *contract.SnsService
	logger                 *slog.Logger
}

func NewPaymentUseCase(r contract.PaymentRepository, e contract.ExternalPaymentService, s *contract.SnsService, logger *slog.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		repository:             r,
		externalPaymentService: e,
		snsService:             s,
		logger:                 logger,
	}
}

func (p *PaymentUseCase) Create(payment *entity.Payment) error {
	if _, err := p.repository.Create(*payment); err != nil {
		return err
	}

	return nil
}

func (p *PaymentUseCase) GetLastPaymentStatus(paymentId int) (enum.PaymentStatus, error) {

	payment, err := p.repository.GetLastPaymentStatus(paymentId)

	if err != nil && payment != nil {
		return payment.CurrentState, err
	}

	if payment != nil && payment.CurrentState == "" {
		return enum.PENDING, nil
	}

	return payment.CurrentState, nil
}

func (p *PaymentUseCase) CreateQRCode(order *entity.Order) (*responsepaymentservice.CreateQRCode, error) {
	lastPaymentStatus, err := p.GetLastPaymentStatus(order.ID)

	if err != nil {
		return nil, err
	}

	if lastPaymentStatus == enum.CONFIRMED {
		p.logger.Error(fmt.Sprintf("Last payment status: %s", lastPaymentStatus))
		return nil, nil
	}

	payment := &entity.Payment{
		OrderID:      order.ID,
		Type:         enum.QRCODE,
		CurrentState: enum.PENDING,
		Amount:       order.Amount,
	}

	p.Create(payment)

	qrCode, _ := p.externalPaymentService.CreateQRCode(*payment)

	return &qrCode, nil
}

func (p *PaymentUseCase) PaymentNotification(paymentId int) error {

	_, err := p.repository.UpdateStatus(paymentId, enum.CONFIRMED)

	if err != nil {
		panic(err)
	}

	res, err := p.snsService.SendMessage(paymentId, enum.CONFIRMED)

	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	return nil
}
