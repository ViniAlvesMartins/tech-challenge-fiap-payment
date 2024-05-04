package service

import (
	response_payment_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
)

type ExternalPayment struct {
}

func NewExternalPayment() *ExternalPayment {
	return &ExternalPayment{}
}

func (e *ExternalPayment) CreateQRCode(p entity.Payment) (response_payment_service.CreateQRCode, error) {
	response := response_payment_service.CreateQRCode{
		QrData: "00020101021243650016COM.MERCADOLIBRE02013063638f1192a-5fd1-4180-a180-8bcae3556bc35204000053039865802BR5925IZABEL AAAA DE MELO6007BARUERI62070503***63040B6D",
	}

	return response, nil
}
