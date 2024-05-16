package use_case

import (
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/contract"
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"
)

type PaymentUseCase struct {
	repository             contract.PaymentRepository
	externalPaymentService contract.ExternalPaymentService
	logger                 *slog.Logger
}

func NewPaymentUseCase(r contract.PaymentRepository, e contract.ExternalPaymentService, logger *slog.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		repository:             r,
		externalPaymentService: e,
		logger:                 logger,
	}
}

func (p *PaymentUseCase) Create(payment *entity.Payment) error {
	p.repository.Create(*payment)
	return nil
}

func (p *PaymentUseCase) GetLastPaymentStatus(orderId int) (enum.PaymentStatus, error) {

	payment, err := p.repository.GetLastPaymentStatus(orderId)

	if err != nil && payment != nil {
		return payment.Status, err
	}

	if payment != nil && payment.Status == "" {
		return enum.PENDING, nil
	}

	return payment.Status, nil
}

func (p *PaymentUseCase) CreateQRCode(order *entity.Order) (*response_payment_service.CreateQRCode, error) {
	//lastPaymentStatus, err := p.GetLastPaymentStatus(order.ID)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//if lastPaymentStatus == enum.CONFIRMED {
	//	p.logger.Error("Last payment status: %v", lastPaymentStatus)
	//	return nil, nil
	//}

	payment := &entity.Payment{
		Type:   enum.QRCODE,
		Status: enum.PENDING,
		Amount: order.Amount,
	}

	p.Create(payment)

	qrCode, _ := p.externalPaymentService.CreateQRCode(*payment)

	return &qrCode, nil
}

func (p *PaymentUseCase) PaymentNotification(order *entity.Order) error {
	payment := &entity.Payment{
		Type:   enum.QRCODE,
		Status: enum.CONFIRMED,
		Amount: order.Amount,
	}

	p.Create(payment)

	// ENVIAR MSG PARA FILA

	return nil
}
