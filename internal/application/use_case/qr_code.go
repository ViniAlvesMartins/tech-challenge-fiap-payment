package use_case

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	mercadopago "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/external/service/external_payment/mercado_pago"
	"strconv"
)

type QRCode struct {
	paymentGateway *mercadopago.PaymentGateway
}

func NewQRCode(g *mercadopago.PaymentGateway) *QRCode {
	return &QRCode{paymentGateway: g}
}

func (q *QRCode) Process(p entity.Payment) (*entity.QRCodePayment, error) {
	response, err := q.paymentGateway.CreateQRCode(p)
	if err != nil {
		return nil, err
	}

	orderID, err := strconv.Atoi(response.InStoreOrderID)
	if err != nil {
		return nil, err
	}

	return &entity.QRCodePayment{
		QRCode:  response.QRData,
		OrderID: orderID,
	}, nil
}
